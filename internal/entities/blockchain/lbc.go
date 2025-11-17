package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

const (
	GetPegoutDepositsErrorTemplate = "error executing getting deposits in range [%d, %d]"
)

type RegisterPeginParams struct {
	QuoteSignature        []byte
	BitcoinRawTransaction []byte
	PartialMerkleTree     []byte
	BlockHeight           *big.Int
	Quote                 quote.PeginQuote
}

func (params RegisterPeginParams) String() string {
	return fmt.Sprintf(
		"RegisterPeginParams { QuoteSignature: %s, BitcoinRawTransaction: %s, "+
			"PartialMerkleTree: %s, BlockHeight: %v, Quote: %+v }",
		hex.EncodeToString(params.QuoteSignature),
		hex.EncodeToString(params.BitcoinRawTransaction),
		hex.EncodeToString(params.PartialMerkleTree),
		params.BlockHeight,
		params.Quote,
	)
}

type RefundPegoutParams struct {
	QuoteHash          [32]byte
	BtcRawTx           []byte
	BtcBlockHeaderHash [32]byte
	MerkleBranchPath   *big.Int
	MerkleBranchHashes [][32]byte
}

func (params RefundPegoutParams) String() string {
	return fmt.Sprintf(
		"RefundPegoutParams { QuoteHash: %s, BtcRawTx: %s, "+
			"BtcBlockHeaderHash: %s, MerkleBranchPath: %v, MerkleBranchHashes: %v }",
		hex.EncodeToString(params.QuoteHash[:]),
		hex.EncodeToString(params.BtcRawTx),
		hex.EncodeToString(params.BtcBlockHeaderHash[:]),
		params.MerkleBranchPath,
		params.MerkleBranchHashes,
	)
}

type ProviderRegistrationParams struct {
	Name       string                          `validate:"required"`
	ApiBaseUrl string                          `validate:"required"`
	Status     bool                            `validate:"required"`
	Type       liquidity_provider.ProviderType `validate:"oneof=0 1 2"`
}

func NewProviderRegistrationParams(
	name string,
	apiBaseUrl string,
	status bool,
	providerType liquidity_provider.ProviderType,
) ProviderRegistrationParams {
	return ProviderRegistrationParams{
		Name:       name,
		ApiBaseUrl: apiBaseUrl,
		Status:     status,
		Type:       providerType,
	}
}

type PeginContract interface {
	GetAddress() string
	GetBalance(address string) (*entities.Wei, error)
	HashPeginQuote(peginQuote quote.PeginQuote) (string, error)
	CallForUser(txConfig TransactionConfig, peginQuote quote.PeginQuote) (TransactionReceipt, error)
	RegisterPegin(params RegisterPeginParams) (TransactionReceipt, error)
	DaoFeePercentage() (uint64, error)
}

type PegoutContract interface {
	GetAddress() string
	HashPegoutQuote(pegoutQuote quote.PegoutQuote) (string, error)
	RefundUserPegOut(quoteHash string) (string, error)
	IsPegOutQuoteCompleted(quoteHash string) (bool, error)
	GetDepositEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]quote.PegoutDeposit, error)
	RefundPegout(txConfig TransactionConfig, params RefundPegoutParams) (TransactionReceipt, error)
	DaoFeePercentage() (uint64, error)
}

type DiscoveryContract interface {
	GetAddress() string
	SetProviderStatus(id uint64, newStatus bool) error
	UpdateProvider(name, url string) (string, error)
	RegisterProvider(txConfig TransactionConfig, params ProviderRegistrationParams) (int64, error)
	GetProviders() ([]liquidity_provider.RegisteredLiquidityProvider, error)
	GetProvider(address string) (liquidity_provider.RegisteredLiquidityProvider, error)
	IsOperational(providerType liquidity_provider.ProviderType, address string) (bool, error)
}

type CollateralManagementContract interface {
	GetAddress() string
	GetPenalizedEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]penalization.PenalizedEvent, error)
	ProviderResign() error
	WithdrawCollateral() error
	AddCollateral(amount *entities.Wei) error
	AddPegoutCollateral(amount *entities.Wei) error
	GetCollateral(address string) (*entities.Wei, error)
	GetPegoutCollateral(address string) (*entities.Wei, error)
	GetMinimumCollateral() (*entities.Wei, error)
}
