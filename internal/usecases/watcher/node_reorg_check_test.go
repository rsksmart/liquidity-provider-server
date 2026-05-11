package watcher

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func matchNodeReorgCheckEventExpected(expected blockchain.NodeReorgCheckEvent) func(entities.Event) bool {
	return func(event entities.Event) bool {
		reorgEvent, ok := event.(blockchain.NodeReorgCheckEvent)
		return ok &&
			reorgEvent.Id() == blockchain.NodeReorgCheckEventId &&
			!reorgEvent.CreationTimestamp().IsZero() &&
			reorgEvent.NodeType == expected.NodeType &&
			reorgEvent.CurrentDepth == expected.CurrentDepth &&
			reorgEvent.MaxAllowedDepth == expected.MaxAllowedDepth &&
			reorgEvent.AboveThreshold == expected.AboveThreshold
	}
}

func matchNodeReorgCheckErrorEvent(nodeType entities.NodeType) func(entities.Event) bool {
	return func(event entities.Event) bool {
		reorgErrorEvent, ok := event.(blockchain.NodeReorgCheckErrorEvent)
		return ok &&
			reorgErrorEvent.Id() == blockchain.NodeReorgCheckErrorEventId &&
			!reorgErrorEvent.CreationTimestamp().IsZero() &&
			reorgErrorEvent.NodeType == nodeType
	}
}

func matchNodeReorgAlertSentEvent(nodeType entities.NodeType, detectedDepth uint64) func(entities.Event) bool {
	return func(event entities.Event) bool {
		reorgAlertEvent, ok := event.(blockchain.NodeReorgAlertSentEvent)
		return ok &&
			reorgAlertEvent.Id() == blockchain.NodeReorgAlertSentEventId &&
			!reorgAlertEvent.CreationTimestamp().IsZero() &&
			reorgAlertEvent.NodeType == nodeType &&
			reorgAlertEvent.DetectedDepth == detectedDepth
	}
}

func TestReorgHistoryWindow(t *testing.T) {
	t.Parallel()
	assert.Equal(t, uint64(12), reorgHistoryWindow(0))
	assert.Equal(t, uint64(12), reorgHistoryWindow(2))
	assert.Equal(t, uint64(15), reorgHistoryWindow(5))
}

func TestNodeReorgCheckUseCase_Run_UnsupportedNodeType(t *testing.T) {
	t.Parallel()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{},
		&mocks.AlertSenderMock{},
		"",
		&mocks.EventBusMock{},
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeType("unknown"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

//nolint:funlen
func TestNodeReorgCheckUseCase_Run_Bitcoin(t *testing.T) {
	t.Parallel()

	t.Run("get blockchain info error", testNodeReorgCheckUseCaseRunBitcoinGetBlockchainInfoError)
	t.Run("first run", testNodeReorgCheckUseCaseRunBitcoinFirstRun)

	t.Run("same tip no extra rpc", testNodeReorgCheckUseCaseRunBitcoinSameTipNoExtraRPC)

	t.Run("one block append", testNodeReorgCheckUseCaseRunBitcoinOneBlockAppend)

	t.Run("reorg within history window", testNodeReorgCheckUseCaseRunBitcoinReorgWithinHistoryWindow)

	t.Run("advance beyond history window", testNodeReorgCheckUseCaseRunBitcoinAdvanceBeyondHistoryWindow)

	t.Run("first run refresh failure", testNodeReorgCheckUseCaseRunBitcoinFirstRunRefreshFailure)

	t.Run("one block append refresh failure", testNodeReorgCheckUseCaseRunBitcoinOneBlockAppendRefreshFailure)

	t.Run("divergence scan rpc failure", testNodeReorgCheckUseCaseRunBitcoinDivergenceScanRPCFailure)

	t.Run("advance beyond history window refresh failure", testNodeReorgCheckUseCaseRunBitcoinAdvanceBeyondHistoryWindowRefreshFailure)

	t.Run("above threshold reorg sends alert", testNodeReorgCheckUseCaseRunBitcoinAboveThresholdReorgSendsAlert)

	t.Run("get block header verbose error", testNodeReorgCheckUseCaseRunBitcoinGetBlockHeaderVerboseError)
}

func testNodeReorgCheckUseCaseRunBitcoinGetBlockchainInfoError(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	btc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeBitcoin))).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeBitcoin)
	require.ErrorIs(t, err, assert.AnError)
	btc.AssertNotCalled(t, "GetBlockHashAtHeight", mock.Anything)
	btc.AssertNotCalled(t, "GetBlockHeaderVerbose", mock.Anything)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinFirstRun(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_tip",
		ValidatedBlocks: new(big.Int).SetUint64(tipHeight),
	}
	btc.On("GetBlockchainInfo").Return(info, nil).Once()
	historyWindow := reorgHistoryWindow(2)
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := tipHeight - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeBitcoin)
	require.NoError(t, err)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinSameTipNoExtraRPC(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_tip",
		ValidatedBlocks: new(big.Int).SetUint64(tipHeight),
	}
	btc.On("GetBlockchainInfo").Return(info, nil).Twice()
	historyWindow := reorgHistoryWindow(2)
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := tipHeight - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Twice()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)
	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	btc.AssertNotCalled(t, "GetBlockHeaderVerbose", mock.Anything)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinOneBlockAppend(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_101",
		ValidatedBlocks: big.NewInt(101),
	}
	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	historyWindow := reorgHistoryWindow(2)
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_101").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_100",
	}, nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(101) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Twice()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)
	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

