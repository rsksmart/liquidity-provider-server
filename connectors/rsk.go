package connectors

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"net/http"
	"net/url"

	gethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/btcsuite/btcutil/base58"
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"

	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider/types"

	log "github.com/sirupsen/logrus"
)

const (
	retries    int = 3
	rpcSleep       = 2 * time.Second
	rpcTimeout     = 5 * time.Second
	ethSleep       = 5 * time.Second
	ethTimeout     = 5 * time.Minute

	newAccountGasCost = uint64(25000)
)

type RSKConnector interface {
	Connect(endpoint string, chainId *big.Int) error
	CheckConnection() error
	Close()
	GetChainId() (*big.Int, error)
	EstimateGas(addr string, value uint64, data []byte) (uint64, error)
	GasPrice() (uint64, error)
	HashQuote(q *types.Quote) (string, error)
	ParseQuote(q *types.Quote) (bindings.LiquidityBridgeContractQuote, error)
	RegisterPegIn(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote, signature []byte, tx []byte, pmt []byte, height *big.Int) (*gethTypes.Transaction, error)
	GetFedSize() (int, error)
	GetFedThreshold() (int, error)
	GetFedPublicKey(index int) (string, error)
	GetFedAddress() (string, error)
	GetActiveFederationCreationBlockHeight() (int, error)
	GetLBCAddress() string
	GetRequiredBridgeConfirmations() int64
	CallForUser(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote) (*gethTypes.Transaction, error)
	RegisterPegInWithoutTx(q bindings.LiquidityBridgeContractQuote, signature []byte, tx []byte, pmt []byte, newInt *big.Int) error
	GetCollateral(addr string) (*big.Int, *big.Int, error)
	RegisterProvider(opts *bind.TransactOpts) error
	AddCollateral(opts *bind.TransactOpts) error
	GetAvailableLiquidity(addr string) (*big.Int, error)
	GetTxStatus(ctx context.Context, tx *gethTypes.Transaction) (bool, error)
	GetMinimumLockTxValue() (*big.Int, error)
	FetchFederationInfo() (*FedInfo, error)
}

type RSK struct {
	c                           *ethclient.Client
	lbc                         *bindings.LBC
	lbcAddress                  common.Address
	bridge                      *bindings.RskBridge
	bridgeAddress               common.Address
	requiredBridgeConfirmations int64
	irisActivationHeight        int
	erpKeys                     []string
}

func NewRSK(lbcAddress string, bridgeAddress string, requiredBridgeConfirmations int64, irisActivationHeight int, erpKeys []string) (*RSK, error) {
	if !common.IsHexAddress(lbcAddress) {
		return nil, errors.New("invalid LBC contract address")
	}
	if !common.IsHexAddress(bridgeAddress) {
		return nil, errors.New("invalid Bridge contract address")
	}

	return &RSK{
		lbcAddress:                  common.HexToAddress(lbcAddress),
		bridgeAddress:               common.HexToAddress(bridgeAddress),
		requiredBridgeConfirmations: requiredBridgeConfirmations,
		irisActivationHeight:        irisActivationHeight,
		erpKeys:                     erpKeys,
	}, nil
}

func (rsk *RSK) Connect(endpoint string, chainId *big.Int) error {
	log.Debug("connecting to RSK node on ", endpoint)

	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	var ethC *ethclient.Client
	switch u.Scheme {
	case "http", "https":
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.DisableKeepAlives = true

		httpC := new(http.Client)
		httpC.Transport = transport

		c, err := rpc.DialHTTPWithClient(endpoint, httpC)
		if err != nil {
			return err
		}

		ethC = ethclient.NewClient(c)
	default:
		ethC, err = ethclient.Dial(endpoint)
		if err != nil {
			return err
		}
	}

	rsk.c = ethC

	log.Debug("verifying connection to RSK node")
	// test connection
	rskChainId, err := rsk.GetChainId()
	if err != nil {
		return err
	}
	//check chain id
	if chainId.Cmp(rskChainId) != 0 {
		return fmt.Errorf("chain id mismatch; expected chain id: %v, rsk node chain id: %v", chainId, rskChainId)
	}

	log.Debug("initializing RSK contracts")
	rsk.bridge, err = bindings.NewRskBridge(rsk.bridgeAddress, rsk.c)
	if err != nil {
		return err
	}
	rsk.lbc, err = bindings.NewLBC(rsk.lbcAddress, rsk.c)
	if err != nil {
		return err
	}
	return nil
}

func (rsk *RSK) CheckConnection() error {
	_, err := rsk.GetChainId()
	return err
}

func (rsk *RSK) Close() {
	log.Debug("closing RSK connection")
	rsk.c.Close()
}

