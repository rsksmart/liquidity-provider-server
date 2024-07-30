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
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	log "github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// nolint:funlen
func (s *IntegrationTestSuite) TestSuccessfulPegOutFlow() {

	var quote pkg.GetPegoutQuoteResponse
	var acceptedQuote pkg.AcceptPegoutResponse
	URL := s.config.Lps.Url

	s.Run("Should be able to get pegout quote", func() {
		body := pkg.PegoutQuoteRequest{
			To:               "n1zjV3WxJgA4dBfS5aMiEHtZsjTUvAL7p7",
			ValueToTransfer:  600000000000000000,
			RskRefundAddress: "0x79568c2989232dCa1840087D73d403602364c0D4",
		}

		result, err := execute[[]pkg.GetPegoutQuoteResponse](Execution{
			Method: http.MethodPost,
			URL:    URL + "/pegout/getQuotes",
			Body:   body,
		})
		s.NoError(err)

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
			"gasFee",
		}

		var rawResponse []map[string]any
		err = json.Unmarshal(result.RawResponse, &rawResponse)
		if err != nil {
			s.FailNow("Response does not have required format")
		}
		s.Equal(http.StatusOK, result.StatusCode)
		s.NotEmpty(rawResponse[0]["quoteHash"])
		s.NotEmpty(rawResponse[0]["quote"])
		quoteFields, ok := rawResponse[0]["quote"].(map[string]any)
		if !ok {
			s.FailNow("Quote is not an object")
		}
		s.AssertFields(expectedFields, quoteFields)
		quote = result.Response[0]
	})

	s.Run("Should be able to accept pegout quote", func() {
		body := pkg.AcceptQuoteRequest{QuoteHash: quote.QuoteHash}
		result, err := execute[pkg.AcceptPegoutResponse](Execution{
			Method: http.MethodPost,
			URL:    URL + "/pegout/acceptQuote",
			Body:   body,
		})
		s.Require().NoError(err)

		expectedFields := []string{
			"signature",
			"lbcAddress",
		}

		s.Equal(http.StatusOK, result.StatusCode)
		var rawResponse map[string]any
		err = json.Unmarshal(result.RawResponse, &rawResponse)
		if err != nil {
			s.FailNow("Response does not have required format")
		}
		s.AssertFields(expectedFields, rawResponse)
		acceptedQuote = result.Response
	})

	s.Run("Should process depositPegOut execution and transfer bitcoin to user", func() {
		var err error
		ctx := context.Background()
		privateKey, err := crypto.HexToECDSA(s.config.Rsk.UserPrivateKey)
		s.NoError(err)
		chainId, err := s.rsk.ChainID(ctx)
		if err != nil {
			s.FailNow(err.Error())
		}
		opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
		if err != nil {
			s.FailNow(err.Error())
		}

		value := entities.NewUWei(quote.Quote.Value)
		callFee := entities.NewUWei(quote.Quote.CallFee)
		gasFee := entities.NewUWei(quote.Quote.GasFee)
		productFee := entities.NewUWei(quote.Quote.ProductFeeAmount)
		totalFees := new(entities.Wei).Add(new(entities.Wei).Add(callFee, gasFee), productFee)
		totalAmount := new(entities.Wei).Add(totalFees, value)
		opts.Value = totalAmount.AsBigInt()

		gasPrice, err := s.rsk.SuggestGasPrice(ctx)
		s.NoError(err)
		opts.GasPrice = gasPrice
		pegoutQuote := parseLbcPegoutQuote(s, quote.Quote)

		signature, err := hex.DecodeString(acceptedQuote.Signature)
		s.NoError(err)

		depositTx, err := s.lbc.DepositPegout(opts, pegoutQuote, signature)
		s.NoError(err)
		log.Debug("[Integration test] Hash of deposit tx ", depositTx.Hash().String())

		address, err := btcutil.DecodeAddress(quote.Quote.DepositAddr, &s.btcParams)
		s.NoError(err)

		txHash := waitForBtcTransactionToAddress(s, address)

		txParsedHash, _ := chainhash.NewHashFromStr(txHash)
		tx, err := s.btc.GetTransaction(txParsedHash)
		s.NoError(err)
		s.NotNil(tx)
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
		s.NoError(err)

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case refund := <-eventChannel:
			subscription.Unsubscribe()
			s.NotNil(refund, "refundPegOut failed")
		case err = <-subscription.Err():
			if err != nil {
				s.FailNow("Error listening for refundPegOut", err)
			}
		case <-done:
			subscription.Unsubscribe()
			s.FailNow("Test cancelled")
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

func parseLbcPegoutQuote(s *IntegrationTestSuite, originalQuote pkg.PegoutQuoteDTO) bindings.QuotesPegOutQuote {
	lpBtcAddress, err := bitcoin.DecodeAddress(originalQuote.LpBTCAddr)
	s.NoError(err)
	btcRefundAddress, err := bitcoin.DecodeAddress(originalQuote.BtcRefundAddr)
	s.NoError(err)
	depositAddress, err := bitcoin.DecodeAddress(originalQuote.DepositAddr)
	s.NoError(err)
	return bindings.QuotesPegOutQuote{
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
		ProductFeeAmount:      big.NewInt(int64(originalQuote.ProductFeeAmount)),
		GasFee:                big.NewInt(int64(originalQuote.GasFee)),
	}
}

func waitForBtcTransactionToAddress(s *IntegrationTestSuite, address btcutil.Address) string {
	var latestBlockHash *chainhash.Hash
	var block *wire.MsgBlock
	info, err := s.btc.GetBlockChainInfo()
	s.NoError(err)
	latestBlockNumber := info.Blocks
	latestBlockHash, _ = chainhash.NewHashFromStr(info.BestBlockHash)
	block, err = s.btc.GetBlock(latestBlockHash)
	s.NoError(err)

	var txHash string
	for txHash == "" {
		txHash = lookForTxToAddress(block, address, &s.btcParams)
		if txHash != "" {
			return txHash
		}
		hash, getBlockError := s.btc.GetBlockHash(int64(latestBlockNumber + 1))
		if getBlockError != nil && !strings.Contains(getBlockError.Error(), "Block height out of range") {
			s.FailNow(getBlockError.Error())
		} else if getBlockError != nil && strings.Contains(getBlockError.Error(), "Block height out of range") {
			time.Sleep(10 * time.Second)
		} else if getBlockError == nil {
			latestBlockHash = hash
			latestBlockNumber++
			if block, err = s.btc.GetBlock(latestBlockHash); err != nil {
				s.FailNow(err.Error())
			}
		}
	}
	return txHash
}
