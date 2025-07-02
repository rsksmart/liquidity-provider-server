package mempool_space

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	merkle "github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type MempoolSpaceApi struct {
	url    string
	client *http.Client
	config *chaincfg.Params
}

type transactionInfoApiResponse struct {
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight uint64 `json:"block_height"`
		BlockHash   string `json:"block_hash"`
	} `json:"status"`
}

func NewMempoolSpaceApi(
	client *http.Client,
	config *chaincfg.Params,
	url string,
) blockchain.BitcoinNetwork {
	return &MempoolSpaceApi{
		client: client,
		config: config,
		url:    strings.TrimSuffix(url, "/"),
	}
}

func (api *MempoolSpaceApi) ValidateAddress(address string) error {
	if err := api.validateNetwork(); err != nil {
		return err
	}
	const errorTemplate = "unable to validate address: %w"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, api.url+"/v1/validate-address/"+address, nil)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	var result struct {
		IsValid bool   `json:"isvalid"`
		Error   string `json:"error"`
	}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	if result.IsValid && blockchain.IsSupportedBtcAddress(address) {
		return nil
	} else if !blockchain.IsSupportedBtcAddress(address) {
		return blockchain.BtcAddressNotSupportedError
	} else {
		return errors.Join(blockchain.BtcAddressInvalidNetworkError, errors.New(result.Error))
	}
}

func (api *MempoolSpaceApi) DecodeAddress(address string) ([]byte, error) {
	if err := api.validateNetwork(); err != nil {
		return nil, err
	}
	return bitcoin.DecodeAddress(address)
}

func (api *MempoolSpaceApi) GetTransactionInfo(hash string) (blockchain.BitcoinTransactionInformation, error) {
	if err := api.validateNetwork(); err != nil {
		return blockchain.BitcoinTransactionInformation{}, err
	}
	const transactionInfoError = "error getting transaction info: %w"

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api.url+"/tx/"+hash, nil)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, fmt.Errorf(transactionInfoError, err)
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, fmt.Errorf(transactionInfoError, err)
	}
	var transactionInfoResult transactionInfoApiResponse
	err = json.NewDecoder(res.Body).Decode(&transactionInfoResult)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, fmt.Errorf(transactionInfoError, err)
	}

	txBytes, err := api.getRawTransaction(hash, true)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, fmt.Errorf(transactionInfoError, err)
	}
	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, fmt.Errorf(transactionInfoError, err)
	}

	outputs, err := api.buildTransactionOutputs(tx)
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, fmt.Errorf(transactionInfoError, err)
	}

	height, err := api.GetHeight()
	if err != nil {
		return blockchain.BitcoinTransactionInformation{}, fmt.Errorf(transactionInfoError, err)
	}

	return blockchain.BitcoinTransactionInformation{
		Hash:          hash,
		Confirmations: height.Uint64() - transactionInfoResult.Status.BlockHeight,
		Outputs:       outputs,
		HasWitness:    tx.HasWitness(),
	}, nil
}

func (api *MempoolSpaceApi) buildTransactionOutputs(tx *btcutil.Tx) (map[string][]*entities.Wei, error) {
	var addresses []btcutil.Address
	var scriptType txscript.ScriptClass
	var amounts []*entities.Wei
	var address string
	var ok bool
	var err error
	outputs := make(map[string][]*entities.Wei, len(tx.MsgTx().TxOut))
	for _, out := range tx.MsgTx().TxOut {
		scriptType, addresses, _, err = txscript.ExtractPkScriptAddrs(out.PkScript, api.config)
		if err != nil {
			return make(map[string][]*entities.Wei, 0), fmt.Errorf("error building transaction outputs: %w", err)
		} else if len(addresses) == 0 && scriptType != txscript.NullDataTy {
			return make(map[string][]*entities.Wei, 0), errors.New("error getting transaction info: no addresses found in output script")
		} else if scriptType == txscript.NullDataTy {
			outputs[""] = []*entities.Wei{entities.NewWei(0)}
		} else if scriptType != txscript.MultiSigTy && scriptType != txscript.NonStandardTy {
			address = addresses[0].EncodeAddress()
			amounts, ok = outputs[address]
			if !ok {
				amounts = make([]*entities.Wei, 0)
			}
			amounts = append(amounts, entities.SatoshiToWei(uint64(out.Value)))
			outputs[address] = amounts
		}
	}
	return outputs, nil
}

func (api *MempoolSpaceApi) GetRawTransaction(hash string) ([]byte, error) {
	return api.getRawTransaction(hash, false)
}

