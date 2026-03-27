//go:build integration

package mongodb_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
)

var (
	defaultNow              = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	defaultLbcAddress       = "0xaabbccdd11223344556677889900aabbccddeeff"
	defaultLpRskAddress     = "0x1234567890abcdef1234567890abcdef12345678"
	defaultBtcRefundAddress = "mzBc4XEFSdzCDcTxAgf6EZXgsZWpztRhef"
	defaultRskRefundAddress = "0xabcdef1234567890abcdef1234567890abcdef12"
	defaultLpBtcAddress     = "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"
	defaultOwnerAddress     = "0x1234567890abcdef1234567890abcdef12345678"
)

func randomHash() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to read random bytes: %v", err))
	}
	return hex.EncodeToString(b)
}

func newTestPeginQuote(hash string) quote.CreatedPeginQuote {
	return quote.CreatedPeginQuote{
		Hash: hash,
		Quote: quote.PeginQuote{
			FedBtcAddress:      "2N1234567890abcdef",
			LbcAddress:         defaultLbcAddress,
			LpRskAddress:       defaultLpRskAddress,
			BtcRefundAddress:   defaultBtcRefundAddress,
			RskRefundAddress:   defaultRskRefundAddress,
			LpBtcAddress:       defaultLpBtcAddress,
			CallFee:            entities.NewWei(1000000),
			PenaltyFee:         entities.NewWei(5000000),
			ContractAddress:    defaultLpRskAddress,
			Data:               "0x",
			GasLimit:           100000,
			Nonce:              42,
			Value:              entities.NewWei(500000000000000000),
			AgreementTimestamp: uint32(defaultNow.Unix()),
			TimeForDeposit:     3600,
			LpCallTime:         7200,
			Confirmations:      10,
			CallOnRegister:     false,
			GasFee:             entities.NewWei(21000),
			ChainId:            31,
		},
		CreationData: quote.PeginCreationData{
			GasPrice:      entities.NewWei(60000000),
			FeePercentage: utils.NewBigFloat64(0.5),
			FixedFee:      entities.NewWei(1000000),
		},
	}
}

func newTestRetainedPeginQuote(hash string, state quote.PeginState) quote.RetainedPeginQuote {
	return quote.RetainedPeginQuote{
		QuoteHash:             hash,
		DepositAddress:        "2N1234567890abcdef",
		Signature:             "0xsignature123",
		RequiredLiquidity:     entities.NewWei(600000000000000000),
		State:                 state,
		UserBtcTxHash:         "btctx" + hash[:8],
		CallForUserTxHash:     "0xcalltx" + hash[:8],
		RegisterPeginTxHash:   "0xregtx" + hash[:8],
		CallForUserGasUsed:    50000,
		CallForUserGasPrice:   entities.NewWei(60000000),
		RegisterPeginGasUsed:  80000,
		RegisterPeginGasPrice: entities.NewWei(60000000),
		OwnerAccountAddress:   defaultOwnerAddress,
	}
}

func newTestPegoutQuote(hash string) quote.CreatedPegoutQuote {
	return quote.CreatedPegoutQuote{
		Hash: hash,
		Quote: quote.PegoutQuote{
			LbcAddress:            defaultLbcAddress,
			LpRskAddress:          defaultLpRskAddress,
			BtcRefundAddress:      defaultBtcRefundAddress,
			RskRefundAddress:      defaultRskRefundAddress,
			LpBtcAddress:          defaultLpBtcAddress,
			CallFee:               entities.NewWei(2000000),
			PenaltyFee:            entities.NewWei(10000000),
			Nonce:                 99,
			DepositAddress:        "0xdeposit1234567890abcdef1234567890abcdef",
			Value:                 entities.NewWei(1000000000000000000),
			AgreementTimestamp:    uint32(defaultNow.Unix()),
			DepositDateLimit:      uint32(defaultNow.Add(2 * time.Hour).Unix()),
			DepositConfirmations:  10,
			TransferConfirmations: 10,
			TransferTime:          7200,
			ExpireDate:            uint32(defaultNow.Add(24 * time.Hour).Unix()),
			ExpireBlock:           100000,
			GasFee:                entities.NewWei(42000),
			ChainId:               31,
		},
		CreationData: quote.PegoutCreationData{
			FeeRate:       utils.NewBigFloat64(0.001),
			FeePercentage: utils.NewBigFloat64(0.5),
			GasPrice:      entities.NewWei(60000000),
			FixedFee:      entities.NewWei(2000000),
		},
	}
}

func newTestRetainedPegoutQuote(hash string, state quote.PegoutState) quote.RetainedPegoutQuote {
	return quote.RetainedPegoutQuote{
		QuoteHash:            hash,
		DepositAddress:       "0xdeposit1234567890abcdef1234567890abcdef",
		Signature:            "0xsignature456",
		RequiredLiquidity:    entities.NewWei(1200000000000000000),
		State:                state,
		UserRskTxHash:        "0xusertx" + hash[:8],
		LpBtcTxHash:          "btctx" + hash[:8],
		RefundPegoutTxHash:   "0xrefundtx" + hash[:8],
		BridgeRefundTxHash:   "0xbridgetx" + hash[:8],
		BridgeRefundGasUsed:  60000,
		BridgeRefundGasPrice: entities.NewWei(60000000),
		RefundPegoutGasUsed:  70000,
		RefundPegoutGasPrice: entities.NewWei(60000000),
		SendPegoutBtcFee:     entities.NewWei(15000),
		BtcReleaseTxHash:     "btcrelease" + hash[:8],
		OwnerAccountAddress:  defaultOwnerAddress,
	}
}

