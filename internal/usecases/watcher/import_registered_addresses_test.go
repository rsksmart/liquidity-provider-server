package watcher_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

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

type importRegisteredAddressesFixture struct {
	repository *mocks.PegInWatchRepositoryMock
	registry   *mocks.PegInAddressRegistryContractMock
	btcNetwork *mocks.BitcoinNetworkMock
	wallet     *mocks.BitcoinWalletMock
	useCase    *watcher.ImportRegisteredAddressesUseCase
}

func newImportRegisteredAddressesFixture(t *testing.T) *importRegisteredAddressesFixture {
	t.Helper()
	fixture := &importRegisteredAddressesFixture{
		repository: mocks.NewPegInWatchRepositoryMock(t),
		registry:   mocks.NewPegInAddressRegistryContractMock(t),
		btcNetwork: mocks.NewBitcoinNetworkMock(t),
		wallet:     mocks.NewBitcoinWalletMock(t),
	}
	fixture.useCase = watcher.NewImportRegisteredAddressesUseCase(
		fixture.repository,
		fixture.registry,
		fixture.btcNetwork,
		fixture.wallet,
	)
	return fixture
}

func registryDepositPayload(index int) ([]byte, string) {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[index]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return payload, decoded.Address
}

func discoveredRegistration(suffix int) rootstock.PegInWatch {
	return rootstock.PegInWatch{
		TxHash:           fmt.Sprintf("0xreg%d", suffix),
		LogIndex:         uint(suffix),
		BlockNumber:      100 + uint64(suffix),
		RskAddress:       fmt.Sprintf("0xrsk%d", suffix),
		Registrant:       "0xregistrant",
		RegistrationRoot: [32]byte{byte(suffix)},
		State:            rootstock.PegInWatchDiscovered,
	}
}

func resolvedRegistration(watch rootstock.PegInWatch, address string) rootstock.PegInWatch {
	watch.BtcAddress = address
	watch.Encoding = uint8(blockchain.PegInAddressRegistryEncodingBase58)
	return watch
}

func importedRegistration(watch rootstock.PegInWatch, address string) rootstock.PegInWatch {
	watch = resolvedRegistration(watch, address)
	watch.State = rootstock.PegInWatchImported
	watch.LastError = ""
	return watch
}

func failedRegistration(watch rootstock.PegInWatch, lastError string) rootstock.PegInWatch {
	watch.LastError = lastError
	return watch
}

func withoutUpdateTime(watch rootstock.PegInWatch) rootstock.PegInWatch {
	watch.UpdatedAt = time.Time{}
	return watch
}

func (fixture *importRegisteredAddressesFixture) withStored(watch rootstock.PegInWatch) {
	stored := watch
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(&stored, nil).Once()
}

