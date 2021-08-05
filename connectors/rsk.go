package connectors

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

type quote struct {
	FedBTCAddr         [20]byte `abi:"fedBtcAddress"`
	LBCAddr            [20]byte `abi:"lbcAddress"`
	LPRSKAddr          [20]byte `abi:"liquidityProviderRskAddress"`
	BTCRefundAddr      []byte   `abi:"btcRefundAddress"`
	RSKRefundAddr      [20]byte `abi:"rskRefundAddress"`
	LPBTCAddr          []byte   `abi:"liquidityProviderBtcAddress"`
	CallFee            *big.Int `abi:"callFee"`
	ContractAddr       [20]byte `abi:"contractAddress"`
	Data               []byte   `abi:"data"`
	GasLimit           *big.Int `abi:"gasLimit"`
	Nonce              *big.Int `abi:"nonce"`
	Value              *big.Int `abi:"value"`
	AgreementTimestamp *big.Int `abi:"agreementTimestamp"`
	TimeForDeposit     *big.Int `abi:"timeForDeposit"`
	CallTime           *big.Int `abi:"callTime"`
	Confirmations      *big.Int `abi:"depositConfirmations"`
}

type RSK struct {
	c          *ethclient.Client
	lbc        *bind.BoundContract
	abi        *abi.ABI
	lbcAddress common.Address
}

func NewRSK(lbcAddress string, abiPath string) (*RSK, error) {
	if !common.IsHexAddress(lbcAddress) {
		return nil, errors.New("invalid contract address")
	}

	def, err := loadLBCABI(abiPath)
	if err != nil {
		return nil, err
	}

	return &RSK{abi: def, lbcAddress: common.HexToAddress(lbcAddress)}, nil
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
	log.Debug("connected to RSK node")
	rsk.lbc = bind.NewBoundContract(rsk.lbcAddress, *rsk.abi, rsk.c, rsk.c, rsk.c)
	return nil
}

func (rsk *RSK) Close() {
	log.Debug("closing RSK connection")
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
	var err error
	results := new([]interface{})
	opts := bind.CallOpts{}

	pq, err := parseQuote(q)
	if err != nil {
		return "", err
	}

	for i := 0; i < retries; i++ {
		err = rsk.lbc.Call(&opts, results, "hashQuote", pq)
		if len(*results) > 0 {
			break
		}
		time.Sleep(sleepTime)
	}
	if len(*results) == 0 {
		return "", fmt.Errorf("error calling hashQuote %v: %v", pq, err)
	}
	arr := *results
	bts := getBytes(arr[0])

	return hex.EncodeToString(bts), nil
}

func getBytes(key interface{}) ([]byte) {
	var bts []byte
	for _, bt := range key.([32]byte) {
		bts = append(bts, bt)
	}

	return bts
}

func parseQuote(q *types.Quote) (*quote, error) {
	pq := quote{}
	var err error

	if err := copyHex(q.FedBTCAddr, pq.FedBTCAddr[:]); err != nil {
		return nil, fmt.Errorf("error parsing federation address: %v", err)
	}
	if err := copyHex(q.LBCAddr, pq.LBCAddr[:]); err != nil {
		return nil, fmt.Errorf("error parsing LBC address: %v", err)
	}
	if err := copyHex(q.LPRSKAddr, pq.LPRSKAddr[:]); err != nil {
		return nil, fmt.Errorf("error parsing provider RSK address: %v", err)
	}
	if err := copyHex(q.RSKRefundAddr, pq.RSKRefundAddr[:]); err != nil {
		return nil, fmt.Errorf("error parsing RSK refund address: %v", err)
	}
	if err := copyHex(q.ContractAddr, pq.ContractAddr[:]); err != nil {
		return nil, fmt.Errorf("error parsing contract address: %v", err)
	}
	if pq.BTCRefundAddr, err = parseHex(q.BTCRefundAddr); err != nil {
		return nil, fmt.Errorf("error parsing BTC refund address: %v", err)
	}
	if pq.LPBTCAddr, err = parseHex(q.LPBTCAddr); err != nil {
		return nil, fmt.Errorf("error parsing provider BTC address: %v", err)
	}
	if pq.Data, err = parseHex(q.Data); err != nil {
		return nil, fmt.Errorf("error parsing data: %v", err)
	}
	pq.CallFee = &q.CallFee
	pq.GasLimit = new(big.Int).SetUint64(uint64(q.GasLimit))
	pq.Nonce = new(big.Int).SetUint64(uint64(q.Nonce))
	pq.Value = &q.Value
	pq.AgreementTimestamp = new(big.Int).SetUint64(uint64(q.AgreementTimestamp))
	pq.CallTime = new(big.Int).SetUint64(uint64(q.CallTime))
	pq.Confirmations = new(big.Int).SetUint64(uint64(q.Confirmations))
	pq.TimeForDeposit = new(big.Int).SetUint64(uint64(q.TimeForDeposit))
	return &pq, nil
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

func loadLBCABI(path string) (*abi.ABI, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	definition, err := abi.JSON(f)
	if err != nil {
		return nil, err
	}
	return &definition, nil
}
