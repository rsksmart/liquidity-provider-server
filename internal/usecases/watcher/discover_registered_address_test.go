package watcher_test

import (
	"context"
	"fmt"
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

func registryDepositPayload(index int) ([]byte, string) {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[index]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return payload, decoded.Address
}

func TestDiscoverRegisteredAddressUseCase_Run_ImportsWithoutBindingAPayment(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	payload, address := registryDepositPayload(0)
	event := blockchain.AddressRegistered{
		TxHash:      "0xreg",
		LogIndex:    3,
		BlockNumber: 101,
		RskAddress:  "0xrsk",
		Registrant:  "0xregistrant",
	}
	discovered := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:      event.TxHash,
		LogIndex:    event.LogIndex,
		BlockNumber: event.BlockNumber,
		RskAddress:  event.RskAddress,
		Registrant:  event.Registrant,
		State:       rootstock.PegInAddressRegistryWatchDiscovered,
	}

	repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(nil, nil).Once()
	repository.EXPECT().Upsert(test.AnyCtx, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
		return entry.RskAddress == event.RskAddress &&
			entry.TxHash == event.TxHash &&
			entry.LogIndex == event.LogIndex &&
			entry.State == rootstock.PegInAddressRegistryWatchDiscovered
	})).Return(nil).Once()
	repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&discovered, nil).Once()
	registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	wallet.EXPECT().ImportAddress(address).Return(nil).Once()

	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet)
	entry, needsRescan, err := useCase.Run(context.Background(), event)

	require.NoError(t, err)
	require.NotNil(t, entry)
	assert.True(t, needsRescan)
	assert.Equal(t, address, entry.BtcAddress)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchDiscovered, entry.State)
	assert.Equal(t, event.TxHash, entry.TxHash)
	assert.Equal(t, event.RegistrationRoot, entry.RegistrationRoot)
}

