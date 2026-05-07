package rootstock

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

var releaseRequestRejectedTopic = crypto.Keccak256Hash([]byte("release_request_rejected(address,uint256,int256)"))

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

// ParseDepositEvent parses a PegOutDeposit event from a transaction receipt.
// It assumes the following event signature:event PegOutDeposit(bytes32 indexed quoteHash, address indexed sender, uint256 amount, uint256 timestamp);
func ParseDepositEvent(receipt blockchain.TransactionReceipt) (blockchain.ParsedLog[quote.PegoutDeposit], error) {
	const (
		eventName   = "PegOutDeposit"
		eventTopics = 4
	)
	abi, err := bindings.IPegOutMetaData.GetAbi()
	if err != nil {
		return blockchain.ParsedLog[quote.PegoutDeposit]{}, err
	}
	index := slices.IndexFunc(receipt.Logs, func(log blockchain.TransactionLog) bool {
		return bytes.Equal(log.Topics[0][:], abi.Events[eventName].ID.Bytes())
	})
	if index < 0 {
		return blockchain.ParsedLog[quote.PegoutDeposit]{}, errors.New("deposit event not found in receipt logs")
	}

	log := receipt.Logs[index]
	if len(log.Topics) != eventTopics {
		return blockchain.ParsedLog[quote.PegoutDeposit]{}, errors.New("invalid number of topics for PegOutDeposit event")
	}

	event := new(bindings.IPegOutPegOutDeposit)
	if err = abi.UnpackIntoInterface(event, eventName, log.Data); err != nil {
		return blockchain.ParsedLog[quote.PegoutDeposit]{}, err
	}

	// indexed args
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

// ParseReleaseRejection scans receipt.Logs for a Bridge release_request_rejected event emitted by bridgeAddress.
//
// Event ABI (from rskj BridgeEvents.RELEASE_REQUEST_REJECTED):
//
//	release_request_rejected(address indexed sender, uint256 amount, int256 reason)
//
// Non-indexed params (amount, reason) are packed into log.Data as two 32-byte words; reason is the last word.
// Reason codes are declared by rskj RejectedPegoutReason.
func ParseReleaseRejection(receipt blockchain.TransactionReceipt, bridgeAddress string) (rejected bool, reason blockchain.RejectedPegoutReason) {
	normalizedBridge := strings.TrimPrefix(strings.ToLower(bridgeAddress), "0x")
	log, found := findReleaseRejectedLog(receipt.Logs, normalizedBridge)
	if !found {
		return false, ""
	}
	if len(log.Data) < 64 {
		return true, blockchain.RejectedPegoutReasonUnknown
	}
	// Decode the trailing int256 reason word using two's complement. big.Int.SetBytes treats
	// its input as unsigned, so a negative int256 (high bit of the leading byte set) would
	// otherwise decode as a near-2^256 positive number. If the sign bit is set, subtract 2^256
	// to recover the signed value. Current Bridge code only emits positive reason codes, but
	// the param is declared int256 so we decode it as signed for forward compatibility.
	reasonBytes := log.Data[len(log.Data)-32:]
	r := new(big.Int).SetBytes(reasonBytes)
	if reasonBytes[0]&0x80 != 0 {
		r.Sub(r, new(big.Int).Lsh(big.NewInt(1), 256))
	}
	return true, parseRejectedPegoutReason(r.Int64())
}

func parseRejectedPegoutReason(reason int64) blockchain.RejectedPegoutReason {
	switch reason {
	case 1:
		return blockchain.RejectedPegoutReasonLowAmount
	case 2:
		return blockchain.RejectedPegoutReasonCallerContract
	case 3:
		return blockchain.RejectedPegoutReasonFeeAboveValue
	default:
		return blockchain.RejectedPegoutReasonUnknown
	}
}

func findReleaseRejectedLog(logs []blockchain.TransactionLog, normalizedBridge string) (blockchain.TransactionLog, bool) {
	for _, log := range logs {
		if len(log.Topics) > 0 &&
			bytes.Equal(log.Topics[0][:], releaseRequestRejectedTopic.Bytes()) &&
			strings.TrimPrefix(strings.ToLower(log.Address), "0x") == normalizedBridge {
			return log, true
		}
	}
	return blockchain.TransactionLog{}, false
}
