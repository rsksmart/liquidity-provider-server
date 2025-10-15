package rootstock

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"math/big"
)

type Bridge interface {
	GetAddress() string
	GetFedAddress() (string, error)
	GetMinimumLockTxValue() (*entities.Wei, error)
	GetFlyoverDerivationAddress(args FlyoverDerivationArgs) (FlyoverDerivation, error)
	GetRequiredTxConfirmations() uint64
	FetchFederationInfo() (FederationInfo, error)
	RegisterBtcCoinbaseTransaction(registrationParams BtcCoinbaseTransactionInformation) (string, error)
	GetBatchPegOutCreatedEvent(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]BatchPegOut, error)
}

type FlyoverDerivationArgs struct {
	FedInfo              FederationInfo
	LbcAdress            []byte
	UserBtcRefundAddress []byte
	LpBtcAddress         []byte
	QuoteHash            []byte
}

type FlyoverDerivation struct {
	Address      string
	RedeemScript string
}

type FederationInfo struct {
	FedSize              int64
	FedThreshold         int64
	PubKeys              []string
	FedAddress           string
	ActiveFedBlockHeight int64
	IrisActivationHeight int64
	ErpKeys              []string
	UseSegwit            bool
}

type BtcCoinbaseTransactionInformation struct {
	BtcTxSerialized      []byte
	BlockHash            [32]byte
	BlockHeight          *big.Int
	SerializedPmt        []byte
	WitnessMerkleRoot    [32]byte
	WitnessReservedValue [32]byte
}

func (params BtcCoinbaseTransactionInformation) String() string {
	return fmt.Sprintf(
		"RegisterPeginParams { BtcTxSerialized: %s, BlockHash: %s, BlockHeight: %d"+
			"SerializedPmt: %s, WitnessMerkleRoot: %s, WitnessReservedValue: %s }",
		hex.EncodeToString(params.BtcTxSerialized),
		hex.EncodeToString(params.BlockHash[:]),
		params.BlockHeight.Uint64(),
		hex.EncodeToString(params.SerializedPmt),
		hex.EncodeToString(params.WitnessMerkleRoot[:]),
		hex.EncodeToString(params.WitnessReservedValue[:]),
	)
}
