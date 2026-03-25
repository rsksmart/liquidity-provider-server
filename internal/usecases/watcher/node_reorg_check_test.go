package watcher

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeType("unknown"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

func TestNodeReorgCheckUseCase_Run_Bitcoin_GetBlockchainInfoError(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	btc.On("GetBlockchainInfo").Return(blockchain.BitcoinBlockchainInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		_, ok := e.(blockchain.NodeReorgCheckErrorEvent)
		return ok
	})).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeBitcoin)
	require.Error(t, err)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestNodeReorgCheckUseCase_Run_Bitcoin_FirstRun(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_tip",
		ValidatedBlocks: new(big.Int).SetUint64(tipHeight),
	}
	btc.On("GetBlockchainInfo").Return(info, nil).Once()
	hw := reorgHistoryWindow(2)
	for i := uint64(0); i < hw; i++ {
		h := tipHeight - i
		btc.On("GetBlockHashAtHeight", int64(h)).Return(fmt.Sprintf("hash_%d", h), nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodeReorgCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin && ev.CurrentDepth == 0 &&
			ev.MaxAllowedDepth == 2 && !ev.AboveThreshold
	})).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeBitcoin)
	require.NoError(t, err)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestNodeReorgCheckUseCase_Run_Bitcoin_SameTipNoExtraRpc(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_tip",
		ValidatedBlocks: new(big.Int).SetUint64(tipHeight),
	}
	btc.On("GetBlockchainInfo").Return(info, nil).Twice()
	hw := reorgHistoryWindow(2)
	for i := uint64(0); i < hw; i++ {
		h := tipHeight - i
		btc.On("GetBlockHashAtHeight", int64(h)).Return(fmt.Sprintf("hash_%d", h), nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodeReorgCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeBitcoin && ev.CurrentDepth == 0 && !ev.AboveThreshold
	})).Return().Twice()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		2,
		time.Hour,
	)
	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestNodeReorgCheckUseCase_Run_Bitcoin_OneBlockAppend(t *testing.T) {
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
	hw := reorgHistoryWindow(2)
	for i := uint64(0); i < hw; i++ {
		h := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(h)).Return(fmt.Sprintf("hash_%d", h), nil).Once()
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_101").Return(blockchain.BitcoinBlockHeaderInfo{
		PreviousHash: "hash_100",
	}, nil).Once()
	for i := uint64(0); i < hw; i++ {
		h := uint64(101) - i
		btc.On("GetBlockHashAtHeight", int64(h)).Return(fmt.Sprintf("hash_%d", h), nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		_, ok := e.(blockchain.NodeReorgCheckEvent)
		return ok
	})).Return().Times(2)
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		2,
		time.Hour,
	)
	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}

func TestNodeReorgCheckUseCase_Run_Rootstock_GetHeightError(t *testing.T) {
	t.Parallel()
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(0), assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		_, ok := e.(blockchain.NodeReorgCheckErrorEvent)
		return ok
	})).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeRootstock)
	require.Error(t, err)
	eventBus.AssertExpectations(t)
}

func TestNodeReorgCheckUseCase_Run_Rootstock_FirstRun(t *testing.T) {
	t.Parallel()
	const tipHeight = uint64(100)
	rsk := &mocks.RootstockRpcServerMock{}
	eventBus := &mocks.EventBusMock{}
	rsk.EXPECT().GetHeight(mock.Anything).Return(tipHeight, nil).Once()
	hw := reorgHistoryWindow(2)
	tipBlock := blockchain.BlockInfo{
		Hash:       "0x_tip",
		ParentHash: "0x_parent",
	}
	rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
		return b != nil && b.Uint64() == tipHeight
	})).Return(tipBlock, nil).Twice()
	for i := uint64(1); i < hw; i++ {
		h := tipHeight - i
		hh := h
		rsk.EXPECT().GetBlockByNumber(mock.Anything, mock.MatchedBy(func(b *big.Int) bool {
			return b != nil && b.Uint64() == hh
		})).Return(blockchain.BlockInfo{
			Hash:       fmt.Sprintf("0x%064d", hh),
			ParentHash: fmt.Sprintf("0xparent%064d", hh),
		}, nil).Once()
	}
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		ev, ok := e.(blockchain.NodeReorgCheckEvent)
		return ok && ev.NodeType == entities.NodeTypeRootstock && ev.CurrentDepth == 0 && !ev.AboveThreshold
	})).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Rsk: rsk},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		2,
		time.Hour,
	)
	err := uc.Run(context.Background(), entities.NodeTypeRootstock)
	require.NoError(t, err)
	eventBus.AssertExpectations(t)
}

func TestNodeReorgCheckUseCase_Run_Bitcoin_GetBlockHeaderVerboseError(t *testing.T) {
	t.Parallel()
	btc := &mocks.BtcRpcMock{}
	eventBus := &mocks.EventBusMock{}
	info1 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100",
		ValidatedBlocks: big.NewInt(100),
	}
	btc.On("GetBlockchainInfo").Return(info1, nil).Once()
	hw := reorgHistoryWindow(2)
	for i := uint64(0); i < hw; i++ {
		h := uint64(100) - i
		btc.On("GetBlockHashAtHeight", int64(h)).Return(fmt.Sprintf("hash_%d", h), nil).Once()
	}
	info2 := blockchain.BitcoinBlockchainInfo{
		BestBlockHash:   "hash_100_fork",
		ValidatedBlocks: big.NewInt(100),
	}
	btc.On("GetBlockchainInfo").Return(info2, nil).Once()
	btc.On("GetBlockHeaderVerbose", "hash_100_fork").Return(blockchain.BitcoinBlockHeaderInfo{}, assert.AnError).Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		_, ok := e.(blockchain.NodeReorgCheckEvent)
		return ok
	})).Return().Once()
	eventBus.On("Publish", mock.MatchedBy(func(e entities.Event) bool {
		_, ok := e.(blockchain.NodeReorgCheckErrorEvent)
		return ok
	})).Return().Once()
	uc := NewNodeReorgCheckUseCase(
		blockchain.Rpc{Btc: btc},
		&mocks.AlertSenderMock{},
		"",
		eventBus,
		2,
		2,
		time.Hour,
	)
	ctx := context.Background()
	require.NoError(t, uc.Run(ctx, entities.NodeTypeBitcoin))
	err := uc.Run(ctx, entities.NodeTypeBitcoin)
	require.Error(t, err)
	btc.AssertExpectations(t)
	eventBus.AssertExpectations(t)
}
