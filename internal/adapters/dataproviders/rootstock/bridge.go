package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/federation"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	log "github.com/sirupsen/logrus"
	"math/big"
	"slices"
	"time"
)

const registerCoinbaseTxGasLimit = 100000

type rskBridgeImpl struct {
	address               string
	requiredConfirmations uint64
	irisActivationHeight  int64
	erpKeys               []string
	contract              RskBridgeAdapter
	client                RpcClientBinding
	btcParams             *chaincfg.Params
	retryParams           RetryParams
	signer                TransactionSigner
	miningTimeout         time.Duration
}

type RskBridgeConfig struct {
	Address               string
	RequiredConfirmations uint64
	IrisActivationHeight  int64
	ErpKeys               []string
}

func NewRskBridgeImpl(
	config RskBridgeConfig,
	contract RskBridgeAdapter,
	client *RskClient,
	btcParams *chaincfg.Params,
	retryParams RetryParams,
	signer TransactionSigner,
	miningTimeout time.Duration,
) rootstock.Bridge {
	return &rskBridgeImpl{
		address:               config.Address,
		requiredConfirmations: config.RequiredConfirmations,
		irisActivationHeight:  config.IrisActivationHeight,
		erpKeys:               config.ErpKeys,
		contract:              contract,
		client:                client.client,
		btcParams:             btcParams,
		retryParams:           retryParams,
		signer:                signer,
		miningTimeout:         miningTimeout,
	}
}

func (bridge *rskBridgeImpl) GetAddress() string {
	return bridge.address
}

func (bridge *rskBridgeImpl) GetFedAddress() (string, error) {
	opts := &bind.CallOpts{}
	return rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
		func() (string, error) {
			return bridge.contract.GetFederationAddress(opts)
		})
}

func (bridge *rskBridgeImpl) GetMinimumLockTxValue() (*entities.Wei, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
		func() (*big.Int, error) {
			return bridge.contract.GetMinimumLockTxValue(opts)
		})
	if err != nil {
		return nil, err
	}
	// This value comes in satoshi from the bridge
	return entities.SatoshiToWei(result.Uint64()), nil
}

func (bridge *rskBridgeImpl) GetFlyoverDerivationAddress(args rootstock.FlyoverDerivationArgs) (rootstock.FlyoverDerivation, error) {
	var err error
	var fedRedeemScript, derivationValue []byte
	derivationValue = federation.GetDerivationValueHash(args)
	opts := &bind.CallOpts{}
	fedRedeemScript, err = rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
		func() ([]byte, error) {
			return bridge.contract.GetActivePowpegRedeemScript(opts)
		})
	if err != nil {
		return rootstock.FlyoverDerivation{}, fmt.Errorf("error retreiving fed redeem script from bridge: %w", err)
	}

	return federation.CalculateFlyoverDerivationAddress(args.FedInfo, *bridge.btcParams, fedRedeemScript, derivationValue)
}

func (bridge *rskBridgeImpl) GetRequiredTxConfirmations() uint64 {
	return bridge.requiredConfirmations
}

func (bridge *rskBridgeImpl) FetchFederationInfo() (rootstock.FederationInfo, error) {
	var err error
	var pubKey []byte
	var pubKeys []string
	var i, federationSize int64

	opts := &bind.CallOpts{}
	fedSize, err := rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
		func() (*big.Int, error) {
			return bridge.contract.GetFederationSize(opts)
		})
	if err != nil {
		return rootstock.FederationInfo{}, err
	}
	federationSize = fedSize.Int64()

	for i = 0; i < federationSize; i++ {
		pubKey, err = rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
			func() ([]byte, error) {
				return bridge.contract.GetFederatorPublicKeyOfType(opts, big.NewInt(i), "btc")
			})
		if err != nil {
			return rootstock.FederationInfo{}, fmt.Errorf("error fetching fed public key: %w", err)
		}
		pubKeys = append(pubKeys, hex.EncodeToString(pubKey))
	}

	fedThreshold, err := rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
		func() (*big.Int, error) {
			return bridge.contract.GetFederationThreshold(opts)
		})
	if err != nil {
		return rootstock.FederationInfo{}, fmt.Errorf("error fetching federation size: %w", err)
	}

	fedAddress, err := rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
		func() (string, error) {
			return bridge.contract.GetFederationAddress(opts)
		})
	if err != nil {
		return rootstock.FederationInfo{}, fmt.Errorf("error fetching federation address: %w", err)
	}

	activeFedBlockHeight, err := rskRetry(bridge.retryParams.Retries, bridge.retryParams.Sleep,
		func() (*big.Int, error) {
			return bridge.contract.GetActiveFederationCreationBlockHeight(opts)
		})
	if err != nil {
		return rootstock.FederationInfo{}, fmt.Errorf("error fetching federation height: %w", err)
	}

	return rootstock.FederationInfo{
		FedThreshold:         fedThreshold.Int64(),
		FedSize:              fedSize.Int64(),
		PubKeys:              pubKeys,
		FedAddress:           fedAddress,
		ActiveFedBlockHeight: activeFedBlockHeight.Int64(),
		IrisActivationHeight: bridge.irisActivationHeight,
		ErpKeys:              bridge.erpKeys,
	}, nil
}

