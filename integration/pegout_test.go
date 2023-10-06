package integration_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/connectors/bindings"
	lps "github.com/rsksmart/liquidity-provider-server/http"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net/http"
	"strings"
	"time"
)

func (s *IntegrationTestSuite) TestSuccessfulPegOutFlow() {

	var quote lps.QuotePegOutResponse
	var acceptedQuote lps.AcceptResPegOut
	URL := s.config.Lps.Url

	s.Run("Should be able to get pegout quote", func() {
		body := lps.QuotePegOutRequest{
			To:                   "mz5RDWsNN38ehxNfKozmt4n1dFnV9BjJ5e",
			ValueToTransfer:      600000000000000000,
			RskRefundAddress:     "0x79568c2989232dCa1840087D73d403602364c0D4",
			BitcoinRefundAddress: "mz5RDWsNN38ehxNfKozmt4n1dFnV9BjJ5e",
		}

		result, err := execute[[]lps.QuotePegOutResponse](Execution{
			Method: http.MethodPost,
			URL:    URL + "/pegout/getQuotes",
			Body:   body,
		})

		if err != nil {
			assert.Fail(s.T(), "Unexpected error: ", err)
		}

		expectedFields := []string{
			"lbcAddress",
			"liquidityProviderRskAddress",
			"btcRefundAddress",
			"rskRefundAddress",
			"lpBtcAddr",
			"callFee",
			"penaltyFee",
			"nonce",
			"depositAddr",
			"value",
			"agreementTimestamp",
			"depositDateLimit",
			"depositConfirmations",
			"transferConfirmations",
			"transferTime",
			"expireDate",
			"expireBlocks",
			"callCost",
		}

		var rawResponse []map[string]any
		err = json.Unmarshal(result.RawResponse, &rawResponse)
		if err != nil {
			assert.Fail(s.T(), "Response does not have required format")
		}
		assert.Equal(s.T(), http.StatusOK, result.StatusCode)
		assert.NotEmpty(s.T(), rawResponse[0]["quoteHash"])
		assert.NotEmpty(s.T(), rawResponse[0]["quote"])
		quoteFields, ok := rawResponse[0]["quote"].(map[string]any)
		if !ok {
			assert.Fail(s.T(), "Quote is not an object")
		}
		s.AssertFields(expectedFields, quoteFields)
		quote = result.Response[0]
	})

	s.Run("Should be able to accept pegout quote", func() {
		body := lps.AcceptReq{QuoteHash: quote.QuoteHash}
		result, err := execute[lps.AcceptResPegOut](Execution{
			Method: http.MethodPost,
			URL:    URL + "/pegout/acceptQuote",
			Body:   body,
		})
		if err != nil {
			assert.Fail(s.T(), "Unexpected error: ", err)
		}

		expectedFields := []string{
			"signature",
			"lbcAddress",
		}

		assert.Equal(s.T(), http.StatusOK, result.StatusCode)
		var rawResponse map[string]any
		err = json.Unmarshal(result.RawResponse, &rawResponse)
		if err != nil {
			assert.Fail(s.T(), "Response does not have required format")
		}
		s.AssertFields(expectedFields, rawResponse)
		acceptedQuote = result.Response
	})

	s.Run("Should process depositPegOut execution and transfer bitcoin to user", func() {
		var err error
		ctx := context.Background()
		privateKey, err := crypto.HexToECDSA(s.config.Rsk.UserPrivateKey)
		if err != nil {
			assert.FailNow(s.T(), "Invalid private key")
		}
		chainId, err := s.rsk.ChainID(ctx)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
		opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}

		value := types.NewWei(int64(quote.Quote.Value))
		callFee := types.NewWei(int64(quote.Quote.CallFee))
		callCost := types.NewWei(int64(quote.Quote.CallCost))
		totalFees := new(types.Wei).Add(callFee, callCost)
		totalAmount := new(types.Wei).Add(totalFees, value)
		opts.Value = totalAmount.AsBigInt()

		gasPrice, err := s.rsk.SuggestGasPrice(ctx)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
		opts.GasPrice = gasPrice

		originalQuote := quote.Quote
		lpBtcAddress, err := connectors.DecodeBTCAddress(originalQuote.LpBTCAddr)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
		btcRefundAddress, err := connectors.DecodeBTCAddress(originalQuote.BtcRefundAddr)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
		depositAddress, err := connectors.DecodeBTCAddress(originalQuote.DepositAddr)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
		pegoutQuote := bindings.QuotesPegOutQuote{
			LbcAddress:            common.HexToAddress(originalQuote.LBCAddr),
			LpRskAddress:          common.HexToAddress(originalQuote.LPRSKAddr),
			BtcRefundAddress:      btcRefundAddress,
			RskRefundAddress:      common.HexToAddress(originalQuote.RSKRefundAddr),
			LpBtcAddress:          lpBtcAddress,
			CallFee:               big.NewInt(int64(originalQuote.CallFee)),
			PenaltyFee:            big.NewInt(int64(originalQuote.PenaltyFee)),
			Nonce:                 originalQuote.Nonce,
			DeposityAddress:       depositAddress,
			Value:                 big.NewInt(int64(originalQuote.Value)),
			AgreementTimestamp:    originalQuote.AgreementTimestamp,
			DepositDateLimit:      originalQuote.DepositDateLimit,
			DepositConfirmations:  originalQuote.DepositConfirmations,
			TransferConfirmations: originalQuote.TransferConfirmations,
			TransferTime:          originalQuote.TransferTime,
			ExpireDate:            originalQuote.ExpireDate,
			ExpireBlock:           originalQuote.ExpireBlock,
		}

		signature, err := hex.DecodeString(acceptedQuote.Signature)
		if err != nil {
			assert.FailNow(s.T(), "invalid signature")
		}

		depositTx, err := s.lbc.DepositPegout(opts, pegoutQuote, signature)
		if err != nil {
			assert.FailNow(s.T(), "error depositing pegout")
		}
		log.Debug("[Integration test] Hash of deposit tx ", depositTx.Hash().String())

		address, err := btcutil.DecodeAddress(quote.Quote.DepositAddr, &s.btcParams)
		if err != nil {
			assert.FailNow(s.T(), "invalid btc address")
		}

		var latestBlockHash *chainhash.Hash
		var block *wire.MsgBlock
		info, err := s.btc.GetBlockChainInfo()
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
		latestBlockNumber := info.Blocks
		latestBlockHash, _ = chainhash.NewHashFromStr(info.BestBlockHash)
		block, err = s.btc.GetBlock(latestBlockHash)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}

		var txHash string
		for txHash == "" {
			txHash = lookForTxToAddress(block, address, &s.btcParams)
			if txHash == "" {
				hash, getBlockError := s.btc.GetBlockHash(int64(latestBlockNumber + 1))
				if getBlockError != nil && !strings.Contains(getBlockError.Error(), "Block height out of range") {
					assert.FailNow(s.T(), getBlockError.Error())
				} else if getBlockError == nil {
					latestBlockHash = hash
					latestBlockNumber++
					block, err = s.btc.GetBlock(latestBlockHash)
					if err != nil {
						assert.FailNow(s.T(), err.Error())
					}
				}
			}
			time.Sleep(10 * time.Second)
		}

		txParsedHash, _ := chainhash.NewHashFromStr(txHash)
		tx, err := s.btc.GetTransaction(txParsedHash)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
		assert.NotNil(s.T(), tx)
	})

	s.Run("Should refund pegout to liquidity provider", func() {
		eventChannel := make(chan *bindings.LiquidityBridgeContractPegOutRefunded)
		var quoteHash [32]byte
		hashBytes, _ := hex.DecodeString(quote.QuoteHash)
		copy(quoteHash[:], hashBytes)
		subscription, err := s.lbc.WatchPegOutRefunded(
			nil,
			eventChannel,
			[][32]byte{quoteHash},
		)
		if err != nil {
			assert.FailNow(s.T(), "Error listening for refundPegOut")
		}

		select {
		case refund := <-eventChannel:
			subscription.Unsubscribe()
			assert.NotNil(s.T(), refund, "refundPegOut failed")
		case err = <-subscription.Err():
			assert.FailNow(s.T(), "Error listening for refundPegOut")
		}
	})
}

func lookForTxToAddress(block *wire.MsgBlock, target btcutil.Address, params *chaincfg.Params) string {
	for _, tx := range block.Transactions {
		for _, output := range tx.TxOut {
			_, addresses, _, _ := txscript.ExtractPkScriptAddrs(output.PkScript, params)
			if len(addresses) != 0 && addresses[0].EncodeAddress() == target.EncodeAddress() {
				return tx.TxHash().String()
			}
		}
	}
	return ""
}
