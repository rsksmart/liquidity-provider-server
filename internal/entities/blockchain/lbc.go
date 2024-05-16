package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"math/big"
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
	Type       liquidity_provider.ProviderType `validate:"required"`
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

type LiquidityBridgeContract interface {
	GetAddress() string
	HashPeginQuote(peginQuote quote.PeginQuote) (string, error)
	HashPegoutQuote(pegoutQuote quote.PegoutQuote) (string, error)
	GetProviders() ([]liquidity_provider.RegisteredLiquidityProvider, error)
	ProviderResign() error
	SetProviderStatus(id uint64, newStatus bool) error
	GetCollateral(address string) (*entities.Wei, error)
	GetPegoutCollateral(address string) (*entities.Wei, error)
	GetMinimumCollateral() (*entities.Wei, error)
	AddCollateral(amount *entities.Wei) error
	AddPegoutCollateral(amount *entities.Wei) error
	WithdrawCollateral() error
	GetBalance(address string) (*entities.Wei, error)
	CallForUser(txConfig TransactionConfig, peginQuote quote.PeginQuote) (string, error)
	RegisterPegin(params RegisterPeginParams) (string, error)
	RefundPegout(txConfig TransactionConfig, params RefundPegoutParams) (string, error)
	IsOperationalPegin(address string) (bool, error)
	IsOperationalPegout(address string) (bool, error)
	RegisterProvider(txConfig TransactionConfig, params ProviderRegistrationParams) (int64, error)
	GetDepositEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]quote.PegoutDeposit, error)
	GetPeginPunishmentEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]liquidity_provider.PunishmentEvent, error)
}

type FeeCollector interface {
	DaoFeePercentage() (uint64, error)
}