// RegisterBtcCoinbaseTransaction registers a new Bitcoin coinbase transaction in the bridge. Returns blockchain.WaitingForBridgeError
// if the transaction has not been observed by the bridge yet. If the transaction was already registered, it returns an empty string instead of the hash.
func (bridge *rskBridgeImpl) RegisterBtcCoinbaseTransaction(params rootstock.BtcCoinbaseTransactionInformation) (string, error) {
	var err error
	var alreadyRegistered bool
	var bestChainHeight *big.Int

	if bestChainHeight, err = bridge.contract.GetBtcBlockchainBestChainHeight(&bind.CallOpts{}); err != nil {
		return "", fmt.Errorf("error validating if coinbase transaction was processed by the bridge: %w", err)
	} else if bestChainHeight.Cmp(params.BlockHeight) < 0 {
		return "", blockchain.WaitingForBridgeError
	}

	if alreadyRegistered, err = bridge.contract.HasBtcBlockCoinbaseTransactionInformation(&bind.CallOpts{}, params.BlockHash); alreadyRegistered {
		log.Info("Coinbase transaction already registered")
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("error validating if coinbase transaction was registered: %w", err)
	}

	log.Infof("Executing RegisterBtcCoinbaseTransaction with params: %s\n", params.String())
	opts := &bind.TransactOpts{
		From:     bridge.signer.Address(),
		Signer:   bridge.signer.Sign,
		GasLimit: registerCoinbaseTxGasLimit,
	}

	receipt, err := awaitTx(bridge.client, bridge.miningTimeout, "RegisterBtcCoinbaseTransaction", func() (*geth.Transaction, error) {
		return bridge.contract.RegisterBtcCoinbaseTransaction(opts, params.BtcTxSerialized, params.BlockHash,
			params.SerializedPmt, params.WitnessMerkleRoot, params.WitnessReservedValue)
	})

	if err != nil {
		return "", fmt.Errorf("register coinbase transaction error: %w", err)
	} else if receipt == nil {
		return "", errors.New("register coinbase transaction error: incomplete receipt")
	} else if receipt.Status == 0 {
		txHash := receipt.TxHash.String()
		return txHash, fmt.Errorf("register coinbase transaction error: transaction reverted (%s)", txHash)
	}
	return receipt.TxHash.String(), nil
}

func (bridge *rskBridgeImpl) GetBatchPegOutCreatedEvent(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]rootstock.BatchPegOut, error) {
	var event *bindings.RskBridgeBatchPegoutCreated
	result := make([]rootstock.BatchPegOut, 0)

	rawIterator, err := bridge.contract.FilterBatchPegoutCreated(&bind.FilterOpts{
		Start:   fromBlock,
		End:     toBlock,
		Context: ctx,
	}, nil)

	iterator := bridge.contract.BatchPegOutCreatedIteratorAdapter(rawIterator)
	defer func() {
		if rawIterator != nil {
			if iteratorError := iterator.Close(); iteratorError != nil {
				log.Error("Error closing BatchPegOutCreated event iterator: ", err)
			}
		}
	}()
	if err != nil {
		return nil, err
	} else if rawIterator == nil {
		return nil, fmt.Errorf("no BatchPegOutCreated events found in the range %d to %d", fromBlock, *toBlock)
	}

	var rskHashes []string
	for iterator.Next() {
		event = iterator.Event()
		if rskHashes, err = parseReleaseRskHashes(event.ReleaseRskTxHashes); err != nil {
			return nil, fmt.Errorf("error parsing release RSK hashes: %w", err)
		}
		result = append(result, rootstock.BatchPegOut{
			TransactionHash:    event.Raw.TxHash.String(),
			BlockHash:          event.Raw.BlockHash.String(),
			BlockNumber:        event.Raw.BlockNumber,
			BtcTxHash:          hex.EncodeToString(event.BtcTxHash[:]),
			ReleaseRskTxHashes: rskHashes,
		})
	}
	if err = iterator.Error(); err != nil {
		return nil, err
	}

	return result, nil
}

func parseReleaseRskHashes(hashes []byte) ([]string, error) {
	const hashSize = 32
	chunks := slices.Chunk(hashes, hashSize)
	result := make([]string, 0)
	for chunk := range chunks {
		if len(chunk) != hashSize {
			return nil, fmt.Errorf("invalid release RSK hash size: expected %d bytes, got %d bytes", hashSize, len(chunk))
		}
		result = append(result, "0x"+hex.EncodeToString(chunk))
	}
	return result, nil
}