func (rsk *RSK) GetAvailableLiquidity(addr string) (*big.Int, error) {
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("invalid address: %v", addr)
	}
	a := common.HexToAddress(addr)
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()
	var err error
	var liq *big.Int
	for i := 0; i < retries; i++ {
		liq, err = rsk.c.BalanceAt(ctx, a, nil)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting balance of %v: %v", addr, err)
	}
	for i := 0; i < retries; i++ {
		var bal *big.Int
		bal, err = rsk.lbc.GetBalance(&bind.CallOpts{}, a)
		if err == nil {
			return liq.Add(liq, bal), nil
		}
		time.Sleep(rpcSleep)
	}
	return nil, fmt.Errorf("error getting %v balance: %v", addr, err)
}

func (rsk *RSK) GetCollateral(addr string) (*big.Int, *big.Int, error) {
	if !common.IsHexAddress(addr) {
		return nil, nil, fmt.Errorf("invalid address: %v", addr)
	}
	a := common.HexToAddress(addr)
	var (
		min *big.Int
		col *big.Int
		err error
	)
	for i := 0; i < retries; i++ {
		min, err = rsk.lbc.GetMinCollateral(&bind.CallOpts{})
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error getting minimum collateral: %v", err)
	}
	for i := 0; i < retries; i++ {
		col, err = rsk.lbc.GetCollateral(&bind.CallOpts{}, a)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error getting collateral: %v", err)
	}
	return col, min, nil
}

