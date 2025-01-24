package main

import (
	"encoding/hex"
	"flag"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/term"
	"math/big"
	"os"
	"testing"
	"time"
)

// nolint:funlen
func TestExecuteRegisterPegIn(t *testing.T) {
	const resultHash = "0x111213"
	var (
		rawTx       = []byte{0x07, 0x08, 0x09}
		pmt         = []byte{0x04, 0x05, 0x06}
		blockHeight = big.NewInt(577)
	)
	parsedInput := ParsedRegisterPegInInput{
		Quote: quote.PeginQuote{
			FedBtcAddress:      test.AnyAddress,
			LbcAddress:         test.AnyAddress,
			LpRskAddress:       test.AnyAddress,
			BtcRefundAddress:   test.AnyAddress,
			RskRefundAddress:   test.AnyAddress,
			LpBtcAddress:       test.AnyAddress,
			CallFee:            entities.NewWei(1),
			PenaltyFee:         entities.NewWei(1),
			ContractAddress:    test.AnyAddress,
			Data:               "",
			GasLimit:           1,
			Nonce:              2,
			Value:              entities.NewWei(1),
			AgreementTimestamp: 3,
			TimeForDeposit:     4,
			LpCallTime:         5,
			Confirmations:      6,
			CallOnRegister:     false,
			GasFee:             entities.NewWei(1),
			ProductFeeAmount:   0,
		},
		Signature: []byte{0x01, 0x02, 0x03},
		BtcTxHash: "bitcoinTxHash",
	}
	blockInfo := blockchain.BitcoinBlockInformation{
		Hash:   [32]byte{0xa, 0xb, 0xc},
		Height: blockHeight,
		Time:   time.Now(),
	}
	t.Run("should execute RegisterPegIn successfully", func(t *testing.T) {
		rpc := new(mocks.BtcRpcMock)
		lbc := new(mocks.LbcMock)
		rpc.On("GetPartialMerkleTree", parsedInput.BtcTxHash).Return(pmt, nil).Once()
		rpc.On("GetRawTransaction", parsedInput.BtcTxHash).Return(rawTx, nil).Once()
		rpc.On("GetTransactionBlockInfo", parsedInput.BtcTxHash).Return(blockInfo, nil).Once()
		lbc.On("RegisterPegin", blockchain.RegisterPeginParams{
			QuoteSignature: parsedInput.Signature, BitcoinRawTransaction: rawTx,
			PartialMerkleTree: pmt, BlockHeight: blockHeight, Quote: parsedInput.Quote,
		}).Return(resultHash, nil).Once()

		result, err := ExecuteRegisterPegIn(rpc, lbc, parsedInput)
		require.NoError(t, err)
		assert.Equal(t, resultHash, result)
		rpc.AssertExpectations(t)
		lbc.AssertExpectations(t)
	})

	t.Run("should return error if GetPartialMerkleTree fails", func(t *testing.T) {
		rpc := new(mocks.BtcRpcMock)
		lbc := new(mocks.LbcMock)
		rpc.On("GetPartialMerkleTree", parsedInput.BtcTxHash).Return([]byte{}, assert.AnError).Once()

		result, err := ExecuteRegisterPegIn(rpc, lbc, parsedInput)
		require.Error(t, err)
		assert.Empty(t, result)
		rpc.AssertExpectations(t)
		lbc.AssertNotCalled(t, "RegisterPegin")
	})

	t.Run("should return error if GetRawTransaction fails", func(t *testing.T) {
		rpc := new(mocks.BtcRpcMock)
		lbc := new(mocks.LbcMock)
		rpc.On("GetPartialMerkleTree", parsedInput.BtcTxHash).Return(pmt, nil).Once()
		rpc.On("GetRawTransaction", parsedInput.BtcTxHash).Return([]byte{}, assert.AnError).Once()

		result, err := ExecuteRegisterPegIn(rpc, lbc, parsedInput)
		require.Error(t, err)
		assert.Empty(t, result)
		rpc.AssertExpectations(t)
		lbc.AssertNotCalled(t, "RegisterPegin")
	})

	t.Run("should return error if GetTransactionBlockInfo fails", func(t *testing.T) {
		rpc := new(mocks.BtcRpcMock)
		lbc := new(mocks.LbcMock)
		rpc.On("GetPartialMerkleTree", parsedInput.BtcTxHash).Return(pmt, nil).Once()
		rpc.On("GetRawTransaction", parsedInput.BtcTxHash).Return(rawTx, nil).Once()
		rpc.On("GetTransactionBlockInfo", parsedInput.BtcTxHash).Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()

		result, err := ExecuteRegisterPegIn(rpc, lbc, parsedInput)

		require.Error(t, err)
		assert.Empty(t, result)
		rpc.AssertExpectations(t)
		lbc.AssertNotCalled(t, "RegisterPegin")
	})

	t.Run("should return error if RegisterPegin fails", func(t *testing.T) {
		rpc := new(mocks.BtcRpcMock)
		lbc := new(mocks.LbcMock)
		rpc.On("GetPartialMerkleTree", parsedInput.BtcTxHash).Return(pmt, nil).Once()
		rpc.On("GetRawTransaction", parsedInput.BtcTxHash).Return(rawTx, nil).Once()
		rpc.On("GetTransactionBlockInfo", parsedInput.BtcTxHash).Return(blockInfo, nil).Once()
		lbc.On("RegisterPegin", mock.Anything).Return("", assert.AnError).Once()

		result, err := ExecuteRegisterPegIn(rpc, lbc, parsedInput)
		require.Error(t, err)
		assert.Empty(t, result)
		rpc.AssertExpectations(t)
		lbc.AssertExpectations(t)
	})
}