func newTestPegoutDeposit(txHash, quoteHash string) quote.PegoutDeposit {
	return quote.PegoutDeposit{
		TxHash:      txHash,
		QuoteHash:   quoteHash,
		Amount:      entities.NewWei(1000000000000000000),
		Timestamp:   defaultNow,
		BlockNumber: 50000,
		From:        defaultOwnerAddress,
	}
}

func newTestPeginConfig() entities.Signed[liquidity_provider.PeginConfiguration] {
	return entities.Signed[liquidity_provider.PeginConfiguration]{
		Value: liquidity_provider.PeginConfiguration{
			TimeForDeposit: 3600,
			CallTime:       7200,
			PenaltyFee:     entities.NewWei(5000000),
			FixedFee:       entities.NewWei(1000000),
			FeePercentage:  utils.NewBigFloat64(0.5),
			MaxValue:       entities.NewUWei(10000000000000000000),
			MinValue:       entities.NewWei(100000000000000000),
		},
		Signature: "pegin-config-sig",
		Hash:      "aabbccdd",
	}
}

func newTestPegoutConfig() entities.Signed[liquidity_provider.PegoutConfiguration] {
	return entities.Signed[liquidity_provider.PegoutConfiguration]{
		Value: liquidity_provider.PegoutConfiguration{
			TimeForDeposit:       3600,
			ExpireTime:           86400,
			PenaltyFee:           entities.NewWei(10000000),
			FixedFee:             entities.NewWei(2000000),
			FeePercentage:        utils.NewBigFloat64(0.5),
			MaxValue:             entities.NewUWei(10000000000000000000),
			MinValue:             entities.NewWei(100000000000000000),
			ExpireBlocks:         1000,
			BridgeTransactionMin: entities.NewWei(50000000000000000),
		},
		Signature: "pegout-config-sig",
		Hash:      "eeff0011",
	}
}

func newTestGeneralConfig() entities.Signed[liquidity_provider.GeneralConfiguration] {
	return entities.Signed[liquidity_provider.GeneralConfiguration]{
		Value: liquidity_provider.GeneralConfiguration{
			RskConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"1000000000000000000": 10,
				"5000000000000000000": 20,
			},
			BtcConfirmations: liquidity_provider.ConfirmationsPerAmount{
				"1000000000000000000": 2,
				"5000000000000000000": 6,
			},
			PublicLiquidityCheck: true,
		},
		Signature: "general-config-sig",
		Hash:      "22334455",
	}
}

func newTestCredentials() entities.Signed[liquidity_provider.HashedCredentials] {
	return entities.Signed[liquidity_provider.HashedCredentials]{
		Value: liquidity_provider.HashedCredentials{
			HashedUsername: "hashed_user_abc123",
			HashedPassword: "hashed_pass_def456",
			UsernameSalt:   "salt_user_111",
			PasswordSalt:   "salt_pass_222",
		},
		Signature: "creds-sig",
		Hash:      "66778899",
	}
}

func newTestTrustedAccount(address string) entities.Signed[liquidity_provider.TrustedAccountDetails] {
	namePrefix := address
	if len(namePrefix) > 8 {
		namePrefix = namePrefix[:8]
	}
	hashPrefix := address
	if len(hashPrefix) > 4 {
		hashPrefix = hashPrefix[:4]
	}

	return entities.Signed[liquidity_provider.TrustedAccountDetails]{
		Value: liquidity_provider.TrustedAccountDetails{
			Address:        address,
			Name:           "TrustedAccount_" + namePrefix,
			BtcLockingCap:  entities.NewWei(5000000000000000000),
			RbtcLockingCap: entities.NewUWei(10000000000000000000),
		},
		Signature: "trusted-sig-" + namePrefix,
		Hash:      "aabb" + hashPrefix,
	}
}

func newTestPenalizedEvent(quoteHash string) penalization.PenalizedEvent {
	return penalization.PenalizedEvent{
		LiquidityProvider: defaultOwnerAddress,
		Penalty:           entities.NewWei(5000000),
		QuoteHash:         quoteHash,
	}
}

func newTestBatchPegOut(txHash string) rootstock.BatchPegOut {
	return rootstock.BatchPegOut{
		TransactionHash:    txHash,
		BlockHash:          "0xblockhash" + txHash[:8],
		BlockNumber:        99999,
		BtcTxHash:          "btcbatch" + txHash[:8],
		ReleaseRskTxHashes: []string{"0xrelease1" + txHash[:8], "0xrelease2" + txHash[:8]},
	}
}