func (rsk *RSK) RegisterProvider(opts *bind.TransactOpts) error {
	var err error
	var tx *gethTypes.Transaction
	for i := 0; i < retries; i++ {
		tx, err = rsk.lbc.Register(opts)
		if err == nil && tx != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if tx == nil || err != nil {
		return fmt.Errorf("error registering provider: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ethTimeout)
	defer cancel()
	s, err := rsk.GetTxStatus(ctx, tx)
	if err != nil || !s {
		return fmt.Errorf("error getting tx status while registering provider: %v", err)
	}
	return nil
}

func (rsk *RSK) AddCollateral(opts *bind.TransactOpts) error {
	var err error
	var tx *gethTypes.Transaction
	for i := 0; i < retries; i++ {
		tx, err = rsk.lbc.AddCollateral(opts)
		if err == nil && tx != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if tx == nil || err != nil {
		return fmt.Errorf("error adding collateral: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ethTimeout)
	defer cancel()
	s, err := rsk.GetTxStatus(ctx, tx)
	if err != nil || !s {
		return fmt.Errorf("error getting tx status while adding collateral: %v", err)
	}
	return nil
}

func (rsk *RSK) GetChainId() (*big.Int, error) {
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var chainId *big.Int
		chainId, err = rsk.c.ChainID(ctx)
		if err == nil {
			return chainId, nil
		}
		time.Sleep(rpcSleep)
	}
	return nil, fmt.Errorf("error retrieving chain id: %v", err)
}

func (rsk *RSK) EstimateGas(addr string, value uint64, data []byte) (uint64, error) {
	if !common.IsHexAddress(addr) {
		return 0, fmt.Errorf("invalid address: %v", addr)
	}

	dst := common.HexToAddress(addr)

	var additionalGas uint64
	if rsk.isNewAccount(dst) {
		additionalGas = newAccountGasCost
	}

	msg := ethereum.CallMsg{
		To:    &dst,
		Data:  data,
		Value: big.NewInt(int64(value)),
	}

	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var gas uint64
		gas, err = rsk.c.EstimateGas(ctx, msg)
		if gas > 0 {
			return gas + additionalGas, nil
		}
		time.Sleep(rpcSleep)
	}
	return 0, fmt.Errorf("error estimating gas: %v", err)
}

func (rsk *RSK) GasPrice() (uint64, error) {
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var price *big.Int
		price, err = rsk.c.SuggestGasPrice(ctx)
		if price != nil && price.Cmp(big.NewInt(0)) >= 0 {
			return price.Uint64(), nil
		}
		time.Sleep(rpcSleep)
	}
	return 0, fmt.Errorf("error estimating gas: %v", err)
}

func (rsk *RSK) HashQuote(q *types.Quote) (string, error) {
	opts := bind.CallOpts{}
	var results [32]byte

	pq, err := rsk.ParseQuote(q)
	if err != nil {
		return "", err
	}

	for i := 0; i < retries; i++ {
		results, err = rsk.lbc.HashQuote(&opts, pq)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if err != nil {
		return "", fmt.Errorf("error calling HashQuote: %v", err)
	}
	return hex.EncodeToString(results[:]), nil
}

func (rsk *RSK) GetFedSize() (int, error) {
	var err error
	opts := bind.CallOpts{}
	var results *big.Int

	for i := 0; i < retries; i++ {
		results, err = rsk.bridge.GetFederationSize(&opts)
		if results != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if err != nil {
		return 0, fmt.Errorf("error calling GetFederationSize: %v", err)
	}

	sizeInt, err := strconv.Atoi(results.String())
	if err != nil {
		return 0, fmt.Errorf("error converting federation size to int. error: %v", err)
	}
	return sizeInt, nil
}

func (rsk *RSK) GetFedThreshold() (int, error) {
	var err error
	opts := bind.CallOpts{}
	var results *big.Int

	for i := 0; i < retries; i++ {
		results, err = rsk.bridge.GetFederationThreshold(&opts)
		if results != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if err != nil {
		return 0, fmt.Errorf("error calling GetFederationThreshold: %v", err)
	}

	sizeInt, err := strconv.Atoi(results.String())
	if err != nil {
		return 0, fmt.Errorf("error converting federation size to int. error: %v", err)
	}

	return sizeInt, nil
}

func (rsk *RSK) GetFedPublicKey(index int) (string, error) {
	var err error
	var results []byte
	opts := bind.CallOpts{}

	for i := 0; i < retries; i++ {
		results, err = rsk.bridge.GetFederatorPublicKeyOfType(&opts, big.NewInt(int64(index)), "btc")
		if len(results) > 0 {
			break
		}
		time.Sleep(rpcSleep)
	}
	if len(results) == 0 {
		return "", fmt.Errorf("error calling GetFederatorPublicKeyOfType: %v", err)
	}

	return hex.EncodeToString(results), nil
}

func (rsk *RSK) GetFedAddress() (string, error) {
	var err error
	var results string
	opts := bind.CallOpts{}

	for i := 0; i < retries; i++ {
		results, err = rsk.bridge.GetFederationAddress(&opts)
		if results != "" {
			break
		}
		time.Sleep(rpcSleep)
	}
	if results == "" {
		return "", fmt.Errorf("error calling GetFederationAddress: %v", err)
	}
	return results, nil
}

func (rsk *RSK) GetActiveFederationCreationBlockHeight() (int, error) {
	var err error
	opts := bind.CallOpts{}
	var results *big.Int
	for i := 0; i < retries; i++ {
		results, err = rsk.bridge.GetActiveFederationCreationBlockHeight(&opts)
		if results != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if results == nil {
		return 0, fmt.Errorf("error calling getActiveFederationCreationBlockHeight: %v", err)
	}
	height, err := strconv.Atoi(results.String())
	if err != nil {
		return 0, err
	}
	return height, nil
}

func (rsk *RSK) GetRequiredBridgeConfirmations() int64 {
	return rsk.requiredBridgeConfirmations
}
func (rsk *RSK) GetLBCAddress() string {
	return rsk.lbcAddress.String()
}

func (rsk *RSK) CallForUser(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote) (*gethTypes.Transaction, error) {
	var err error
	var tx *gethTypes.Transaction
	for i := 0; i < retries; i++ {
		tx, err = rsk.lbc.CallForUser(opt, q)
		if err == nil && tx != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if tx == nil && err != nil {
		return nil, fmt.Errorf("error calling callForUser: %v", err)
	}
	return tx, nil
}

func (rsk *RSK) RegisterPegIn(opt *bind.TransactOpts, q bindings.LiquidityBridgeContractQuote, signature []byte, tx []byte, pmt []byte, height *big.Int) (*gethTypes.Transaction, error) {
	var err error
	var t *gethTypes.Transaction
	for i := 0; i < retries; i++ {
		t, err = rsk.lbc.RegisterPegIn(opt, q, signature, tx, pmt, height)
		if err == nil && t != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if tx == nil && err != nil {
		return nil, fmt.Errorf("error calling registerPegIn: %v", err)
	}
	return t, nil
}

func (rsk *RSK) RegisterPegInWithoutTx(q bindings.LiquidityBridgeContractQuote, signature []byte, tx []byte, pmt []byte, height *big.Int) error {
	var res []interface{}
	lbcCaller := &bindings.LBCCallerRaw{Contract: &rsk.lbc.LBCCaller}
	err := lbcCaller.Call(&bind.CallOpts{}, &res, "registerPegIn", q, signature, tx, pmt, height)
	if err != nil {
		return err
	}
	return nil
}

func (rsk *RSK) GetTxStatus(ctx context.Context, tx *gethTypes.Transaction) (bool, error) {
	ticker := time.NewTicker(ethSleep)

	for {
		select {
		case <-ticker.C:
			cctx, cancel := context.WithTimeout(ctx, rpcTimeout)
			defer cancel()
			r, _ := rsk.c.TransactionReceipt(cctx, tx.Hash())
			if r != nil {
				return r.Status == 1, nil
			}
		case <-ctx.Done():
			ticker.Stop()
			return false, fmt.Errorf("operation cancelled")
		}
	}
}

func (rsk *RSK) isNewAccount(addr common.Address) bool {
	var (
		err  error
		code []byte
		bal  *big.Int
		n    uint64
	)
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		code, err = rsk.c.CodeAt(ctx, addr, nil)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		bal, err = rsk.c.BalanceAt(ctx, addr, nil)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		n, err = rsk.c.NonceAt(ctx, addr, nil)
		if err == nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	return len(code) == 0 && bal.Cmp(common.Big0) == 0 && n == 0
}

func (rsk *RSK) GetMinimumLockTxValue() (*big.Int, error) {
	var err error
	opts := bind.CallOpts{}
	var value *big.Int
	for i := 0; i < retries; i++ {
		value, err = rsk.bridge.GetMinimumLockTxValue(&opts)
		if value != nil {
			break
		}
		time.Sleep(rpcSleep)
	}
	if value == nil {
		return nil, fmt.Errorf("error calling GetMinimumLockTxValue: %v", err)
	}
	return value, nil
}

func DecodeRSKAddress(address string) ([]byte, error) {
	trim := strings.Trim(address, "0x")
	if !common.IsHexAddress(trim) {
		return nil, fmt.Errorf("invalid address: %v", address)
	}
	return common.HexToAddress(trim).Bytes(), nil
}

func (rsk *RSK) ParseQuote(q *types.Quote) (bindings.LiquidityBridgeContractQuote, error) {
	pq := bindings.LiquidityBridgeContractQuote{}
	var err error

	if err := copyBtcAddr(q.FedBTCAddr, pq.FedBtcAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing federation address: %v", err)
	}
	if pq.LiquidityProviderBtcAddress, err = DecodeBTCAddressWithVersion(q.LPBTCAddr); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing bitcoin liquidity provider address: %v", err)
	}
	if pq.BtcRefundAddress, err = DecodeBTCAddressWithVersion(q.BTCRefundAddr); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing bitcoin refund address: %v", err)
	}
	if err := copyHex(q.LBCAddr, pq.LbcAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing LBC address: %v", err)
	}
	if err := copyHex(q.LPRSKAddr, pq.LiquidityProviderRskAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing provider RSK address: %v", err)
	}
	if err := copyHex(q.RSKRefundAddr, pq.RskRefundAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing RSK refund address: %v", err)
	}
	if err := copyHex(q.ContractAddr, pq.ContractAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing contract address: %v", err)
	}
	if pq.Data, err = parseHex(q.Data); err != nil {
		return bindings.LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing data: %v", err)
	}
	pq.CallFee = q.CallFee
	pq.PenaltyFee = q.PenaltyFee
	pq.GasLimit = q.GasLimit
	pq.Nonce = q.Nonce
	pq.Value = q.Value
	pq.AgreementTimestamp = q.AgreementTimestamp
	pq.CallTime = q.CallTime
	pq.DepositConfirmations = q.Confirmations
	pq.TimeForDeposit = q.TimeForDeposit
	return pq, nil
}

func (rsk *RSK) FetchFederationInfo() (*FedInfo, error) {
	log.Debug("getting federation info")
	fedSize, err := rsk.GetFedSize()
	if err != nil {
		return nil, err
	}

	var pubKeys []string
	for i := 0; i < fedSize; i++ {
		pubKey, err := rsk.GetFedPublicKey(i)
		if err != nil {
			log.Error("error fetching fed public key: ", err.Error())
			return nil, err
		}
		pubKeys = append(pubKeys, pubKey)
	}

	fedThreshold, err := rsk.GetFedThreshold()
	if err != nil {
		log.Error("error fetching federation size: ", err.Error())
		return nil, err
	}

	fedAddress, err := rsk.GetFedAddress()
	if err != nil {
		return nil, err
	}

	activeFedBlockHeight, err := rsk.GetActiveFederationCreationBlockHeight()
	if err != nil {
		log.Error("error fetching federation address: ", err.Error())
		return nil, err
	}

	return &FedInfo{
		FedThreshold:         fedThreshold,
		FedSize:              fedSize,
		PubKeys:              pubKeys,
		FedAddress:           fedAddress,
		ActiveFedBlockHeight: activeFedBlockHeight,
		IrisActivationHeight: rsk.irisActivationHeight,
		ErpKeys:              rsk.erpKeys,
	}, nil
}

func copyBtcAddr(addr string, dst []byte) error {
	addressBts, _, err := base58.CheckDecode(addr)
	if err != nil {
		return err
	}
	copy(dst, addressBts)
	return nil
}

func copyHex(str string, dst []byte) error {
	bts, err := parseHex(str)
	if err != nil {
		return err
	}
	copy(dst, bts[:])
	return nil
}

func parseHex(str string) ([]byte, error) {
	bts, err := hex.DecodeString(strings.Trim(str, "0x"))
	if err != nil {
		return nil, err
	}
	return bts, nil
}
