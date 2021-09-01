package connectors

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
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
	retries   int           = 3
	sleepTime time.Duration = 2 * time.Second
)

type RSK struct {
	c             *ethclient.Client
	lbc           *LBC
	lbcAddress    common.Address
	bridge        *RskBridge
	bridgeAddress common.Address
}

func NewRSK(lbcAddress string, bridgeAddress string) (*RSK, error) {
	if !common.IsHexAddress(lbcAddress) {
		return nil, errors.New("invalid LBC contract address")
	}
	if !common.IsHexAddress(bridgeAddress) {
		return nil, errors.New("invalid Bridge contract address")
	}

	return &RSK{
		lbcAddress:    common.HexToAddress(lbcAddress),
		bridgeAddress: common.HexToAddress(bridgeAddress),
	}, nil
}

func (rsk *RSK) Connect(endpoint string) error {
	log.Debug("connecting to RSK node on ", endpoint)
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return err
	}
	rsk.c = client

	// test connection
	if _, err := rsk.GasPrice(); err != nil {
		return err
	}
	log.Debug("Verified connection to node successfully")
	rsk.bridge, err = NewRskBridge(rsk.bridgeAddress, rsk.c)
	if err != nil {
		return err
	}
	rsk.lbc, err = NewLBC(rsk.lbcAddress, rsk.c)
	if err != nil {
		return err
	}

	log.Debug("Connected to RSK contracts")

	return nil
}

func (rsk *RSK) Close() {
	log.Debug("Closing RSK connection")
	rsk.c.Close()
}

func (rsk *RSK) EstimateGas(addr string, value big.Int, data []byte) (uint64, error) {
	if !common.IsHexAddress(addr) {
		return 0, fmt.Errorf("invalid address: %v", addr)
	}

	dst := common.HexToAddress(addr)

	msg := ethereum.CallMsg{
		To:    &dst,
		Data:  data,
		Value: &value,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var err error
	for i := 0; i < retries; i++ {
		var gas uint64
		gas, err = rsk.c.EstimateGas(ctx, msg)
		if gas > 0 {
			return gas, nil
		}
		time.Sleep(sleepTime)
	}
	return 0, fmt.Errorf("error estimating gas: %v", err)
}

func (rsk *RSK) GasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var err error
	for i := 0; i < retries; i++ {
		var price *big.Int
		price, err = rsk.c.SuggestGasPrice(ctx)
		if price != nil && price.Cmp(big.NewInt(0)) > 0 {
			return price, nil
		}
		time.Sleep(sleepTime)
	}
	return nil, fmt.Errorf("error estimating gas: %v", err)
}

func (rsk *RSK) HashQuote(q *types.Quote) (string, error) {
	opts := bind.CallOpts{}
	var results [32]byte

	pq, err := parseQuote(q)
	if err != nil {
		return "", err
	}

	for i := 0; i < retries; i++ {
		results, err = rsk.lbc.HashQuote(&opts, pq)
		if err == nil {
			break
		}
		time.Sleep(sleepTime)
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
		time.Sleep(sleepTime)
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
		time.Sleep(sleepTime)
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

func (rsk *RSK) GetFedPublicKeyOfType(index int) (string, error) {
	var err error
	var results []byte
	opts := bind.CallOpts{}

	for i := 0; i < retries; i++ {
		results, err = rsk.bridge.GetFederatorPublicKeyOfType(&opts, big.NewInt(int64(index)), "btc")
		if len(results) > 0 {
			break
		}
		time.Sleep(sleepTime)
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
		time.Sleep(sleepTime)
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
		time.Sleep(sleepTime)
	}
	if results == nil {
		return 0, fmt.Errorf("error calling getActiveFederationCreationBlockHeight: %v", err)
	}
	height, err := strconv.Atoi(results.String())

	return height, nil
}

func parseQuote(q *types.Quote) (LiquidityBridgeContractQuote, error) {
	pq := LiquidityBridgeContractQuote{}
	var err error

	if err := copyHex(q.FedBTCAddr, pq.FedBtcAddress[:]); err != nil {
		return LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing federation address: %v", err)
	}
	if err := copyHex(q.LBCAddr, pq.LbcAddress[:]); err != nil {
		return LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing LBC address: %v", err)
	}
	if err := copyHex(q.LPRSKAddr, pq.LiquidityProviderRskAddress[:]); err != nil {
		return LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing provider RSK address: %v", err)
	}
	if err := copyHex(q.RSKRefundAddr, pq.RskRefundAddress[:]); err != nil {
		return LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing RSK refund address: %v", err)
	}
	if err := copyHex(q.ContractAddr, pq.ContractAddress[:]); err != nil {
		return LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing contract address: %v", err)
	}
	if pq.Data, err = parseHex(q.Data); err != nil {
		return LiquidityBridgeContractQuote{}, fmt.Errorf("error parsing data: %v", err)
	}
	pq.CallFee = &q.CallFee
	pq.PenaltyFee = &q.PenaltyFee
	pq.GasLimit = new(big.Int).SetUint64(uint64(q.GasLimit))
	pq.Nonce = new(big.Int).SetUint64(uint64(q.Nonce))
	pq.Value = &q.Value
	pq.AgreementTimestamp = new(big.Int).SetUint64(uint64(q.AgreementTimestamp))
	pq.CallTime = new(big.Int).SetUint64(uint64(q.CallTime))
	pq.DepositConfirmations = new(big.Int).SetUint64(uint64(q.Confirmations))
	pq.TimeForDeposit = new(big.Int).SetUint64(uint64(q.TimeForDeposit))
	return pq, nil
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
