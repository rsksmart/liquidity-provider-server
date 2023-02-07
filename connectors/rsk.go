package connectors

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"

	"github.com/ethereum/go-ethereum/rpc"

	gethTypes "github.com/ethereum/go-ethereum/core/types"

	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"

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

type QuotePegOutWatcher interface {
	OnRegisterPegOut(newState types.RQState)
	OnExpire()
	Done() <-chan struct{}
}

type RSKConnector interface {
	Connect(endpoint string, chainId *big.Int) error
	CheckConnection() error
	Close()
	GetChainId() (*big.Int, error)
	EstimateGas(addr string, value *big.Int, data []byte) (uint64, error)
	GasPrice() (*big.Int, error)
	HashQuote(q *pegin.Quote) (string, error)
	HashPegOutQuote(q *pegout.Quote) (string, error)
	ParseQuote(q *pegin.Quote) (bindings.LiquidityBridgeContractQuote, error)
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
	GetLbcBalance(addr string) (*big.Int, error)
	GetAvailableLiquidity(addr string) (*big.Int, error)
	GetTxStatus(ctx context.Context, tx *gethTypes.Transaction) (bool, error)
	GetMinimumLockTxValue() (*big.Int, error)
	FetchFederationInfo() (*FedInfo, error)
	AddQuoteToWatch(hash string, interval time.Duration, exp time.Time, w QuotePegOutWatcher, cb RegisterPegOutQuoteWatcherCompleteCallback) error
	GetRskHeight() (uint64, error)
	GetProviders() ([]bindings.LiquidityBridgeContractProvider, error)
	GetDerivedBitcoinAddress(fedInfo *FedInfo, btcParams chaincfg.Params, userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error)
	GetActivePowpegRedeemScript() ([]byte, error)
}

type RSKClient interface {
	ChainID(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*gethTypes.Receipt, error)
	Close()
}

type RSKBridge interface {
	GetActiveFederationCreationBlockHeight(opts *bind.CallOpts) (*big.Int, error)
	GetFederationSize(opts *bind.CallOpts) (*big.Int, error)
	GetFederationThreshold(opts *bind.CallOpts) (*big.Int, error)
	GetFederationAddress(opts *bind.CallOpts) (string, error)
	GetFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, arg1 string) ([]byte, error)
	GetMinimumLockTxValue(opts *bind.CallOpts) (*big.Int, error)
	GetActivePowpegRedeemScript(opts *bind.CallOpts) ([]byte, error)
}

type RSK struct {
	c                           *ethclient.Client
	lbc                         *bindings.LiquidityBridgeContract
	lbcAddress                  common.Address
	bridge                      *bindings.RskBridge
	bridgeAddress               common.Address
	requiredBridgeConfirmations int64
	irisActivationHeight        int
	erpKeys                     []string
}

type RegisterPegOutQuoteWatcherCompleteCallback = func(w QuotePegOutWatcher)

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
	rsk.lbc, err = bindings.NewLiquidityBridgeContract(rsk.lbcAddress, rsk.c)
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

func (rsk *RSK) GetLbcBalance(addr string) (*big.Int, error) {
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("invalid address: %v", addr)
	}
	a := common.HexToAddress(addr)
	var err error
	for i := 0; i < retries; i++ {
		var bal *big.Int
		bal, err = rsk.lbc.GetBalance(&bind.CallOpts{}, a)
		if err == nil {
			return bal, nil
		}
		time.Sleep(rpcSleep)
	}
	return nil, fmt.Errorf("error getting %v balance: %v", addr, err)
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

func (rsk *RSK) GetProcessedPegOutQuotes(quoteHash [32]byte) (uint8, error) {
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var status uint8
		status, err = rsk.lbc.GetPegOutProcessedQuote(&bind.CallOpts{
			Context: ctx,
		}, quoteHash)
		if err == nil {
			return status, nil
		}

		log.Debugf("Exp:: GetProcessedPegOutQuotes error ::: %v", err)
		time.Sleep(rpcSleep)
	}

	return 0, fmt.Errorf("error retrieving processed pegout status: %v", err)
}