func (fixture *importRegisteredAddressesFixture) withDepositAddress(rskAddress string, index int) string {
	payload, address := registryDepositPayload(index)
	fixture.registry.EXPECT().GetPegInAddress(rskAddress).Return(blockchain.PegInAddress{
		Payload:  payload,
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	return address
}

func (fixture *importRegisteredAddressesFixture) withImportedAddress(rskAddress string, index int) string {
	address := fixture.withDepositAddress(rskAddress, index)
	fixture.wallet.EXPECT().ImportAddress(address).Return(nil).Once()
	return address
}

func (fixture *importRegisteredAddressesFixture) withRescan(tip int64, fromHeight int64) {
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(tip), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(fromHeight).
		Return(blockchain.BitcoinRescanResult{StartHeight: fromHeight, StopHeight: tip}, nil).
		Once()
}

func (fixture *importRegisteredAddressesFixture) expectUpdate(expected rootstock.PegInWatch, err error) {
	fixture.repository.EXPECT().Update(test.AnyCtx, mock.MatchedBy(func(actual rootstock.PegInWatch) bool {
		return assert.ObjectsAreEqual(withoutUpdateTime(expected), withoutUpdateTime(actual))
	})).Return(err).Once()
}

func assertImportResult(t *testing.T, expected []rootstock.PegInWatch, result watcher.ImportRegisteredAddressesResult) {
	t.Helper()
	actual := make([]rootstock.PegInWatch, 0, len(result.Watches))
	for _, watch := range result.Watches {
		actual = append(actual, withoutUpdateTime(watch))
	}
	expectedWithoutTime := make([]rootstock.PegInWatch, 0, len(expected))
	for _, watch := range expected {
		expectedWithoutTime = append(expectedWithoutTime, withoutUpdateTime(watch))
	}
	assert.Equal(t, expectedWithoutTime, actual)
}

func TestImportRegisteredAddressesUseCase_Run_ImportsDiscoveredAddressesWithOneRescan(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	first := failedRegistration(discoveredRegistration(1), "previous failure")
	second := discoveredRegistration(2)
	fixture.withStored(first)
	fixture.withStored(second)
	firstAddress := fixture.withImportedAddress(first.RskAddress, 0)
	secondAddress := fixture.withImportedAddress(second.RskAddress, 1)
	fixture.withRescan(200, 100)
	fixture.expectUpdate(importedRegistration(first, firstAddress), nil)
	fixture.expectUpdate(importedRegistration(second, secondAddress), nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{first, second})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{
		importedRegistration(first, firstAddress),
		importedRegistration(second, secondAddress),
	}, result)
}

func TestImportRegisteredAddressesUseCase_Run_CreatesTheStoredWatchWhenMissing(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(nil, nil).Once()
	fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.MatchedBy(func(created rootstock.PegInWatch) bool {
		created.CreatedAt, created.UpdatedAt = time.Time{}, time.Time{}
		return assert.ObjectsAreEqual(watch, created)
	})).Return(nil).Once()
	fixture.withStored(watch)
	address := fixture.withImportedAddress(watch.RskAddress, 0)
	fixture.withRescan(200, 100)
	fixture.expectUpdate(importedRegistration(watch, address), nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{importedRegistration(watch, address)}, result)
}

func TestImportRegisteredAddressesUseCase_Run_SkipsWatchesThatAreNotDiscovered(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	imported := discoveredRegistration(1)
	imported.State = rootstock.PegInWatchImported
	unsupported := discoveredRegistration(2)
	unsupported.State = rootstock.PegInWatchUnsupportedEncoding

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{imported, unsupported})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{}, result)
	fixture.repository.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
	fixture.btcNetwork.AssertNotCalled(t, "GetHeight")
}

func TestImportRegisteredAddressesUseCase_Run_ReturnsTheStoredWatchWhenAlreadyImported(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	_, address := registryDepositPayload(0)
	stored := importedRegistration(watch, address)
	fixture.withStored(stored)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{stored}, result)
	fixture.registry.AssertNotCalled(t, "GetPegInAddress", mock.Anything)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.btcNetwork.AssertNotCalled(t, "GetHeight")
}

func TestImportRegisteredAddressesUseCase_Run_ImportsTheStoredRowOfARepeatedRskAddress(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	stored := discoveredRegistration(1)
	replay := stored
	replay.TxHash, replay.LogIndex = "0xreplay", 9
	fixture.withStored(stored)
	address := fixture.withImportedAddress(stored.RskAddress, 0)
	fixture.withRescan(200, 100)
	fixture.expectUpdate(importedRegistration(stored, address), nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{replay})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{importedRegistration(stored, address)}, result)
	fixture.repository.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestImportRegisteredAddressesUseCase_Run_PersistsUnsupportedEncoding(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	fixture.registry.EXPECT().GetPegInAddress(watch.RskAddress).Return(blockchain.PegInAddress{
		Encoding: blockchain.PegInAddressRegistryEncodingBech32,
	}, nil).Once()
	unsupported := watch
	unsupported.Encoding = uint8(blockchain.PegInAddressRegistryEncodingBech32)
	unsupported.State = rootstock.PegInWatchUnsupportedEncoding
	unsupported.LastError = fmt.Sprintf("unsupported encoding %d", blockchain.PegInAddressRegistryEncodingBech32)
	fixture.expectUpdate(unsupported, nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{unsupported}, result)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.btcNetwork.AssertNotCalled(t, "GetHeight")
}

func TestImportRegisteredAddressesUseCase_Run_RecordsUnencodablePayload(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	payload, _ := registryDepositPayload(0)
	fixture.withStored(watch)
	fixture.registry.EXPECT().GetPegInAddress(watch.RskAddress).Return(blockchain.PegInAddress{
		Payload:  payload[:20],
		Encoding: blockchain.PegInAddressRegistryEncodingBase58,
	}, nil).Once()
	failed := failedRegistration(
		resolvedRegistration(watch, ""),
		"encode PegIn address for event 0xreg1/1: base58check payload must be 25 bytes, got 20",
	)
	fixture.expectUpdate(failed, nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{failed}, result)
	fixture.wallet.AssertNotCalled(t, "ImportAddress", mock.Anything)
	fixture.btcNetwork.AssertNotCalled(t, "GetHeight")
}