func (api *MempoolSpaceApi) getRawTransaction(hash string, includeWitness bool) ([]byte, error) {
	if err := api.validateNetwork(); err != nil {
		return nil, err
	}
	const rawTransactionError = "error getting raw transaction: %w"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/tx/%s/hex", api.url, hash), nil)
	if err != nil {
		return []byte{}, fmt.Errorf(rawTransactionError, err)
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return []byte{}, fmt.Errorf(rawTransactionError, err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf(rawTransactionError, err)
	}
	txBytes, err := hex.DecodeString(string(data))
	if err != nil {
		return []byte{}, fmt.Errorf(rawTransactionError, err)
	}

	if includeWitness {
		return txBytes, nil
	}

	tx, err := btcutil.NewTxFromBytes(txBytes)
	if err != nil {
		return []byte{}, fmt.Errorf(rawTransactionError, err)
	}

	buff := new(bytes.Buffer)
	if err = tx.MsgTx().SerializeNoWitness(buff); err != nil {
		return []byte{}, fmt.Errorf(rawTransactionError, err)
	}
	return buff.Bytes(), nil
}

func (api *MempoolSpaceApi) GetPartialMerkleTree(hash string) ([]byte, error) {
	if err := api.validateNetwork(); err != nil {
		return nil, err
	}
	const getPmtError = "error getting partial merkle tree for tx %s: %w"
	blockInfo, err := api.GetTransactionBlockInfo(hash)
	if err != nil {
		return nil, fmt.Errorf(getPmtError, hash, err)
	}
	block, err := api.getBlock(hex.EncodeToString(blockInfo.Hash[:]))
	if err != nil {
		return nil, fmt.Errorf(getPmtError, hash, err)
	}
	parsedHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return []byte{}, fmt.Errorf(getPmtError, hash, err)
	}
	result, err := bitcoin.SerializePartialMerkleTree(parsedHash, block)
	if err != nil {
		return []byte{}, fmt.Errorf(getPmtError, hash, err)
	}
	return result, nil
}

func (api *MempoolSpaceApi) GetHeight() (*big.Int, error) {
	if err := api.validateNetwork(); err != nil {
		return nil, err
	}
	const getHeightError = "error getting blockchain height: %w"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, api.url+"/blocks/tip/height", nil)
	if err != nil {
		return nil, fmt.Errorf(getHeightError, err)
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return nil, fmt.Errorf(getHeightError, err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(getHeightError, err)
	}
	result := new(big.Int)
	_, ok := result.SetString(string(data), 10)
	if !ok {
		return nil, errors.New(getHeightError + "invalid response format")
	}
	return result, nil
}

func (api *MempoolSpaceApi) BuildMerkleBranch(txHash string) (blockchain.MerkleBranch, error) {
	if err := api.validateNetwork(); err != nil {
		return blockchain.MerkleBranch{}, err
	}
	const merkleBranchError = "error building merkle branch: %w"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/tx/%s/merkle-proof", api.url, txHash), nil)
	if err != nil {
		return blockchain.MerkleBranch{}, fmt.Errorf(merkleBranchError, err)
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return blockchain.MerkleBranch{}, fmt.Errorf(merkleBranchError, err)
	} else if res.StatusCode == http.StatusBadRequest {
		return blockchain.MerkleBranch{}, fmt.Errorf("error building merkle branch: transaction %s not found", txHash)
	}
	var hashBytes []byte
	var merkleProof struct {
		Merkle []string `json:"merkle"`
		Pos    int64    `json:"pos"`
	}
	err = json.NewDecoder(res.Body).Decode(&merkleProof)
	if err != nil {
		return blockchain.MerkleBranch{}, fmt.Errorf(merkleBranchError, err)
	}
	hashes := make([][32]byte, len(merkleProof.Merkle))
	for i, hashStr := range merkleProof.Merkle {
		hashBytes, err = hex.DecodeString(hashStr)
		if err != nil {
			return blockchain.MerkleBranch{}, fmt.Errorf(merkleBranchError, err)
		}
		hashes[i] = utils.To32Bytes(hashBytes)
	}
	return blockchain.MerkleBranch{
		Hashes: hashes,
		Path:   big.NewInt(merkleProof.Pos),
	}, nil

}

func (api *MempoolSpaceApi) GetTransactionBlockInfo(txHash string) (blockchain.BitcoinBlockInformation, error) {
	if err := api.validateNetwork(); err != nil {
		return blockchain.BitcoinBlockInformation{}, err
	}
	const getTransactionBlockError = "error getting block of the transaction %s: %w"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, api.url+"/tx/"+txHash, nil)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf(getTransactionBlockError, txHash, err)
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf(getTransactionBlockError, txHash, err)
	}
	var transactionInfoResult transactionInfoApiResponse
	err = json.NewDecoder(res.Body).Decode(&transactionInfoResult)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf(getTransactionBlockError, txHash, err)
	} else if !transactionInfoResult.Status.Confirmed {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf("error getting block of the transaction %s: transaction not confirmet", txHash)
	}

	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, api.url+"/block/"+transactionInfoResult.Status.BlockHash, nil)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf(getTransactionBlockError, txHash, err)
	}
	res, err = api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf(getTransactionBlockError, txHash, err)
	}
	var blockInfo struct {
		Height    uint64 `json:"height"`
		Timestamp int64  `json:"timestamp"`
	}
	err = json.NewDecoder(res.Body).Decode(&blockInfo)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf(getTransactionBlockError, txHash, err)
	}

	hashBytes, err := hex.DecodeString(transactionInfoResult.Status.BlockHash)
	if err != nil {
		return blockchain.BitcoinBlockInformation{}, fmt.Errorf(getTransactionBlockError, txHash, err)
	}

	return blockchain.BitcoinBlockInformation{
		Hash:   utils.To32Bytes(hashBytes),
		Height: new(big.Int).SetUint64(blockInfo.Height),
		Time:   time.Unix(blockInfo.Timestamp, 0),
	}, nil
}