//nolint:funlen
func testNodeReorgCheckUseCaseRunBitcoinReorgWithinHistoryWindow(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_101_fork",
		ValidatedBlocks: big.NewInt(101),
	}

	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	historyWindow := reorgHistoryWindow(2)
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}

	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_101_fork").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_100_fork",
	}, nil).Once()
	btc.On("GetBlockHashAtHeight", int64(100)).Return("hash_100_fork", nil).Once()
	btc.On("GetBlockHashAtHeight", int64(99)).Return("hash_99", nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(101) - i
		hash := fmt.Sprintf("hash_%d", blockHeight)
		switch blockHeight {
		case 101:
			hash = "hash_101_fork"
		case 100:
			hash = "hash_100_fork"
		default:
		}
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(hash, nil).Once()
	}

	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    1,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinAdvanceBeyondHistoryWindow(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_113",
		ValidatedBlocks: big.NewInt(113),
	}

	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	historyWindow := reorgHistoryWindow(2)
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}

	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_113").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_112",
	}, nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(113) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}

	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	eventBus.AssertNumberOfCalls(t, "Publish", 1)
}

func testNodeReorgCheckUseCaseRunBitcoinFirstRunRefreshFailure(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_tip",
		ValidatedBlocks: new(big.Int).SetUint64(tipHeight),
	}
	btc.On("GetBlockchainInfo").Return(info, nil).Once()
	btc.On("GetBlockHashAtHeight", int64(tipHeight)).Return("", assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeBitcoin))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	err := uc.Run(context.Background(), entities.NodeTypeBitcoin)
	require.ErrorIs(t, err, assert.AnError)
	btc.AssertNotCalled(t, "GetBlockHeaderVerbose", mock.Anything)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinOneBlockAppendRefreshFailure(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_101",
		ValidatedBlocks: big.NewInt(101),
	}
	historyWindow := reorgHistoryWindow(2)

	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_101").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_100",
	}, nil).Once()
	btc.On("GetBlockHashAtHeight", int64(101)).Return("", assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeBitcoin))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	err := uc.Run(ctx, entities.NodeTypeBitcoin)
	require.ErrorIs(t, err, assert.AnError)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinDivergenceScanRPCFailure(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_101_fork",
		ValidatedBlocks: big.NewInt(101),
	}
	historyWindow := reorgHistoryWindow(2)

	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_101_fork").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_100_fork",
	}, nil).Once()
	btc.On("GetBlockHashAtHeight", int64(100)).Return("", assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeBitcoin))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	err := uc.Run(ctx, entities.NodeTypeBitcoin)
	require.ErrorIs(t, err, assert.AnError)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinAdvanceBeyondHistoryWindowRefreshFailure(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_113",
		ValidatedBlocks: big.NewInt(113),
	}
	historyWindow := reorgHistoryWindow(2)

	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_113").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_112",
	}, nil).Once()
	btc.On("GetBlockHashAtHeight", int64(113)).Return("", assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeBitcoin))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	err := uc.Run(ctx, entities.NodeTypeBitcoin)
	require.ErrorIs(t, err, assert.AnError)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