func TestImportRegisteredAddressesUseCase_Run_RecordsResolveErrorAndImportsTheNextAddress(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	broken := discoveredRegistration(1)
	valid := discoveredRegistration(2)
	fixture.withStored(broken)
	fixture.registry.EXPECT().GetPegInAddress(broken.RskAddress).Return(blockchain.PegInAddress{}, assert.AnError).Once()
	brokenAfterError := failedRegistration(broken, fmt.Sprintf("resolve PegIn address for event 0xreg1/1: %v", assert.AnError))
	fixture.expectUpdate(brokenAfterError, nil)
	fixture.withStored(valid)
	address := fixture.withImportedAddress(valid.RskAddress, 1)
	fixture.withRescan(200, 100)
	fixture.expectUpdate(importedRegistration(valid, address), nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{broken, valid})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{brokenAfterError, importedRegistration(valid, address)}, result)
}

func TestImportRegisteredAddressesUseCase_Run_SuppressesARepeatedWatchError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := failedRegistration(discoveredRegistration(1), "resolve PegIn address for event 0xreg1/1: boom")
	fixture.withStored(watch)
	fixture.registry.EXPECT().GetPegInAddress(watch.RskAddress).Return(blockchain.PegInAddress{}, errors.New("boom")).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{watch}, result)
	fixture.repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestImportRegisteredAddressesUseCase_Run_RescansWhenTheWalletAlreadyImportedTheAddress(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	address := fixture.withDepositAddress(watch.RskAddress, 0)
	fixture.wallet.EXPECT().ImportAddress(address).Return(errors.New("address already imported")).Once()
	fixture.withRescan(200, 100)
	fixture.expectUpdate(importedRegistration(watch, address), nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{importedRegistration(watch, address)}, result)
}

func TestImportRegisteredAddressesUseCase_Run_RecordsUnexpectedImportError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	address := fixture.withDepositAddress(watch.RskAddress, 0)
	fixture.wallet.EXPECT().ImportAddress(address).Return(errors.New("wallet locked")).Once()
	failed := failedRegistration(resolvedRegistration(watch, address), "import PegIn address for event 0xreg1/1: wallet locked")
	fixture.expectUpdate(failed, nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{failed}, result)
	fixture.btcNetwork.AssertNotCalled(t, "GetHeight")
}

func TestImportRegisteredAddressesUseCase_Run_RescansFromTheFirstBlockWhenTheTipIsShallow(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	address := fixture.withImportedAddress(watch.RskAddress, 0)
	fixture.withRescan(40, 0)
	fixture.expectUpdate(importedRegistration(watch, address), nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{importedRegistration(watch, address)}, result)
}

func TestImportRegisteredAddressesUseCase_Run_RecordsTheHeightErrorOnEveryPendingWatch(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	first := discoveredRegistration(1)
	second := discoveredRegistration(2)
	fixture.withStored(first)
	fixture.withStored(second)
	firstAddress := fixture.withImportedAddress(first.RskAddress, 0)
	secondAddress := fixture.withImportedAddress(second.RskAddress, 1)
	fixture.btcNetwork.EXPECT().GetHeight().Return(nil, assert.AnError).Once()
	lastError := fmt.Sprintf("get BTC height for registry rescan: %v", assert.AnError)
	firstFailed := failedRegistration(resolvedRegistration(first, firstAddress), lastError)
	secondFailed := failedRegistration(resolvedRegistration(second, secondAddress), lastError)
	fixture.expectUpdate(firstFailed, nil)
	fixture.expectUpdate(secondFailed, nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{first, second})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{firstFailed, secondFailed}, result)
	fixture.wallet.AssertNotCalled(t, "RescanBlockchain", mock.Anything)
}

