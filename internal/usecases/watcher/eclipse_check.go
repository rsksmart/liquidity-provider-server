package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	NodeEclipseDetectedError = errors.New("node eclipse detected")
	EclipseAlertSubject      = "Node Eclipse Detected"
	EclipseAlertBodyTemplate = "Your %s node is under eclipse attack. Please, check your node's connectivity and synchronization."
)

type EclipseCheckConfig struct {
	RskToleranceThreshold    uint8
	RskMaxMsWaitForBlock     uint64
	RskWaitPollingMsInterval uint64
	BtcToleranceThreshold    uint8
	BtcMaxMsWaitForBlock     uint64
	BtcWaitPollingMsInterval uint64
}

type blockIds struct {
	Hash   string
	Number uint64
}

type nodeEclipseCheckResult struct {
	OurBlock       blockIds
	ExternalBlocks []blockIds
}

type EclipseCheckUseCase struct {
	config             EclipseCheckConfig
	mainRpc            blockchain.Rpc
	externalBtcSources []blockchain.BitcoinNetwork
	externalRskSources []blockchain.RootstockRpcServer
	eventBus           entities.EventBus
	alertSender        entities.AlertSender
	alertRecipient     string
	eclipsedBlock      blockIds
	eclipsedBlockMutex sync.Locker
}

func NewEclipseCheckUseCase(
	config EclipseCheckConfig,
	mainRpc blockchain.Rpc,
	externalBtcSources []blockchain.BitcoinNetwork,
	externalRskSources []blockchain.RootstockRpcServer,
	eventBus entities.EventBus,
	alertSender entities.AlertSender,
	alertRecipient string,
	eclipsedBlockMutex sync.Locker,
) *EclipseCheckUseCase {
	return &EclipseCheckUseCase{
		config:             config,
		mainRpc:            mainRpc,
		externalBtcSources: externalBtcSources,
		externalRskSources: externalRskSources,
		eventBus:           eventBus,
		alertSender:        alertSender,
		alertRecipient:     alertRecipient,
		eclipsedBlockMutex: eclipsedBlockMutex,
	}
}

func (useCase *EclipseCheckUseCase) Run(ctx context.Context, nodeType entities.NodeType) error {
	useCase.eclipsedBlockMutex.Lock()
	defer useCase.eclipsedBlockMutex.Unlock()

	var err, alertErr error
	if nodeType == entities.NodeTypeBitcoin {
		err = useCase.checkBtcNode(ctx)
	} else if nodeType == entities.NodeTypeRootstock {
		err = useCase.checkRskNode(ctx)
	} else {
		return fmt.Errorf("unsupported node type: %s", nodeType)
	}

	if errors.Is(err, NodeEclipseDetectedError) {
		alertErr = useCase.triggerEclipseAlert(ctx, nodeType)
	}
	err = errors.Join(err, alertErr)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.EclipseCheckId, err)
	}
	return nil
}

func (useCase *EclipseCheckUseCase) triggerEclipseAlert(ctx context.Context, nodeType entities.NodeType) error {
	useCase.eventBus.Publish(blockchain.NodeEclipseEvent{
		BaseEvent:           entities.NewBaseEvent(blockchain.NodeEclipseEventId),
		NodeType:            nodeType,
		EclipsedBlockNumber: useCase.eclipsedBlock.Number,
		EclipsedBlockHash:   useCase.eclipsedBlock.Hash,
		DetectionTime:       time.Now(),
	})
	useCase.eclipsedBlock = blockIds{}
	return useCase.alertSender.SendAlert(
		ctx,
		EclipseAlertSubject,
		fmt.Sprintf(EclipseAlertBodyTemplate, nodeType),
		[]string{useCase.alertRecipient},
	)
}

func (useCase *EclipseCheckUseCase) checkBtcNode(ctx context.Context) error {
	log.Debugf("Eclipse check started for BTC node with %d external sources", len(useCase.externalBtcSources))
	checkResult, err := useCase.pullBtcBlocks()
	if err != nil {
		return err
	}

	successRate := getEclipseCheckSuccessRate(useCase.externalBtcSources, checkResult)
	if successRate >= useCase.config.BtcToleranceThreshold {
		log.Debugf("BTC node is above the tolerance threshold: %d%% (threshold: %d%%). No action needed.", successRate, useCase.config.BtcToleranceThreshold)
		return nil
	}

	log.Debugf("BTC node is under the tolerance threshold: %d%% (threshold: %d%%). Starting propagation deadline...", successRate, useCase.config.BtcToleranceThreshold)
	if err = useCase.waitForBtcSync(ctx); err != nil {
		return err
	}
	log.Debug("BTC node is synced with external sources again")
	return nil
}

func (useCase *EclipseCheckUseCase) checkRskNode(ctx context.Context) error {
	log.Debugf("Eclipse check started for RSK node with %d external sources", len(useCase.externalRskSources))
	checkResult, err := useCase.pullRskBlocks(ctx)
	if err != nil {
		return err
	}

	successRate := getEclipseCheckSuccessRate(useCase.externalRskSources, checkResult)
	if successRate >= useCase.config.RskToleranceThreshold {
		log.Debugf("RSK node is above the tolerance threshold: %d%% (threshold: %d%%). No action needed.", successRate, useCase.config.RskToleranceThreshold)
		return nil
	}

	log.Debugf("RSK node is under the tolerance threshold: %d%% (threshold: %d%%). Starting propagation deadline...", successRate, useCase.config.RskToleranceThreshold)
	if err = useCase.waitForRskSync(ctx); err != nil {
		return err
	}
	log.Debug("RSK node is synced with external sources again")
	return nil
}

