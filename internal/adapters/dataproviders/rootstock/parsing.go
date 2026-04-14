package rootstock

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

func ParseReceipt(tx *geth.Transaction, receipt *geth.Receipt) (blockchain.TransactionReceipt, error) {
	if tx == nil || receipt == nil {
		return blockchain.TransactionReceipt{}, errors.New("invalid parameters")
	}

	gasUsed := new(big.Int)
	gasUsed.SetUint64(receipt.GasUsed)
	cumulativeGasUsed := new(big.Int)
	cumulativeGasUsed.SetUint64(receipt.CumulativeGasUsed)
	from, err := geth.Sender(geth.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		if from, err = geth.Sender(geth.HomesteadSigner{}, tx); err != nil {
			return blockchain.TransactionReceipt{}, err
		}
	}

	result := blockchain.TransactionReceipt{
		TransactionHash:   receipt.TxHash.String(),
		BlockHash:         receipt.BlockHash.String(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		From:              from.String(),
		To:                tx.To().String(),
		CumulativeGasUsed: cumulativeGasUsed,
		GasUsed:           gasUsed,
		Value:             entities.NewBigWei(tx.Value()),
		Logs:              convertReceiptLogs(receipt),
		GasPrice:          entities.NewBigWei(tx.GasPrice()),
	}

	return result, nil
}

func convertReceiptLogs(receipt *geth.Receipt) []blockchain.TransactionLog {
	logs := make([]blockchain.TransactionLog, len(receipt.Logs))

	for i, eventLog := range receipt.Logs {
		topics := make([][32]byte, len(eventLog.Topics))
		for j, topic := range eventLog.Topics {
			topics[j] = topic
		}
		logs[i] = blockchain.TransactionLog{
			Address:     eventLog.Address.String(),
			Topics:      topics,
			Data:        eventLog.Data,
			BlockNumber: receipt.BlockNumber.Uint64(),
			TxHash:      receipt.TxHash.String(),
			TxIndex:     eventLog.TxIndex,
			BlockHash:   receipt.BlockHash.String(),
			Index:       eventLog.Index,
			Removed:     eventLog.Removed,
		}
	}

	return logs
}

// ParseDepositEventByQuoteHash iterates all PegOutDeposit events in the receipt and returns the
// one whose QuoteHash matches quoteHash AND whose emitting address matches lbcAddress
// (both comparisons are hex, case-insensitive, 0x-agnostic).
// It assumes the following event signature: event PegOutDeposit(bytes32 indexed quoteHash, address indexed sender, uint256 indexed timestamp, uint256 amount)
func ParseDepositEventByQuoteHash(
	receipt blockchain.TransactionReceipt,
	quoteHash string,
	lbcAddress string,
) (blockchain.ParsedLog[quote.PegoutDeposit], error) {
	const eventName = "PegOutDeposit"
	abi, err := bindings.PegoutContractMetaData.ParseABI()
	if err != nil {
		return blockchain.ParsedLog[quote.PegoutDeposit]{}, err
	}
	eventID := abi.Events[eventName].ID.Bytes()
	normalizedHash := strings.TrimPrefix(strings.ToLower(quoteHash), "0x")
	normalizedLBC := strings.TrimPrefix(strings.ToLower(lbcAddress), "0x")
	log, found := findDepositLog(receipt.Logs, eventID, normalizedHash, normalizedLBC)
	if !found {
		return blockchain.ParsedLog[quote.PegoutDeposit]{}, fmt.Errorf("deposit event not found for quote %s", quoteHash)
	}
	event := new(bindings.PegoutContractPegOutDeposit)
	if err = abi.UnpackIntoInterface(event, eventName, log.Data); err != nil {
		return blockchain.ParsedLog[quote.PegoutDeposit]{}, err
	}
	event.QuoteHash = common.BytesToHash(log.Topics[1][:])
	event.Sender = common.BytesToAddress(log.Topics[2][:])
	timestamp := new(big.Int)
	timestamp.SetBytes(log.Topics[3][:])
	event.Timestamp = timestamp
	return blockchain.ParsedLog[quote.PegoutDeposit]{
		Log: quote.PegoutDeposit{
			TxHash:      receipt.TransactionHash,
			QuoteHash:   hex.EncodeToString(event.QuoteHash[:]),
			Amount:      entities.NewBigWei(event.Amount),
			Timestamp:   time.Unix(event.Timestamp.Int64(), 0),
			BlockNumber: receipt.BlockNumber,
			From:        receipt.From,
		},
		RawLog: log,
	}, nil
}

func findDepositLog(
	logs []blockchain.TransactionLog,
	eventID []byte,
	normalizedHash string,
	normalizedLBC string,
) (blockchain.TransactionLog, bool) {
	const eventTopics = 4
	for _, log := range logs {
		if len(log.Topics) == eventTopics &&
			bytes.Equal(log.Topics[0][:], eventID) &&
			strings.TrimPrefix(strings.ToLower(log.Address), "0x") == normalizedLBC &&
			strings.EqualFold(hex.EncodeToString(log.Topics[1][:]), normalizedHash) {
			return log, true
		}
	}
	return blockchain.TransactionLog{}, false
}