func TestImportRegisteredAddressesUseCase_Run_RecordsTheRescanErrorOnEveryPendingWatch(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	first := discoveredRegistration(1)
	second := discoveredRegistration(2)
	fixture.withStored(first)
	fixture.withStored(second)
	firstAddress := fixture.withImportedAddress(first.RskAddress, 0)
	secondAddress := fixture.withImportedAddress(second.RskAddress, 1)
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).Return(blockchain.BitcoinRescanResult{}, assert.AnError).Once()
	lastError := fmt.Sprintf("rescan PegIn addresses: %v", assert.AnError)
	firstFailed := failedRegistration(resolvedRegistration(first, firstAddress), lastError)
	secondFailed := failedRegistration(resolvedRegistration(second, secondAddress), lastError)
	fixture.expectUpdate(firstFailed, nil)
	fixture.expectUpdate(secondFailed, nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{first, second})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{firstFailed, secondFailed}, result)
}

func TestImportRegisteredAddressesUseCase_Run_SuppressesARepeatedRescanError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := failedRegistration(discoveredRegistration(1), fmt.Sprintf("rescan PegIn addresses: %v", assert.AnError))
	fixture.withStored(watch)
	address := fixture.withImportedAddress(watch.RskAddress, 0)
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).Return(blockchain.BitcoinRescanResult{}, assert.AnError).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{resolvedRegistration(watch, address)}, result)
	fixture.repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestImportRegisteredAddressesUseCase_Run_WrapsLoadError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(nil, assert.AnError).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ImportRegisteredAddressesId)+": load AddressRegistered event 0xrsk1")
	assertImportResult(t, []rootstock.PegInWatch{}, result)
}

func TestImportRegisteredAddressesUseCase_Run_WrapsUpsertError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(nil, nil).Once()
	fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, "persist AddressRegistered event 0xreg1/1")
	assertImportResult(t, []rootstock.PegInWatch{}, result)
}

func TestImportRegisteredAddressesUseCase_Run_WrapsReloadAfterUpsertError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(nil, nil).Once()
	fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.Anything).Return(nil).Once()
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(nil, assert.AnError).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, "load AddressRegistered event 0xrsk1")
	assertImportResult(t, []rootstock.PegInWatch{}, result)
}

func TestImportRegisteredAddressesUseCase_Run_WrapsMissingWatchAfterUpsert(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(nil, nil).Once()
	fixture.repository.EXPECT().Upsert(test.AnyCtx, mock.Anything).Return(nil).Once()
	fixture.repository.EXPECT().Get(test.AnyCtx, watch.RskAddress).Return(nil, nil).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.EqualError(t, err, string(usecases.ImportRegisteredAddressesId)+": load AddressRegistered event 0xrsk1: watch not found after upsert")
	assertImportResult(t, []rootstock.PegInWatch{}, result)
}

func TestImportRegisteredAddressesUseCase_Run_WrapsUnsupportedEncodingPersistError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	fixture.registry.EXPECT().GetPegInAddress(watch.RskAddress).Return(blockchain.PegInAddress{
		Encoding: blockchain.PegInAddressRegistryEncodingBech32,
	}, nil).Once()
	fixture.repository.EXPECT().Update(test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, "persist unsupported encoding for event 0xreg1/1")
	assertImportResult(t, []rootstock.PegInWatch{}, result)
}

func TestImportRegisteredAddressesUseCase_Run_WrapsWatchErrorPersistError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	fixture.registry.EXPECT().GetPegInAddress(watch.RskAddress).Return(blockchain.PegInAddress{}, errors.New("boom")).Once()
	fixture.repository.EXPECT().Update(test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, "persist PegIn address registry watch error for 0xreg1/1")
	assertImportResult(t, []rootstock.PegInWatch{}, result)
}

func TestImportRegisteredAddressesUseCase_Run_StopsBeforeTheRescanWhenAStoreReadFails(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	first := discoveredRegistration(1)
	second := discoveredRegistration(2)
	fixture.withStored(first)
	address := fixture.withImportedAddress(first.RskAddress, 0)
	fixture.repository.EXPECT().Get(test.AnyCtx, second.RskAddress).Return(nil, assert.AnError).Once()

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{first, second})

	require.ErrorIs(t, err, assert.AnError)
	assertImportResult(t, []rootstock.PegInWatch{resolvedRegistration(first, address)}, result)
	fixture.btcNetwork.AssertNotCalled(t, "GetHeight")
	fixture.repository.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestImportRegisteredAddressesUseCase_Run_KeepsTheWatchDiscoveredWhenTheImportedStateIsNotSaved(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	address := fixture.withImportedAddress(watch.RskAddress, 0)
	fixture.withRescan(200, 100)
	fixture.expectUpdate(importedRegistration(watch, address), assert.AnError)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ImportRegisteredAddressesId)+": persist imported state for event 0xreg1/1")
	assertImportResult(t, []rootstock.PegInWatch{resolvedRegistration(watch, address)}, result)
}