func TestRegisterPegInScriptInput_ToEnv(t *testing.T) {
	const btcRpcHost = "http://localhost:1111"
	t.Run("Should parse RegisterPegIn script input successfully", func(t *testing.T) {
		scriptInput := new(RegisterPegInScriptInput)
		ReadRegisterPegInScriptInput(scriptInput)
		require.NoError(t, flag.Set("network", "regtest"))
		require.NoError(t, flag.Set("input-file", "/file/path"))
		require.NoError(t, flag.Set("btc-rpc-host", btcRpcHost))
		require.NoError(t, flag.Set("btc-rpc-user", "btcUser"))
		require.NoError(t, flag.Set("btc-rpc-password", "btcPassword"))
		env, err := scriptInput.ToEnv(term.ReadPassword)
		require.NoError(t, err)
		assert.Equal(t, "regtest", env.LpsStage)
		assert.Equal(t, "regtest", env.Btc.Network)
		assert.Equal(t, btcRpcHost, env.Btc.Endpoint)
		assert.Equal(t, "btcUser", env.Btc.Username)
		assert.Equal(t, "btcPassword", env.Btc.Password)
		// we assert the rsk defaults to ensure the method called its parent method
		assert.Equal(t, uint64(33), env.Rsk.ChainId)
		assert.Equal(t, "0x8901a2bbf639bfd21a97004ba4d7ae2bd00b8da8", env.Rsk.LbcAddress)
		assert.Equal(t, "0x0000000000000000000000000000000001000006", env.Rsk.BridgeAddress)
		assert.Equal(t, 0, env.Rsk.AccountNumber)
	})

	test.ResetFlagSet()
	t.Run("Should use dummy credentials if credentials were not provided", func(t *testing.T) {
		scriptInput := new(RegisterPegInScriptInput)
		ReadRegisterPegInScriptInput(scriptInput)
		require.NoError(t, flag.Set("network", "regtest"))
		require.NoError(t, flag.Set("input-file", "/file/path"))
		require.NoError(t, flag.Set("btc-rpc-host", btcRpcHost))
		env, err := scriptInput.ToEnv(term.ReadPassword)
		require.NoError(t, err)
		assert.Equal(t, "regtest", env.LpsStage)
		assert.Equal(t, "regtest", env.Btc.Network)
		assert.Equal(t, btcRpcHost, env.Btc.Endpoint)
		assert.Equal(t, "none", env.Btc.Username)
		assert.Equal(t, "none", env.Btc.Password)
	})
}

