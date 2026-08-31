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
	repository *mocks.PegInWatchRepositoryMock
	registry   *mocks.PegInAddressRegistryContractMock
	wallet     *mocks.BitcoinWalletMock
	useCase    *watcher.DiscoverRegisteredAddressUseCase
}

func newDiscoverRegisteredAddressFixture(t *testing.T) *discoverRegisteredAddressFixture {
	t.Helper()
	repository := mocks.NewPegInWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	return &discoverRegisteredAddressFixture{
		repository: repository,
		registry:   registry,
		wallet:     wallet,
		useCase:    watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet),
	}
}

func registryDepositPayload(index int) ([]byte, string) {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[index]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return payload, decoded.Address
}

func TestDiscoverRegisteredAddressUseCase_Run_ImportsAddressAndRequestsRescan(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	payload, address := registryDepositPayload(0)
	event := blockchain.AddressRegistered{
		TxHash:      "0xreg",
		LogIndex:    3,
		BlockNumber: 101,
		RskAddress:  "0xrsk",
		Registrant:  "0xregistrant",
	}
	discovered := rootstock.PegInWatch{
		TxHash:      event.TxHash,
		LogIndex:    event.LogIndex,
		BlockNumber: event.BlockNumber,
		RskAddress:  event.RskAddress,
		Registrant:  event.Registrant,
		State:       rootstock.PegInWatchDiscovered,
	}

	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(nil, nil).Once()
	fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.MatchedBy(func(watch rootstock.PegInWatch) bool {
		return watch.RskAddress == event.RskAddress &&
			watch.TxHash == event.TxHash &&
			watch.LogIndex == event.LogIndex &&
			watch.State == rootstock.PegInWatchDiscovered
	})).Return(nil).Once()
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(address).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	require.NotNil(t, result.Watch)
	assert.True(t, result.NeedsRescan)
	assert.Equal(t, address, result.Watch.BtcAddress)
	assert.Equal(t, rootstock.PegInWatchDiscovered, result.Watch.State)
	assert.Equal(t, event.TxHash, result.Watch.TxHash)
	assert.Equal(t, event.RegistrationRoot, result.Watch.RegistrationRoot)
}

func TestDiscoverRegisteredAddressUseCase_Run_AlreadyImportedIsNoOp(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "0xoriginal", LogIndex: 1, RskAddress: "0xrsk"}
	imported := rootstock.PegInWatch{
		TxHash:     event.TxHash,
		LogIndex:   event.LogIndex,
		RskAddress: event.RskAddress,
		State:      rootstock.PegInWatchImported,
		BtcAddress: "n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7",
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&imported, nil).Once()

	result, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, result.NeedsRescan)
	assert.Equal(t, imported.TxHash, result.Watch.TxHash)
	assert.Equal(t, rootstock.PegInWatchImported, result.Watch.State)
	fixture.registry.AssertNotCalled(t, "GetPegInAddress", mock.Anything)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_PersistsUnsupportedEncoding(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "0xbech32", LogIndex: 1, RskAddress: "0xbech32"}
	discovered := rootstock.PegInWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInWatchDiscovered,
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Encoding: blockchain.PegInAddressRegistryEncodingBech32,
	}, nil).Once()
	fixture.repository.On("Update", test.AnyCtx, mock.MatchedBy(func(watch rootstock.PegInWatch) bool {
		return watch.TxHash == event.TxHash &&
			watch.State == rootstock.PegInWatchUnsupportedEncoding &&
			watch.BtcAddress == "" &&
			watch.LastError != ""
	})).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, result.NeedsRescan)
	assert.Equal(t, rootstock.PegInWatchUnsupportedEncoding, result.Watch.State)
	assert.NotEmpty(t, result.Watch.LastError)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_RecordsUnencodablePayload(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "truncated-payload", LogIndex: 1, RskAddress: "truncated-rsk"}
	discovered := rootstock.PegInWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInWatchDiscovered,
	}
	payload, _ := registryDepositPayload(0)
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload[:20],
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.repository.On("Update", test.AnyCtx, mock.MatchedBy(func(watch rootstock.PegInWatch) bool {
		return watch.TxHash == event.TxHash &&
			watch.State == rootstock.PegInWatchDiscovered &&
			strings.Contains(watch.LastError, "encode PegIn address for event truncated-payload/1")
	})).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, result.NeedsRescan)
	assert.NotEmpty(t, result.Watch.LastError)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_WrapsPersistErrors(t *testing.T) {
	event := blockchain.AddressRegistered{TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk"}
	discoveredWatch := func() *rootstock.PegInWatch {
		return &rootstock.PegInWatch{
			TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
			State: rootstock.PegInWatchDiscovered,
		}
	}

	t.Run("load", func(t *testing.T) {
		fixture := newDiscoverRegisteredAddressFixture(t)
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(nil, assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background(), event)
		require.Error(t, err)
		assert.ErrorContains(t, err, string(usecases.DiscoverRegisteredAddressId))
	})

	t.Run("upsert", func(t *testing.T) {
		fixture := newDiscoverRegisteredAddressFixture(t)
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(nil, nil).Once()
		fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background(), event)
		require.Error(t, err)
		assert.ErrorContains(t, err, "persist AddressRegistered")
	})

	t.Run("reload after upsert", func(t *testing.T) {
		fixture := newDiscoverRegisteredAddressFixture(t)
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(nil, nil).Once()
		fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.Anything).Return(nil).Once()
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(nil, assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background(), event)
		require.Error(t, err)
		assert.ErrorContains(t, err, "load AddressRegistered")
	})

	t.Run("missing after upsert", func(t *testing.T) {
		fixture := newDiscoverRegisteredAddressFixture(t)
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(nil, nil).Once()
		fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.Anything).Return(nil).Once()
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(nil, nil).Once()

		_, err := fixture.useCase.Run(context.Background(), event)
		require.Error(t, err)
		assert.ErrorContains(t, err, "watch not found after upsert")
	})

	t.Run("unsupported encoding persist", func(t *testing.T) {
		fixture := newDiscoverRegisteredAddressFixture(t)
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(discoveredWatch(), nil).Once()
		fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
			Encoding: blockchain.PegInAddressRegistryEncodingBech32,
		}, nil).Once()
		fixture.repository.EXPECT().Update(test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background(), event)
		require.Error(t, err)
		assert.ErrorContains(t, err, "persist unsupported encoding")
	})

	t.Run("watch error persist", func(t *testing.T) {
		fixture := newDiscoverRegisteredAddressFixture(t)
		fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(discoveredWatch(), nil).Once()
		fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{}, assert.AnError).Once()
		fixture.repository.EXPECT().Update(test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

		_, err := fixture.useCase.Run(context.Background(), event)
		require.Error(t, err)
		assert.ErrorContains(t, err, "persist PegIn address registry watch error")
	})
}