//nolint:funlen
func testNodeReorgCheckUseCaseRunBitcoinAboveThresholdReorgSendsAlert(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	alertSender := &mocks.AlertSenderMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_103_fork",
		ValidatedBlocks: big.NewInt(103),
	}
	historyWindow := reorgHistoryWindow(2)
	recipient := "alerts@example.com"
	alertBody := fmt.Sprintf(nodeReorgAlertBodyTemplate, entities.NodeTypeBitcoin, uint64(3), uint64(2))

	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_103_fork").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_102_fork",
	}, nil).Once()
	btc.On("GetBlockHashAtHeight", int64(100)).Return("hash_100_fork", nil).Once()
	btc.On("GetBlockHashAtHeight", int64(99)).Return("hash_99_fork", nil).Once()
	btc.On("GetBlockHashAtHeight", int64(98)).Return("hash_98_fork", nil).Once()
	btc.On("GetBlockHashAtHeight", int64(97)).Return("hash_97", nil).Once()
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(103) - i
		hash := fmt.Sprintf("hash_%d", blockHeight)
		switch blockHeight {
		case 103:
			hash = "hash_103_fork"
		case 102:
			hash = "hash_102_fork"
		case 101:
			hash = "hash_101_fork"
		case 100:
			hash = "hash_100_fork"
		case 99:
			hash = "hash_99_fork"
		case 98:
			hash = "hash_98_fork"
		default:
		}
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(hash, nil).Once()
	}
	alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectNodeReorg, alertBody, []string{recipient}).Return(nil).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    3,
		MaxAllowedDepth: 2,
		AboveThreshold:  true,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgAlertSentEvent(entities.NodeTypeBitcoin, 3))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		alertSender,
		recipient,
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	btc.AssertExpectations(t)
	alertSender.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunBitcoinGetBlockHeaderVerboseError(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	historyWindow := reorgHistoryWindow(2)
	for i := uint64(0); i < historyWindow; i++ {
		blockHeight := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(blockHeight)).Return(fmt.Sprintf("hash_%d", blockHeight), nil).Once()
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100_fork",
		ValidatedBlocks: big.NewInt(100),
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_100_fork").Return(blockchain.BitcoinBlockHeaderInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeBitcoin,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeBitcoin))).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)
	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	err := uc.Run(ctx, entities.NodeTypeBitcoin)
	require.ErrorIs(t, err, assert.AnError)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

//nolint:funlen
func TestNodeReorgCheckUseCase_HandleReorgAlert(t *testing.T) {
	t.Parallel()

	t.Run("alert cooldown suppresses repeated alerts", func(t *testing.T) {
		t.Parallel()
		alertSender := &mocks.AlertSenderMock{}
		eventBus := &mocks.EventBusMock{}
		recipient := "alerts@example.com"
		alertBody := fmt.Sprintf(nodeReorgAlertBodyTemplate, entities.NodeTypeBitcoin, uint64(3), uint64(2))
		alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectNodeReorg, alertBody, []string{recipient}).Return(nil).Once()
		eventBus.On("Publish", mock.MatchedBy(matchNodeReorgAlertSentEvent(entities.NodeTypeBitcoin, 3))).Return().Once()

		uc := NewNodeReorgCheckUseCase(
			blockchain.Rpc{},
			alertSender,
			recipient,
			eventBus,
			2,
			time.Hour,
		)

		err := uc.handleReorgAlert(context.Background(), entities.NodeTypeBitcoin, 3, 2)
		require.NoError(t, err)
		err = uc.handleReorgAlert(context.Background(), entities.NodeTypeBitcoin, 3, 2)
		require.NoError(t, err)
		alertSender.AssertExpectations(t)
		eventBus.AssertExpectations(t)
		alertSender.AssertNumberOfCalls(t, "SendAlert", 1)
		eventBus.AssertNumberOfCalls(t, "Publish", 1)
	})

	t.Run("alert send failure", func(t *testing.T) {
		t.Parallel()
		alertSender := &mocks.AlertSenderMock{}
		eventBus := &mocks.EventBusMock{}
		recipient := "alerts@example.com"
		alertBody := fmt.Sprintf(nodeReorgAlertBodyTemplate, entities.NodeTypeBitcoin, uint64(3), uint64(2))
		alertSender.On("SendAlert", mock.Anything, alerts.AlertSubjectNodeReorg, alertBody, []string{recipient}).Return(assert.AnError).Once()

		uc := NewNodeReorgCheckUseCase(
			blockchain.Rpc{},
			alertSender,
			recipient,
			eventBus,
			2,
			time.Hour,
		)

		err := uc.handleReorgAlert(context.Background(), entities.NodeTypeBitcoin, 3, 2)
		require.ErrorIs(t, err, assert.AnError)
		eventBus.AssertNotCalled(t, "Publish", mock.Anything)
		alertSender.AssertExpectations(t)
	})
}