func TestDiscoverRegisteredAddressUseCase_Run_AlreadyImportedIsNoOp(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	event := blockchain.AddressRegistered{TxHash: "0xnew", LogIndex: 9, RskAddress: "0xrsk"}
	imported := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:     "0xoriginal",
		LogIndex:   1,
		RskAddress: event.RskAddress,
		State:      rootstock.PegInAddressRegistryWatchImported,
		BtcAddress: "n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7",
	}
	repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&imported, nil).Once()

	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet)
	entry, needsRescan, err := useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, needsRescan)
	assert.Equal(t, imported.TxHash, entry.TxHash)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entry.State)
	registry.AssertNotCalled(t, "GetPegInAddress", mock.Anything)
	wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_DuplicateRskAddressKeepsFirstRow(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	payload, address := registryDepositPayload(0)
	first := blockchain.AddressRegistered{TxHash: "0xfirst", LogIndex: 1, RskAddress: "0xshared"}
	replay := blockchain.AddressRegistered{TxHash: "0xreplay", LogIndex: 2, RskAddress: "0xshared"}
	discovered := rootstock.PegInAddressRegistryWatchEntry{
		TxHash:     first.TxHash,
		LogIndex:   first.LogIndex,
		RskAddress: first.RskAddress,
		State:      rootstock.PegInAddressRegistryWatchDiscovered,
	}
	repository.EXPECT().Get(test.AnyCtx, first.RskAddress).Return(&discovered, nil).Once()
	registry.EXPECT().GetPegInAddress(first.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	wallet.EXPECT().ImportAddress(address).Return(nil).Once()

	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet)
	entry, needsRescan, err := useCase.Run(context.Background(), replay)

	require.NoError(t, err)
	assert.True(t, needsRescan)
	assert.Equal(t, first.TxHash, entry.TxHash)
	assert.Equal(t, first.LogIndex, entry.LogIndex)
	repository.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_PersistsUnsupportedEncoding(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	event := blockchain.AddressRegistered{TxHash: "0xbech32", LogIndex: 1, RskAddress: "0xbech32"}
	discovered := rootstock.PegInAddressRegistryWatchEntry{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&discovered, nil).Once()
	registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Encoding: blockchain.PegInAddressRegistryEncodingBech32,
	}, nil).Once()
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
		return entry.RskAddress == event.RskAddress &&
			entry.State == rootstock.PegInAddressRegistryWatchUnsupportedEncoding &&
			entry.BtcAddress == ""
	})).Return(nil).Once()

	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet)
	entry, needsRescan, err := useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, needsRescan)
	assert.Equal(t, rootstock.PegInAddressRegistryWatchUnsupportedEncoding, entry.State)
	wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_RecordsUnencodablePayload(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	event := blockchain.AddressRegistered{TxHash: "truncated-payload", LogIndex: 1, RskAddress: "truncated-rsk"}
	discovered := rootstock.PegInAddressRegistryWatchEntry{
		TxHash: event.TxHash, LogIndex: event.LogIndex, RskAddress: event.RskAddress,
		State: rootstock.PegInAddressRegistryWatchDiscovered,
	}
	payload, _ := registryDepositPayload(0)
	repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(&discovered, nil).Once()
	registry.EXPECT().GetPegInAddress(event.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload[:20],
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(entry rootstock.PegInAddressRegistryWatchEntry) bool {
		return entry.RskAddress == event.RskAddress &&
			entry.State == rootstock.PegInAddressRegistryWatchDiscovered &&
			strings.Contains(entry.LastError, "encode PegIn address for event truncated-payload/1")
	})).Return(nil).Once()

	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet)
	entry, needsRescan, err := useCase.Run(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, needsRescan)
	assert.NotEmpty(t, entry.LastError)
	wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_MarkImported(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	entry := &rootstock.PegInAddressRegistryWatchEntry{
		TxHash: "0xreg", LogIndex: 1, RskAddress: "0xrsk",
		State:     rootstock.PegInAddressRegistryWatchDiscovered,
		LastError: "previous",
	}
	repository.On("Update", test.AnyCtx, mock.MatchedBy(func(updated rootstock.PegInAddressRegistryWatchEntry) bool {
		return updated.RskAddress == entry.RskAddress &&
			updated.State == rootstock.PegInAddressRegistryWatchImported &&
			updated.LastError == ""
	})).Return(nil).Once()

	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, nil, nil)
	require.NoError(t, useCase.MarkImported(context.Background(), entry))
	assert.Equal(t, rootstock.PegInAddressRegistryWatchImported, entry.State)
}

func TestDiscoverRegisteredAddressUseCase_RecordError_SuppressesIdenticalWrites(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	entry := &rootstock.PegInAddressRegistryWatchEntry{
		TxHash: "stuck", LogIndex: 1, RskAddress: "stuck-rsk",
		State:     rootstock.PegInAddressRegistryWatchDiscovered,
		LastError: fmt.Sprintf("import PegIn address for event stuck/1: %v", assert.AnError),
	}
	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, nil, nil)
	require.NoError(t, useCase.RecordError(context.Background(), entry, fmt.Errorf("import PegIn address for event stuck/1: %w", assert.AnError)))
	repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestDiscoverRegisteredAddressUseCase_Run_WrapsPersistErrors(t *testing.T) {
	repository := mocks.NewPegInAddressRegistryWatchRepositoryMock(t)
	registry := mocks.NewPegInAddressRegistryContractMock(t)
	wallet := mocks.NewBitcoinWalletMock(t)
	event := blockchain.AddressRegistered{TxHash: "0xreg", RskAddress: "0xrsk"}
	repository.EXPECT().Get(test.AnyCtx, event.RskAddress).Return(nil, assert.AnError).Once()

	useCase := watcher.NewDiscoverRegisteredAddressUseCase(repository, registry, wallet)
	_, _, err := useCase.Run(context.Background(), event)
	require.Error(t, err)
	assert.ErrorContains(t, err, string(usecases.DiscoverRegisteredAddressId))
}