func (useCase *EclipseCheckUseCase) pullBtcBlocks() (nodeEclipseCheckResult, error) {
	var err error
	var ourChain, externalChain blockchain.BitcoinBlockchainInfo
	result := nodeEclipseCheckResult{}

	if ourChain, err = useCase.mainRpc.Btc.GetBlockchainInfo(); err != nil {
		return nodeEclipseCheckResult{}, fmt.Errorf("error getting latest Bitcoin block from our node: %w", err)
	}
	result.OurBlock = useCase.getBtcBlockIds(ourChain)
	result.ExternalBlocks = make([]blockIds, 0)
	for _, btcSource := range useCase.externalBtcSources {
		if externalChain, err = btcSource.GetBlockchainInfo(); err != nil {
			log.Error("Error getting latest block from external Bitcoin source: ", err)
		} else {
			result.ExternalBlocks = append(result.ExternalBlocks, useCase.getBtcBlockIds(externalChain))
		}
	}
	return result, nil
}

func (useCase *EclipseCheckUseCase) pullRskBlocks(ctx context.Context) (nodeEclipseCheckResult, error) {
	var err error
	var ourBlock, externalBlock blockchain.BlockInfo
	result := nodeEclipseCheckResult{}

	if ourBlock, err = useCase.mainRpc.Rsk.GetBlockByNumber(ctx, nil); err != nil {
		return nodeEclipseCheckResult{}, fmt.Errorf("error getting RSK block from our node: %w", err)
	}
	result.OurBlock = useCase.getRskBlockIds(ourBlock)
	result.ExternalBlocks = make([]blockIds, 0)
	for _, rskSource := range useCase.externalRskSources {
		if externalBlock, err = rskSource.GetBlockByNumber(ctx, nil); err != nil {
			log.Error("Error getting block from external RSK source: ", err)
		} else {
			result.ExternalBlocks = append(result.ExternalBlocks, useCase.getRskBlockIds(externalBlock))
		}
	}
	return result, nil
}

func (useCase *EclipseCheckUseCase) waitForBtcSync(ctx context.Context) error {
	var ourLatestBlock blockIds
	ticker := utils.NewTickerWrapper(time.Duration(useCase.config.BtcWaitPollingMsInterval) * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(time.Duration(useCase.config.BtcMaxMsWaitForBlock) * time.Millisecond)
	defer timer.Stop()

btcBlockWaitLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			useCase.eclipsedBlock = ourLatestBlock
			return NodeEclipseDetectedError
		case <-ticker.C():
			checkResult, err := useCase.pullBtcBlocks()
			if err != nil {
				return err
			}
			ourLatestBlock = checkResult.OurBlock
			successRate := getEclipseCheckSuccessRate(useCase.externalBtcSources, checkResult)
			if successRate >= useCase.config.BtcToleranceThreshold {
				break btcBlockWaitLoop
			}
		}
	}
	return nil
}

func (useCase *EclipseCheckUseCase) waitForRskSync(ctx context.Context) error {
	var ourLatestBlock blockIds
	ticker := utils.NewTickerWrapper(time.Duration(useCase.config.RskWaitPollingMsInterval) * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(time.Duration(useCase.config.RskMaxMsWaitForBlock) * time.Millisecond)
	defer timer.Stop()

rskBlockWaitLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			useCase.eclipsedBlock = ourLatestBlock
			return NodeEclipseDetectedError
		case <-ticker.C():
			checkResult, err := useCase.pullRskBlocks(ctx)
			if err != nil {
				return err
			}
			ourLatestBlock = checkResult.OurBlock
			successRate := getEclipseCheckSuccessRate(useCase.externalRskSources, checkResult)
			if successRate >= useCase.config.RskToleranceThreshold {
				break rskBlockWaitLoop
			}
		}
	}
	return nil
}

func (useCase *EclipseCheckUseCase) getRskBlockIds(block blockchain.BlockInfo) blockIds {
	return blockIds{
		Hash:   block.Hash,
		Number: block.Number,
	}
}

func (useCase *EclipseCheckUseCase) getBtcBlockIds(btcChain blockchain.BitcoinBlockchainInfo) blockIds {
	return blockIds{
		Hash:   btcChain.BestBlockHash,
		Number: btcChain.ValidatedBlocks.Uint64(),
	}
}

func getEclipseCheckSuccessRate[T []blockchain.BitcoinNetwork | []blockchain.RootstockRpcServer](source T, checkResult nodeEclipseCheckResult) uint8 {
	total := uint8(len(source))
	if total == 0 {
		return 0
	}
	var successCount uint8 = 0
	for _, externalBlock := range checkResult.ExternalBlocks {
		if externalBlock.Number == checkResult.OurBlock.Number && externalBlock.Hash == checkResult.OurBlock.Hash {
			successCount++
		} else {
			log.Debugf("Block mismatch: our block %d (%s) vs external block %d (%s)",
				checkResult.OurBlock.Number, checkResult.OurBlock.Hash, externalBlock.Number, externalBlock.Hash)
		}
	}
	return (successCount * 100) / total
}