//nolint:funlen
func TestNodeReorgCheckUseCase_Run_Rootstock(t *testing.T) {
	t.Parallel()

	t.Run("get height error", testNodeReorgCheckUseCaseRunRootstockGetHeightError)
	t.Run("first run", testNodeReorgCheckUseCaseRunRootstockFirstRun)

	t.Run("tip block fetch error", testNodeReorgCheckUseCaseRunRootstockTipBlockFetchError)

	t.Run("same tip no extra rpc", testNodeReorgCheckUseCaseRunRootstockSameTipNoExtraRPC)

	t.Run("one block append", testNodeReorgCheckUseCaseRunRootstockOneBlockAppend)

	t.Run("first run refresh failure", testNodeReorgCheckUseCaseRunRootstockFirstRunRefreshFailure)

	t.Run("reorg within history window", testNodeReorgCheckUseCaseRunRootstockReorgWithinHistoryWindow)

	t.Run("advance beyond history window", testNodeReorgCheckUseCaseRunRootstockAdvanceBeyondHistoryWindow)

	t.Run("one block append refresh failure", testNodeReorgCheckUseCaseRunRootstockOneBlockAppendRefreshFailure)

	t.Run("divergence scan rpc failure", testNodeReorgCheckUseCaseRunRootstockDivergenceScanRPCFailure)

	t.Run("advance beyond history window refresh failure", testNodeReorgCheckUseCaseRunRootstockAdvanceBeyondHistoryWindowRefreshFailure)
}

func testNodeReorgCheckUseCaseRunRootstockGetHeightError(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeRootstock))).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeRootstock)
	require.ErrorIs(t, err, assert.AnError)
	rsk.AssertNotCalled(t, "GetBlockByNumber", mock.Anything, mock.Anything)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunRootstockFirstRun(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	rsk.EXPECT().GetHeight(mock.Anything).Return(tipHeight, nil).Once()
	historyWindow := reorgHistoryWindow(2)
	tipBlock := blockchain.BlockInfo{
		Hash:       "0x_tip",
		ParentHash: "0x_parent",
	}
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == tipHeight
	})).Return(tipBlock, nil).Times(3)
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := tipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0x%064d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xparent%064d", currentBlockHeight),
		}, nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeRootstock)
	require.NoError(t, err)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunRootstockTipBlockFetchError(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == 100
	})).Return(blockchain.BlockInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeRootstock))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	err := uc.Run(context.Background(), entities.NodeTypeRootstock)
	require.ErrorIs(t, err, assert.AnError)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunRootstockSameTipNoExtraRPC(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	historyWindow := reorgHistoryWindow(2)
	tipBlock := blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}
	rsk.EXPECT().GetHeight(mock.Anything).Return(tipHeight, nil).Twice()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == tipHeight
	})).Return(tipBlock, nil).Times(3)
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := tipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Twice()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

//nolint:funlen
func testNodeReorgCheckUseCaseRunRootstockOneBlockAppend(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	initialTipHeight := uint64(100)
	nextTipHeight := uint64(101)
	historyWindow := reorgHistoryWindow(2)
	rsk.EXPECT().GetHeight(mock.Anything).Return(initialTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == initialTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := initialTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}
	rsk.EXPECT().GetHeight(mock.Anything).Return(nextTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == nextTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_101",
		ParentHash: "0xhash_100",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := nextTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Twice()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunRootstockFirstRunRefreshFailure(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	rsk.EXPECT().GetHeight(mock.Anything).Return(tipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == tipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == 99
	})).Return(blockchain.BlockInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeRootstock))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	err := uc.Run(context.Background(), entities.NodeTypeRootstock)
	require.ErrorIs(t, err, assert.AnError)
	eventBus.AssertExpectations(t)
}