func (rsk *RSK) EstimateGas(addr string, value *big.Int, data []byte) (uint64, error) {
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
		Value: new(big.Int).Set(value),
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

func (rsk *RSK) GasPrice() (*big.Int, error) {
	var err error
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		var price *big.Int
		price, err = rsk.c.SuggestGasPrice(ctx)
		if price != nil && price.Cmp(big.NewInt(0)) >= 0 {
			return price, nil
		}
		time.Sleep(rpcSleep)
	}
	return nil, fmt.Errorf("error estimating gas: %v", err)
}

func (rsk *RSK) HashPegOutQuote(q *pegout.Quote) (string, error) {
	opts := bind.CallOpts{}
	var results [32]byte

	pq, err := rsk.ParsePegOutQuote(q)
	if err != nil {
		return "", err
	}

	for i := 0; i < retries; i++ {
		results, err = rsk.lbc.HashPegoutQuote(&opts, pq)
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

func (rsk *RSK) HashQuote(q *pegin.Quote) (string, error) {
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
	lbcCaller := &bindings.LiquidityBridgeContractCallerRaw{Contract: &rsk.lbc.LiquidityBridgeContractCaller}
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

func (rsk *RSK) GetDerivedBitcoinAddress(fedInfo *FedInfo, btcParams chaincfg.Params, userBtcRefundAddr []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) (string, error) {
	derivationValue, err := getDerivationValueHash(userBtcRefundAddr, lbcAddress, lpBtcAddress, derivationArgumentsHash)
	if err != nil {
		return "", fmt.Errorf("error computing derivation value: %v", err)
	}
	var fedRedeemScript []byte
	fedRedeemScript, err = rsk.GetActivePowpegRedeemScript()
	if err != nil {
		return "", fmt.Errorf("error retreiving fed redeem script from bridge: %v", err)
	}
	if len(fedRedeemScript) == 0 {
		fedRedeemScript, err = fedInfo.getFedRedeemScript(btcParams)
		if err != nil {
			return "", fmt.Errorf("error generating fed redeem script: %v", err)
		}
	} else {
		err = fedInfo.validateRedeemScript(btcParams, fedRedeemScript)
		if err != nil {
			return "", fmt.Errorf("error validating fed redeem script: %v", err)
		}
	}
	flyoverScript, err := fedInfo.getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if err != nil {
		return "", fmt.Errorf("error generating flyover redeem script: %v", err)
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(flyoverScript, &btcParams)
	if err != nil {
		return "", err
	}
	return addressScriptHash.EncodeAddress(), nil
}

// GetActivePowpegRedeemScript returns a PowPeg redeem script fetched from the RSK bridge.
// It returns a PowPeg redeem script, if the method is activated on the bridge. Otherwise - empty result.
// It returns an error, if encountered a communication issue with the bridge.
func (rsk *RSK) GetActivePowpegRedeemScript() ([]byte, error) {
	var err error
	opts := bind.CallOpts{}
	var value []byte
	for i := 0; i < retries; i++ {
		value, err = rsk.bridge.GetActivePowpegRedeemScript(&opts)
		if err == nil || isNoContractError(err) {
			break
		}
		time.Sleep(rpcSleep)
	}
	if err != nil {
		if isNoContractError(err) {
			return []byte{}, nil
		}
		return nil, fmt.Errorf("error calling GetActivePowpegRedeemScript: %v", err)
	}
	return value, nil
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
	trim := strings.TrimPrefix(address, "0x")
	if !common.IsHexAddress(trim) {
		return nil, fmt.Errorf("invalid address: %v", address)
	}
	return common.HexToAddress(trim).Bytes(), nil
}

func (rsk *RSK) ParseQuote(q *pegin.Quote) (bindings.LiquidityBridgeContractQuote, error) {
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
	pq.CallFee = q.CallFee.Copy().AsBigInt()
	pq.PenaltyFee = q.PenaltyFee.Copy().AsBigInt()
	pq.GasLimit = q.GasLimit
	pq.Nonce = q.Nonce
	pq.Value = q.Value.Copy().AsBigInt()
	pq.AgreementTimestamp = q.AgreementTimestamp
	pq.CallTime = q.CallTime
	pq.DepositConfirmations = q.Confirmations
	pq.TimeForDeposit = q.TimeForDeposit
	return pq, nil
}

func (rsk *RSK) ParsePegOutQuote(q *pegout.Quote) (bindings.LiquidityBridgeContractPegOutQuote, error) {
	pq := bindings.LiquidityBridgeContractPegOutQuote{}

	if err := copyHex(q.LBCAddr, pq.LbcAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractPegOutQuote{}, fmt.Errorf("error parsing LBC address: %v", err)
	}
	if err := copyHex(q.LPRSKAddr, pq.LiquidityProviderRskAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractPegOutQuote{}, fmt.Errorf("error parsing provider RSK address: %v", err)
	}
	if err := copyHex(q.RSKRefundAddr, pq.RskRefundAddress[:]); err != nil {
		return bindings.LiquidityBridgeContractPegOutQuote{}, fmt.Errorf("error parsing RSK refund address: %v", err)
	}

	pq.Fee = q.Fee
	pq.PenaltyFee = q.PenaltyFee
	pq.Nonce = q.Nonce
	pq.ValueToTransfer = q.Value
	pq.AgreementTimestamp = q.AgreementTimestamp
	pq.DepositDateLimit = q.DepositDateLimit
	pq.DepositConfirmations = q.DepositConfirmations
	pq.TransferConfirmations = q.TransferConfirmations
	pq.TransferTime = q.TransferTime
	pq.ExpireDate = q.ExpireDate
	pq.ExpireBlocks = q.ExpireBlocks

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

func (rsk *RSK) AddQuoteToWatch(hash string, interval time.Duration, exp time.Time, w QuotePegOutWatcher, cb RegisterPegOutQuoteWatcherCompleteCallback) error {
	go func(w QuotePegOutWatcher) {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				rsk.checkPegoutRegister(hash, w, exp, time.Now)
			case <-w.Done():
				ticker.Stop()
				cb(w)
				return
			}
		}
	}(w)
	return nil
}

func (rsk *RSK) GetRskHeight() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
	defer cancel()

	number, err := rsk.c.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (rsk *RSK) checkPegoutRegister(quoteHash string, w QuotePegOutWatcher, expTime time.Time, now func() time.Time) {
	log.Debugf("checkPegoutRegister was started %v", quoteHash)

	log.Debugf("Exp:: time %v\n", expTime)
	log.Debugf("Exp:: now %v\n", now())
	log.Debugf("Exp:: quoteHash %s\n", quoteHash)

	if now().After(expTime) {
		log.Errorf("time for registerPegout expired %s", quoteHash)
		w.OnExpire()
		return
	}

	b, err := hex.DecodeString(quoteHash)

	if err != nil {
		log.Errorln(err.Error())
		return
	}

	var b32 [32]byte

	copy(b32[:], b)

	status, err := rsk.GetProcessedPegOutQuotes(b32)

	if err != nil {
		w.OnRegisterPegOut(types.RQStateCallForUserFailed)
		return
	}

	log.Debugf("Exp:: status %v\n", status)

	if status == 2 {
		w.OnRegisterPegOut(types.RQStateCallForUserSucceeded)
	}
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
	bts, err := hex.DecodeString(strings.TrimPrefix(str, "0x"))
	if err != nil {
		return nil, err
	}
	return bts, nil
}

func isNoContractError(err error) bool {
	return "no contract code at given address" == err.Error()
}

func (rsk *RSK) GetProviders() ([]bindings.LiquidityBridgeContractProvider, error) {
	opts := bind.CallOpts{}
	var err error
	var providers []bindings.LiquidityBridgeContractProvider

	providers, err = rsk.lbc.GetProviders(&opts)

	if err != nil {
		log.Debug("Error RSK.go", err)
	}

	return providers, err
}
