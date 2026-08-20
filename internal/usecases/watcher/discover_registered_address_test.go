package watcher_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type discoverRegisteredAddressFixture struct {
	repository *mocks.PegInAddressRegistryWatchRepositoryMock
	registry   *mocks.PegInAddressRegistryContractMock
	wallet     *mocks.BitcoinWalletMock
	useCase    *watcher.DiscoverRegisteredAddressUseCase
}

func newDiscoverRegisteredAddressFixture(t *testing.T) *discoverRegisteredAddressFixture {
	t.Helper()
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	return &discoverRegisteredAddressFixture{
		repository: repository,
		registry:   registry,
		wallet:     wallet,
		useCase:    watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet),
	}
}

func registryDepositPayload() ([]byte, string) {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[0]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return payload, decoded.Address
}

func TestDiscoverRegisteredAddressUseCase_Run_ImportsAddressAndRequestsRescan(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	payload, address := registryDepositPayload()
	event := blockchain.AddressRegistered{
		TxHash:      "0xreg",
		LogIndex:    3,
		BlockNumber: 101,
		RskAddress:  "0xrsk",
		Registrant:  "0xregistrant",
	}
	discovered := rootstock.PegInAddressRegistryWatch{
		TxHash:      event.TxHash,
		LogIndex:    event.LogIndex,
		BlockNumber: event.BlockNumber,
		RskAddress:  event.RskAddress,
		Registrant:  event.Registrant,
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
	}

	fixture.repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(nil, nil).Once()
	fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.MatchedBy(func(watch rootstock.PegInAddressRegistryWatch) bool {
		return watch.RskAddress == event.RskAddress &&
			watch.TxHash == event.TxHash &&
			watch.LogIndex == event.LogIndex &&
			watch.State == rootstock.PegInAddressRegistryWatchDiscovered
	})).Return(nil).Once()
	fixture.repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(address).Return(nil).Once()

	watch, needsRescan, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	require.NotNil(t, watch)
	assert.True(t, needsRescan)
	assert.Equal(t, address, watch.BtcAddress)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchDiscovered, watch.State)
	assert.Equal(t, event.TxHash, watch.TxHash)
	assert.Equal(t, event.RegistrationRoot, watch.RegistrationRoot)
}

func TestDiscoverRegisteredAddressUseCase_Run_AlreadyImportedIsNoOp(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "0xoriginal", LogIndex: 1, RskAddress: "0xrsk"}
	imported := rootstock.PegInAddressRegistryWatch{
		TxHash:     event.TxHash,
		LogIndex:   event.LogIndex,
		RskAddress: event.RskAddress,
		State:      rootstock.PegInAddressRegistryWatchImported,
		BtcAddress: "n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7",
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&imported, nil).Once()

	watch, needsRescan, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, needsRescan)
	assert.Equal(t, imported.TxHash, watch.TxHash)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, watch.State)
	fixture.registry.AssertNotCalled(t, "GetPegInAddress", mock.Anything)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_DuplicateRskAddressKeepsFirstRow(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	payload, address := registryDepositPayload()
	first := blockchain.AddressRegistered{TxHash: "0xfirst", LogIndex: 1, RskAddress: "0xshared"}
	replay := blockchain.AddressRegistered{TxHash: "0xreplay", LogIndex: 2, RskAddress: "0xshared"}
	discovered := rootstock.PegInAddressRegistryWatch{
		TxHash:     first.TxHash,
		LogIndex:   first.LogIndex,
		RskAddress: first.RskAddress,
		State:      rootstock.PegInAddressRegistryWatchDiscovered,
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, first.RskAddress).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(first.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(address).Return(nil).Once()

	watch, needsRescan, err := fixture.useCase.Run(context.Background(), replay)

	require.NoError(t, err)
	assert.True(t, needsRescan)
	assert.Equal(t, first.TxHash, watch.TxHash)
	assert.Equal(t, first.LogIndex, watch.LogIndex)
	fixture.repository.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_PersistsUnsupportedEncoding(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "0xbech32", LogIndex: 1, RskAddress: "0xbech32"}
	discovered := rootstock.PegInAddressRegistryWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Encoding: blockchain.PegInAddressRegistryEncodingBech32,
	}, nil).Once()
	fixture.repository.On("Update", test.AnyCtx, mock.MatchedBy(func(watch rootstock.PegInAddressRegistryWatch) bool {
		return watch.TxHash == event.TxHash &&
			watch.State == rootstock.PegInAddressRegistryWatchUnsupportedEncoding &&
			watch.BtcAddress == ""
	})).Return(nil).Once()

	watch, needsRescan, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, needsRescan)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchUnsupportedEncoding, watch.State)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_RecordsUnencodablePayload(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "truncated-payload", LogIndex: 1, RskAddress: "truncated-rsk"}
	discovered := rootstock.PegInAddressRegistryWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	payload, _ := registryDepositPayload()
	fixture.repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload[:20],
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.repository.On("Update", test.AnyCtx, mock.MatchedBy(func(watch rootstock.PegInAddressRegistryWatch) bool {
		return watch.TxHash == event.TxHash &&
			watch.State == rootstock.PegInAddressRegistryWatchDiscovered &&
			strings.Contains(watch.LastError, "encode PegIn address for event truncated-payload/1")
	})).Return(nil).Once()

	watch, needsRescan, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, needsRescan)
	assert.NotEmpty(t, watch.LastError)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_WrapsPersistErrors(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk"}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(nil, assert.AnError).Once()

	_, _, err := fixture.useCase.Run(context.Background(), event)
	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.DiscoverRegisteredAddressId))
}

func TestDiscoverRegisteredAddressUseCase_Run_WalletAlreadyImportedStillNeedsRescan(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	payload, address := registryDepositPayload()
	event := blockchain.AddressRegistered{TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk"}
	discovered := rootstock.PegInAddressRegistryWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(address).Return(errors.New("address already imported")).Once()

	watch, needsRescan, err := fixture.useCase.Run(context.Background(), event)
	require.NoError(t, err)
	assert.True(t, needsRescan)
	assert.Equal(t, address, watch.BtcAddress)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchDiscovered, watch.State)
}