//nolint:funlen
func testNodeReorgCheckUseCaseRunRootstockReorgWithinHistoryWindow(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	initialTipHeight := uint64(100)
	forkTipHeight := uint64(101)
	historyWindow := reorgHistoryWindow(2)

	rsk.EXPECT().GetHeight(mock.Anything).Return(initialTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == initialTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := initialTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}

	rsk.EXPECT().GetHeight(mock.Anything).Return(forkTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == forkTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_101_fork",
		ParentHash: "0xhash_100_fork",
	}, nil).Twice()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == 100
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100_fork",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == 99
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_99",
		ParentHash: "0xhash_98",
	}, nil).Twice()
	for i := uint64(2); i < historyWindow; i++ {
		blockHeight := forkTipHeight - i
		if blockHeight == 99 {
			continue
		}
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}

	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    1,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

//nolint:funlen
func testNodeReorgCheckUseCaseRunRootstockAdvanceBeyondHistoryWindow(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	initialTipHeight := uint64(100)
	advancedTipHeight := uint64(113)
	historyWindow := reorgHistoryWindow(2)

	rsk.EXPECT().GetHeight(mock.Anything).Return(initialTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == initialTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := initialTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}

	rsk.EXPECT().GetHeight(mock.Anything).Return(advancedTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == advancedTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_113",
		ParentHash: "0xhash_112",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := advancedTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}

	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	eventBus.AssertNumberOfCalls(t, "Publish", 1)
}

func testNodeReorgCheckUseCaseRunRootstockOneBlockAppendRefreshFailure(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	initialTipHeight := uint64(100)
	nextTipHeight := uint64(101)
	historyWindow := reorgHistoryWindow(2)
	rsk.EXPECT().GetHeight(mock.Anything).Return(initialTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == initialTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := initialTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}
	rsk.EXPECT().GetHeight(mock.Anything).Return(nextTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == nextTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_101",
		ParentHash: "0xhash_100",
	}, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == nextTipHeight
	})).Return(blockchain.BlockInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeRootstock))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	err := uc.Run(ctx, entities.NodeTypeRootstock)
	require.ErrorIs(t, err, assert.AnError)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunRootstockDivergenceScanRPCFailure(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	initialTipHeight := uint64(100)
	forkTipHeight := uint64(101)
	historyWindow := reorgHistoryWindow(2)
	rsk.EXPECT().GetHeight(mock.Anything).Return(initialTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == initialTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := initialTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}
	rsk.EXPECT().GetHeight(mock.Anything).Return(forkTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == forkTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_101_fork",
		ParentHash: "0xhash_100_fork",
	}, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == 100
	})).Return(blockchain.BlockInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeRootstock))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	err := uc.Run(ctx, entities.NodeTypeRootstock)
	require.ErrorIs(t, err, assert.AnError)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func testNodeReorgCheckUseCaseRunRootstockAdvanceBeyondHistoryWindowRefreshFailure(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	initialTipHeight := uint64(100)
	advancedTipHeight := uint64(113)
	historyWindow := reorgHistoryWindow(2)
	rsk.EXPECT().GetHeight(mock.Anything).Return(initialTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == initialTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_100",
		ParentHash: "0xhash_99",
	}, nil).Twice()
	for i := uint64(1); i < historyWindow; i++ {
		blockHeight := initialTipHeight - i
		currentBlockHeight := blockHeight
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
			return blockNumber != nil && blockNumber.Uint64() == currentBlockHeight
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0xhash_%d", currentBlockHeight),
			ParentHash: fmt.Sprintf("0xhash_%d", currentBlockHeight-1),
		}, nil).Once()
	}
	rsk.EXPECT().GetHeight(mock.Anything).Return(advancedTipHeight, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == advancedTipHeight
	})).Return(blockchain.BlockInfo{
		Hash:       "0xhash_113",
		ParentHash: "0xhash_112",
	}, nil).Once()
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(blockNumber *big.Int) bool {
		return blockNumber != nil && blockNumber.Uint64() == advancedTipHeight
	})).Return(blockchain.BlockInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckEventExpected(blockchain.NodeReorgCheckEvent{
		NodeType:        entities.NodeTypeRootstock,
		CurrentDepth:    0,
		MaxAllowedDepth: 2,
		AboveThreshold:  false,
	}))).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(matchNodeReorgCheckErrorEvent(entities.NodeTypeRootstock))).Return().Once()

	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		time.Hour,
	)

	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeRootstock))
	err := uc.Run(ctx, entities.NodeTypeRootstock)
	require.ErrorIs(t, err, assert.AnError)
	rsk.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}
