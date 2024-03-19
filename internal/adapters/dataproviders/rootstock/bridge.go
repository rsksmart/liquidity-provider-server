package rootstock

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
)

type rskBridgeImpl struct {
	address               string
	requiredConfirmations uint64
	irisActivationHeight  int64
	erpKeys               []string
	contract              *bindings.RskBridge
	client                *ethclient.Client
	btcParams             *chaincfg.Params
}

func NewRskBridgeImpl(
	address string,
	requiredConfirmations uint64,
	irisActivationHeight int64,
	erpKeys []string,
	contract *bindings.RskBridge,
	client *RskClient,
	btcParams *chaincfg.Params,
) blockchain.RootstockBridge {
	return &rskBridgeImpl{
		address:               address,
		requiredConfirmations: requiredConfirmations,
		irisActivationHeight:  irisActivationHeight,
		erpKeys:               erpKeys,
		contract:              contract,
		client:                client.client,
		btcParams:             btcParams,
	}
}

func (bridge *rskBridgeImpl) GetAddress() string {
	return bridge.address
}

func (bridge *rskBridgeImpl) GetFedAddress() (string, error) {
	opts := &bind.CallOpts{}
	return rskRetry(func() (string, error) {
		return bridge.contract.GetFederationAddress(opts)
	})
}

func (bridge *rskBridgeImpl) GetMinimumLockTxValue() (*entities.Wei, error) {
	opts := &bind.CallOpts{}
	result, err := rskRetry(func() (*big.Int, error) {
		return bridge.contract.GetMinimumLockTxValue(opts)
	})
	if err != nil {
		return nil, err
	}
	// This value comes in satoshi from the bridge
	return entities.SatoshiToWei(result.Uint64()), nil
}

func (bridge *rskBridgeImpl) GetFlyoverDerivationAddress(args blockchain.FlyoverDerivationArgs) (blockchain.FlyoverDerivation, error) {
	var err error
	var fedRedeemScript, derivationValue, flyoverScript []byte
	var addressScriptHash *btcutil.AddressScriptHash

	if derivationValue = bridge.getDerivationValueHash(args); err != nil {
		return blockchain.FlyoverDerivation{}, fmt.Errorf("error computing derivation value: %w", err)
	}
	opts := &bind.CallOpts{}
	fedRedeemScript, err = rskRetry(func() ([]byte, error) {
		return bridge.contract.GetActivePowpegRedeemScript(opts)
	})
	if err != nil {
		return blockchain.FlyoverDerivation{}, fmt.Errorf("error retreiving fed redeem script from bridge: %w", err)
	}

	if len(fedRedeemScript) == 0 {
		if fedRedeemScript, err = getFedRedeemScript(args.FedInfo, *bridge.btcParams); err != nil {
			return blockchain.FlyoverDerivation{}, fmt.Errorf("error generating fed redeem script: %w", err)
		}
	} else {
		if err = validateRedeemScript(args.FedInfo, *bridge.btcParams, fedRedeemScript); err != nil {
			return blockchain.FlyoverDerivation{}, fmt.Errorf("error validating fed redeem script: %w", err)
		}
	}

	flyoverScript = getFlyoverRedeemScript(derivationValue, fedRedeemScript)
	if addressScriptHash, err = btcutil.NewAddressScriptHash(flyoverScript, bridge.btcParams); err != nil {
		return blockchain.FlyoverDerivation{}, err
	}
	return blockchain.FlyoverDerivation{
		Address:      addressScriptHash.EncodeAddress(),
		RedeemScript: hex.EncodeToString(flyoverScript),
	}, nil
}

func (bridge *rskBridgeImpl) GetRequiredTxConfirmations() uint64 {
	return bridge.requiredConfirmations
}

func (bridge *rskBridgeImpl) FetchFederationInfo() (blockchain.FederationInfo, error) {
	var err error
	var pubKey []byte
	var pubKeys []string
	var i, federationSize int64

	opts := &bind.CallOpts{}
	fedSize, err := rskRetry(func() (*big.Int, error) {
		return bridge.contract.GetFederationSize(opts)
	})
	if err != nil {
		return blockchain.FederationInfo{}, err
	}
	federationSize = fedSize.Int64()

	for i = 0; i < federationSize; i++ {
		pubKey, err = rskRetry(func() ([]byte, error) {
			return bridge.contract.GetFederatorPublicKeyOfType(opts, big.NewInt(i), "btc")
		})
		if err != nil {
			return blockchain.FederationInfo{}, fmt.Errorf("error fetching fed public key: %w", err)
		}
		pubKeys = append(pubKeys, hex.EncodeToString(pubKey))
	}

	fedThreshold, err := rskRetry(func() (*big.Int, error) {
		return bridge.contract.GetFederationThreshold(opts)
	})
	if err != nil {
		return blockchain.FederationInfo{}, fmt.Errorf("error fetching federation size: %w", err)
	}

	fedAddress, err := rskRetry(func() (string, error) {
		return bridge.contract.GetFederationAddress(opts)
	})
	if err != nil {
		return blockchain.FederationInfo{}, fmt.Errorf("error fetching federation address: %w", err)
	}

	activeFedBlockHeight, err := rskRetry(func() (*big.Int, error) {
		return bridge.contract.GetActiveFederationCreationBlockHeight(opts)
	})
	if err != nil {
		return blockchain.FederationInfo{}, fmt.Errorf("error fetching federation height: %w", err)
	}

	return blockchain.FederationInfo{
		FedThreshold:         fedThreshold.Int64(),
		FedSize:              fedSize.Int64(),
		PubKeys:              pubKeys,
		FedAddress:           fedAddress,
		ActiveFedBlockHeight: activeFedBlockHeight.Int64(),
		IrisActivationHeight: bridge.irisActivationHeight,
		ErpKeys:              bridge.erpKeys,
	}, nil
}

func (bridge *rskBridgeImpl) getDerivationValueHash(args blockchain.FlyoverDerivationArgs) []byte {
	var buf bytes.Buffer
	buf.Write(args.QuoteHash)
	buf.Write(args.UserBtcRefundAddress)
	buf.Write(args.LbcAdress)
	buf.Write(args.LpBtcAddress)

	derivationValueHash := crypto.Keccak256(buf.Bytes())

	return derivationValueHash
}