func (api *MempoolSpaceApi) GetCoinbaseInformation(txHash string) (blockchain.BtcCoinbaseTransactionInformation, error) {
	if err := api.validateNetwork(); err != nil {
		return blockchain.BtcCoinbaseTransactionInformation{}, err
	}
	const getCoinbaseInformationError = "error getting coinbase information of the transaction %s: %w"
	var coinbaseTxHash *chainhash.Hash
	var witnessReservedValue [32]byte

	blockInfo, err := api.GetTransactionBlockInfo(txHash)
	if err != nil {
		return blockchain.BtcCoinbaseTransactionInformation{}, fmt.Errorf(getCoinbaseInformationError, txHash, err)
	}
	block, err := api.getBlock(hex.EncodeToString(blockInfo.Hash[:]))
	if err != nil {
		return blockchain.BtcCoinbaseTransactionInformation{}, fmt.Errorf(getCoinbaseInformationError, txHash, err)
	}

	serializedCoinbase := new(bytes.Buffer)
	for _, tx := range block.Transactions() {
		if merkle.IsCoinBaseTx(tx.MsgTx()) {
			if err = tx.MsgTx().SerializeNoWitness(serializedCoinbase); err != nil {
				return blockchain.BtcCoinbaseTransactionInformation{}, err
			}
			coinbaseTxHash = tx.Hash()
			copy(witnessReservedValue[:], [][]byte(tx.MsgTx().TxIn[0].Witness)[0])
		}
	}

	pmt, err := bitcoin.SerializePartialMerkleTree(coinbaseTxHash, block)
	if err != nil {
		return blockchain.BtcCoinbaseTransactionInformation{}, err
	}

	return blockchain.BtcCoinbaseTransactionInformation{
		BtcTxSerialized:      serializedCoinbase.Bytes(),
		BlockHash:            bitcoin.ToSwappedBytes32(block.Hash()),
		BlockHeight:          blockInfo.Height,
		SerializedPmt:        pmt,
		WitnessMerkleRoot:    bitcoin.ToSwappedBytes32(merkle.CalcMerkleRoot(block.Transactions(), true)),
		WitnessReservedValue: witnessReservedValue,
	}, nil
}

func (api *MempoolSpaceApi) NetworkName() string {
	return strings.ToLower(api.config.Name)
}

func (api *MempoolSpaceApi) GetBlockchainInfo() (blockchain.BitcoinBlockchainInfo, error) {
	if err := api.validateNetwork(); err != nil {
		return blockchain.BitcoinBlockchainInfo{}, err
	}
	const getBlockchainInfoError = "error getting blockchain info: %w"

	blocks, err := api.GetHeight()
	if err != nil {
		return blockchain.BitcoinBlockchainInfo{}, fmt.Errorf(getBlockchainInfoError, err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, api.url+"/blocks/tip/hash", nil)
	if err != nil {
		return blockchain.BitcoinBlockchainInfo{}, fmt.Errorf(getBlockchainInfoError, err)
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return blockchain.BitcoinBlockchainInfo{}, fmt.Errorf(getBlockchainInfoError, err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return blockchain.BitcoinBlockchainInfo{}, fmt.Errorf(getBlockchainInfoError, err)
	}

	return blockchain.BitcoinBlockchainInfo{
		NetworkName:      api.NetworkName(),
		ValidatedBlocks:  blocks,
		ValidatedHeaders: blocks,
		BestBlockHash:    string(data),
	}, nil
}

func (api *MempoolSpaceApi) getBlock(blockHash string) (*btcutil.Block, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("%s/block/%s/raw", api.url, blockHash), nil)
	if err != nil {
		return nil, err
	}
	res, err := api.client.Do(req)
	defer closeBody(res)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return btcutil.NewBlockFromBytes(data)
}

func (api *MempoolSpaceApi) validateNetwork() error {
	supportedNetworks := []*chaincfg.Params{&chaincfg.TestNet3Params, &chaincfg.MainNetParams}
	for _, network := range supportedNetworks {
		if api.config.Name == network.Name {
			return nil
		}
	}
	return errors.New("unsupported network")
}

func closeBody(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}
	if err := res.Body.Close(); err != nil {
		log.Error("Error closing body in MempoolSpace API:", err)
	}
}