func TestImportRegisteredAddressesUseCase_Run_KeepsTheUnsavedWatchesDiscoveredWhenMarkingImportedFails(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	first := discoveredRegistration(1)
	second := discoveredRegistration(2)
	fixture.withStored(first)
	fixture.withStored(second)
	firstAddress := fixture.withImportedAddress(first.RskAddress, 0)
	secondAddress := fixture.withImportedAddress(second.RskAddress, 1)
	fixture.withRescan(200, 100)
	fixture.expectUpdate(importedRegistration(first, firstAddress), nil)
	fixture.expectUpdate(importedRegistration(second, secondAddress), assert.AnError)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{first, second})

	require.ErrorIs(t, err, assert.AnError)
	assertImportResult(t, []rootstock.PegInWatch{
		importedRegistration(first, firstAddress),
		resolvedRegistration(second, secondAddress),
	}, result)
}

func TestImportRegisteredAddressesUseCase_Run_KeepsTheUnsavedWatchesUnchangedWhenSavingTheRescanErrorFails(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	first := discoveredRegistration(1)
	second := discoveredRegistration(2)
	fixture.withStored(first)
	fixture.withStored(second)
	firstAddress := fixture.withImportedAddress(first.RskAddress, 0)
	secondAddress := fixture.withImportedAddress(second.RskAddress, 1)
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).Return(blockchain.BitcoinRescanResult{}, errors.New("rescan aborted")).Once()
	lastError := "rescan PegIn addresses: rescan aborted"
	firstFailed := failedRegistration(resolvedRegistration(first, firstAddress), lastError)
	fixture.expectUpdate(firstFailed, nil)
	fixture.expectUpdate(failedRegistration(resolvedRegistration(second, secondAddress), lastError), assert.AnError)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{first, second})

	require.ErrorIs(t, err, assert.AnError)
	assertImportResult(t, []rootstock.PegInWatch{firstFailed, resolvedRegistration(second, secondAddress)}, result)
}

func TestImportRegisteredAddressesUseCase_Run_ReturnsTheHandledWatchesWithTheRescanFailures(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	handled := discoveredRegistration(1)
	pending := discoveredRegistration(2)
	_, handledAddress := registryDepositPayload(0)
	storedHandled := importedRegistration(handled, handledAddress)
	fixture.withStored(storedHandled)
	fixture.withStored(pending)
	pendingAddress := fixture.withImportedAddress(pending.RskAddress, 1)
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).Return(blockchain.BitcoinRescanResult{}, errors.New("rescan aborted")).Once()
	pendingFailed := failedRegistration(resolvedRegistration(pending, pendingAddress), "rescan PegIn addresses: rescan aborted")
	fixture.expectUpdate(pendingFailed, nil)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{handled, pending})

	require.NoError(t, err)
	assertImportResult(t, []rootstock.PegInWatch{storedHandled, pendingFailed}, result)
}

func TestImportRegisteredAddressesUseCase_Run_WrapsRescanErrorPersistError(t *testing.T) {
	fixture := newImportRegisteredAddressesFixture(t)
	watch := discoveredRegistration(1)
	fixture.withStored(watch)
	address := fixture.withImportedAddress(watch.RskAddress, 0)
	fixture.btcNetwork.EXPECT().GetHeight().Return(big.NewInt(200), nil).Once()
	fixture.wallet.EXPECT().RescanBlockchain(int64(100)).Return(blockchain.BitcoinRescanResult{}, errors.New("rescan aborted")).Once()
	failed := failedRegistration(resolvedRegistration(watch, address), "rescan PegIn addresses: rescan aborted")
	fixture.expectUpdate(failed, assert.AnError)

	result, err := fixture.useCase.Run(context.Background(), []rootstock.PegInWatch{watch})

	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, "persist PegIn address registry watch error for 0xreg1/1")
	assertImportResult(t, []rootstock.PegInWatch{resolvedRegistration(watch, address)}, result)
}