func TestParseRegisterPegInScriptInput(t *testing.T) {
	parse := func() { /* mock function to prevent calling flag.Parse inside a test */ }
	input := RegisterPegInScriptInput{
		BaseInput: scripts.BaseInput{
			AwsLocalEndpoint: "http://localhost:4566", RskEndpoint: "http://localhost:4444", Network: "regtest",
			EncryptedJsonSecret: "secret1", EncryptedJsonPasswordSecret: "secret2", SecretSource: "aws",
		},
		InputFilePath: "input-example.json",
		BtcRcpHost:    "localhost:5555", BtcRpcUser: "user", BtcRpcPassword: "password",
	}
	t.Run("should fail on missing required fields", func(t *testing.T) {
		inputCopy := input
		inputCopy.InputFilePath = ""
		inputCopy.BtcRcpHost = ""
		result, err := ParseRegisterPegInScriptInput(parse, &inputCopy, os.ReadFile)
		require.ErrorContains(t, err, "Error:Field validation for 'InputFilePath'")
		require.ErrorContains(t, err, "Error:Field validation for 'BtcRcpHost'")
		assert.Empty(t, result)
	})
	t.Run("should fail on missing required fields of parent structure", func(t *testing.T) {
		inputCopy := input
		inputCopy.BaseInput.Network = ""
		inputCopy.BaseInput.SecretSource = ""
		inputCopy.BaseInput.RskEndpoint = ""
		result, err := ParseRegisterPegInScriptInput(parse, &inputCopy, os.ReadFile)
		require.ErrorContains(t, err, "Error:Field validation for 'Network'")
		require.ErrorContains(t, err, "Error:Field validation for 'SecretSource'")
		require.ErrorContains(t, err, "Error:Field validation for 'RskEndpoint'")
		assert.Empty(t, result)
	})
	t.Run("should fail if input file doesn't exists", func(t *testing.T) {
		inputCopy := input
		inputCopy.InputFilePath = "non-existing.json"
		result, err := ParseRegisterPegInScriptInput(parse, &inputCopy, os.ReadFile)
		require.ErrorContains(t, err, "no such file or directory")
		assert.Empty(t, result)
	})
	t.Run("should fail if input file is not a json", func(t *testing.T) {
		inputCopy := input
		inputCopy.InputFilePath = "register_pegin.go"
		result, err := ParseRegisterPegInScriptInput(parse, &inputCopy, os.ReadFile)
		require.Error(t, err)
		assert.Empty(t, result)
	})
	t.Run("should parse example successfully", func(t *testing.T) {
		result, err := ParseRegisterPegInScriptInput(parse, &input, os.ReadFile)
		require.NoError(t, err)
		assert.Equal(t, "e57767cefb13bb962e9729d99adbb7147f6054af6e8f4d7c4cd47e74cf9ccaa4", result.BtcTxHash)
		// nolint:errcheck
		signature, _ := hex.DecodeString("7290e2c28751d7e4ba2ea5fe5f8b1d3a0bfcd55089fddc0e74fe6809afb8195622801d2dd8267ea3cc4088f5e4b133e0e22dcc403ee0f838efbb277f493c8cde1b")
		assert.Equal(t, signature, result.Signature)
		assert.Equal(t, quote.PeginQuote{
			FedBtcAddress: "2N5muMepJizJE1gR7FbHJU6CD18V3BpNF9p", LbcAddress: "0x7557fcE0BbFAe81a9508FF469D481f2c72a8B5f3",
			LpRskAddress: "0x9d93929a9099be4355fc2389fbf253982f9df47c", BtcRefundAddress: "mfWxJ45yp2SFn7UciZyNpvDKrzbhyfKrY8",
			RskRefundAddress: "0x79568c2989232dCa1840087D73d403602364c0D4", LpBtcAddress: "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6",
			CallFee: entities.NewWei(10000000000000000), PenaltyFee: entities.NewWei(1000000000000000),
			ContractAddress: "0x79568c2989232dCa1840087D73d403602364c0D4",
			Data:            "", GasLimit: 46000, Value: entities.NewWei(600000000000000000),
			Nonce: 8941842587185974000, AgreementTimestamp: 1732101992, TimeForDeposit: 3600, LpCallTime: 7200,
			Confirmations: 10, CallOnRegister: false, GasFee: entities.NewWei(0), ProductFeeAmount: 0,
		}, result.Quote)
	})
}