func TestDiscoverRegisteredAddressUseCase_Run_RecordsResolveErrorWithoutFailingRun(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk"}
	discovered := rootstock.PegInWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInWatchDiscovered,
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{}, assert.AnError).Once()
	fixture.repository.EXPECT().Update(test.AnyCtx, mock.Anything).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background(), event)
	require.NoError(t, err)
	assert.False(t, result.NeedsRescan)
	assert.NotEmpty(t, result.Watch.LastError)
}

func TestDiscoverRegisteredAddressUseCase_Run_SuppressesDuplicateWatchError(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	event := blockchain.AddressRegistered{TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk"}
	resolveErr := errors.New("resolve PegIn address for event 0xreg/1: boom")
	discovered := rootstock.PegInWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State:     rootstock.PegInWatchDiscovered,
		LastError: resolveErr.Error(),
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{}, errors.New("boom")).Once()

	result, err := fixture.useCase.Run(context.Background(), event)
	require.NoError(t, err)
	assert.False(t, result.NeedsRescan)
	fixture.repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_WalletAlreadyImportedStillNeedsRescan(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	payload, address := registryDepositPayload(0)
	event := blockchain.AddressRegistered{TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk"}
	discovered := rootstock.PegInWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInWatchDiscovered,
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(address).Return(errors.New("address already imported")).Once()

	result, err := fixture.useCase.Run(context.Background(), event)
	require.NoError(t, err)
	assert.True(t, result.NeedsRescan)
	assert.Equal(t, address, result.Watch.BtcAddress)
	assert.Equal(t, rootstock.PegInWatchDiscovered, result.Watch.State)
}

func TestDiscoverRegisteredAddressUseCase_Run_RecordsUnexpectedImportError(t *testing.T) {
	fixture := newDiscoverRegisteredAddressFixture(t)
	payload, address := registryDepositPayload(0)
	event := blockchain.AddressRegistered{TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk"}
	discovered := rootstock.PegInWatch{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInWatchDiscovered,
	}
	fixture.repository.EXPECT().Get(test.AnyCtx, event.TxHash, event.LogIndex).Return(&discovered, nil).Once()
	fixture.registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	fixture.wallet.EXPECT().ImportAddress(address).Return(errors.New("wallet locked")).Once()
	fixture.repository.EXPECT().Update(test.AnyCtx, mock.Anything).Return(nil).Once()

	result, err := fixture.useCase.Run(context.Background(), event)
	require.NoError(t, err)
	assert.False(t, result.NeedsRescan)
	assert.Contains(t, result.Watch.LastError, "import PegIn address")
}
