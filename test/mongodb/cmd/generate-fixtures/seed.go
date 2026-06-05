package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/configuration/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
)

const (
	fixtureLbcAddress        = "0xaabbccdd11223344556677889900aabbccddeeff"
	fixtureLpRskAddress      = "0x1234567890abcdef1234567890abcdef12345678"
	fixtureRskRefundAddress  = "0xabcdef1234567890abcdef1234567890abcdef12"
	fixtureBtcRefundAddress  = "mzBc4XEFSdzCDcTxAgf6EZXgsZWpztRhef"
	fixtureLpBtcAddress      = "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"
	fixturePeginDepositAddr  = "2N6RWTxbUem64JxLxMJ2V43erUDfKjSKqjm"
	fixturePegoutDepositAddr = "0xdeposit1234567890abcdef1234567890abcdef"

	fixtureChainID  uint64 = 31
	fixtureGasPrice int64  = 60000000

	batchPegoutTxHash = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	pegoutDepositTx   = "9999999999999999999999999999999999999999999999999999999999999999"
)

var (
	baseTime = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	peginHashes = []string{
		"1111111111111111111111111111111111111111111111111111111111111111",
		"2222222222222222222222222222222222222222222222222222222222222222",
		"3333333333333333333333333333333333333333333333333333333333333333",
		"4444444444444444444444444444444444444444444444444444444444444444",
		"5555555555555555555555555555555555555555555555555555555555555555",
	}

	pegoutHashes = []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	}
)

func writeRepresentativeData(ctx context.Context, db *registry.Database) error {
	if err := writePeginQuotes(ctx, db, peginHashes); err != nil {
		return err
	}
	if err := writeRetainedPeginQuotes(ctx, db, peginHashes); err != nil {
		return err
	}

	if err := writePegoutQuotes(ctx, db, pegoutHashes); err != nil {
		return err
	}
	if err := writeRetainedPegoutQuotes(ctx, db, pegoutHashes); err != nil {
		return err
	}
	if err := writePegoutDeposits(ctx, db, pegoutHashes[:3]); err != nil {
		return err
	}

	if err := writeLPConfigurations(ctx, db); err != nil {
		return err
	}
	if err := writeTrustedAccounts(ctx, db); err != nil {
		return err
	}
	if err := writePenalizations(ctx, db, peginHashes[:2]); err != nil {
		return err
	}
	return writeBatchPegOutEvents(ctx, db)
}

func writePeginQuotes(ctx context.Context, db *registry.Database, hashes []string) error {
	for i, h := range hashes {
		ts := uint32(baseTime.Add(time.Duration(-i) * time.Hour).Unix())
		createdQuote := newPeginQuote(h, ts)
		if err := db.PeginRepository.InsertQuote(ctx, createdQuote); err != nil {
			return fmt.Errorf("insert pegin quote %d: %w", i, err)
		}
	}
	return nil
}

func writeRetainedPeginQuotes(ctx context.Context, db *registry.Database, hashes []string) error {
	peginStates := []quote.PeginState{
		quote.PeginStateWaitingForDeposit,
		quote.PeginStateCallForUserSucceeded,
		quote.PeginStateRegisterPegInSucceeded,
		quote.PeginStateCallForUserFailed,
		quote.PeginStateTimeForDepositElapsed,
	}
	if len(hashes) != len(peginStates) {
		return fmt.Errorf("retained pegin fixtures mismatch: %d hashes, %d states", len(hashes), len(peginStates))
	}
	for i, h := range hashes {
		retainedQuote := newRetainedPeginQuote(h, peginStates[i])
		if err := db.PeginRepository.InsertRetainedQuote(ctx, retainedQuote); err != nil {
			return fmt.Errorf("insert retained pegin quote %d: %w", i, err)
		}
	}
	return nil
}

func writePegoutQuotes(ctx context.Context, db *registry.Database, hashes []string) error {
	for i, h := range hashes {
		ts := uint32(baseTime.Add(time.Duration(-i) * time.Hour).Unix())
		createdQuote := newPegoutQuote(h, ts)
		if err := db.PegoutRepository.InsertQuote(ctx, createdQuote); err != nil {
			return fmt.Errorf("insert pegout quote %d: %w", i, err)
		}
	}
	return nil
}

func writeRetainedPegoutQuotes(ctx context.Context, db *registry.Database, hashes []string) error {
	pegoutStates := []quote.PegoutState{
		quote.PegoutStateWaitingForDeposit,
		quote.PegoutStateSendPegoutSucceeded,
		quote.PegoutStateBridgeTxSucceeded,
		quote.PegoutStateSendPegoutFailed,
		quote.PegoutStateBtcReleased,
	}
	if len(hashes) != len(pegoutStates) {
		return fmt.Errorf("retained pegout fixtures mismatch: %d hashes, %d states", len(hashes), len(pegoutStates))
	}
	for i, h := range hashes {
		retainedQuote := newRetainedPegoutQuote(h, pegoutStates[i])
		if err := db.PegoutRepository.InsertRetainedQuote(ctx, retainedQuote); err != nil {
			return fmt.Errorf("insert retained pegout quote %d: %w", i, err)
		}
	}
	return nil
}

func writePegoutDeposits(ctx context.Context, db *registry.Database, quoteHashes []string) error {
	for i, h := range quoteHashes {
		txHash := "0x" + pegoutDepositTx[:len(pegoutDepositTx)-2] + fmt.Sprintf("%02x", i+1)
		deposit := newPegoutDeposit(txHash, h)
		if err := db.PegoutRepository.UpsertPegoutDeposit(ctx, deposit); err != nil {
			return fmt.Errorf("upsert pegout deposit for quote %s: %w", h, err)
		}
	}
	return nil
}

func writeLPConfigurations(ctx context.Context, db *registry.Database) error {
	if err := db.LiquidityProviderRepository.UpsertPeginConfiguration(ctx, newPeginConfig()); err != nil {
		return fmt.Errorf("upsert pegin configuration: %w", err)
	}
	if err := db.LiquidityProviderRepository.UpsertPegoutConfiguration(ctx, newPegoutConfig()); err != nil {
		return fmt.Errorf("upsert pegout configuration: %w", err)
	}
	if err := db.LiquidityProviderRepository.UpsertGeneralConfiguration(ctx, newGeneralConfig()); err != nil {
		return fmt.Errorf("upsert general configuration: %w", err)
	}
	if err := db.LiquidityProviderRepository.UpsertCredentials(ctx, newCredentials()); err != nil {
		return fmt.Errorf("upsert credentials: %w", err)
	}
	return nil
}

func writeTrustedAccounts(ctx context.Context, db *registry.Database) error {
	addresses := []string{
		"0xaaaa111122223333444455556666777788889999",
		"0xbbbb111122223333444455556666777788889999",
	}
	for i, address := range addresses {
		if err := db.TrustedAccountRepository.AddTrustedAccount(ctx, newTrustedAccount(address)); err != nil {
			return fmt.Errorf("add trusted account %d: %w", i+1, err)
		}
	}
	return nil
}

func writePenalizations(ctx context.Context, db *registry.Database, quoteHashes []string) error {
	for _, h := range quoteHashes {
		if err := db.PenalizedEventRepository.InsertPenalization(ctx, newPenalizedEvent(h)); err != nil {
			return fmt.Errorf("insert penalization for quote %s: %w", h, err)
		}
	}
	return nil
}

func writeBatchPegOutEvents(ctx context.Context, db *registry.Database) error {
	if err := db.BatchPegOutRepository.UpsertBatch(ctx, newBatchPegOut("0x"+batchPegoutTxHash)); err != nil {
		return fmt.Errorf("upsert batch pegout: %w", err)
	}
	return nil
}

// --- fixture constructors ---

func newPeginQuote(hash string, agreementTimestamp uint32) quote.CreatedPeginQuote {
	return quote.CreatedPeginQuote{
		Hash: hash,
		Quote: quote.PeginQuote{
			FedBtcAddress:      fixturePeginDepositAddr,
			LbcAddress:         fixtureLbcAddress,
			LpRskAddress:       fixtureLpRskAddress,
			BtcRefundAddress:   fixtureBtcRefundAddress,
			RskRefundAddress:   fixtureRskRefundAddress,
			LpBtcAddress:       fixtureLpBtcAddress,
			CallFee:            entities.NewWei(1000000),
			PenaltyFee:         entities.NewWei(5000000),
			ContractAddress:    fixtureLpRskAddress,
			Data:               "0x",
			GasLimit:           100000,
			Nonce:              42,
			Value:              entities.NewWei(500000000000000000),
			AgreementTimestamp: agreementTimestamp,
			TimeForDeposit:     3600,
			LpCallTime:         7200,
			Confirmations:      10,
			CallOnRegister:     false,
			GasFee:             entities.NewWei(21000),
			ChainId:            fixtureChainID,
		},
		CreationData: quote.PeginCreationData{
			GasPrice:      entities.NewWei(fixtureGasPrice),
			FeePercentage: utils.NewBigFloat64(0.5),
			FixedFee:      entities.NewWei(1000000),
		},
	}
}

func newRetainedPeginQuote(hash string, state quote.PeginState) quote.RetainedPeginQuote {
	hashPrefix := safeSlice(hash, 0, 8)
	return quote.RetainedPeginQuote{
		QuoteHash:             hash,
		DepositAddress:        fixturePeginDepositAddr,
		Signature:             "0xsignature_pegin_" + hashPrefix,
		RequiredLiquidity:     entities.NewWei(600000000000000000),
		State:                 state,
		UserBtcTxHash:         "btctx" + hashPrefix,
		CallForUserTxHash:     "0xcalltx" + hashPrefix,
		RegisterPeginTxHash:   "0xregtx" + hashPrefix,
		CallForUserGasUsed:    50000,
		CallForUserGasPrice:   entities.NewWei(fixtureGasPrice),
		RegisterPeginGasUsed:  80000,
		RegisterPeginGasPrice: entities.NewWei(fixtureGasPrice),
		OwnerAccountAddress:   fixtureLpRskAddress,
	}
}

func newPegoutQuote(hash string, agreementTimestamp uint32) quote.CreatedPegoutQuote {
	return quote.CreatedPegoutQuote{
		Hash: hash,
		Quote: quote.PegoutQuote{
			LbcAddress:            fixtureLbcAddress,
			LpRskAddress:          fixtureLpRskAddress,
			BtcRefundAddress:      fixtureBtcRefundAddress,
			RskRefundAddress:      fixtureRskRefundAddress,
			LpBtcAddress:          fixtureLpBtcAddress,
			CallFee:               entities.NewWei(2000000),
			PenaltyFee:            entities.NewWei(10000000),
			Nonce:                 99,
			DepositAddress:        fixturePegoutDepositAddr,
			Value:                 entities.NewWei(1000000000000000000),
			AgreementTimestamp:    agreementTimestamp,
			DepositDateLimit:      agreementTimestamp + 7200,
			DepositConfirmations:  10,
			TransferConfirmations: 10,
			TransferTime:          7200,
			ExpireDate:            agreementTimestamp + 86400,
			ExpireBlock:           100000,
			GasFee:                entities.NewWei(42000),
			ChainId:               fixtureChainID,
		},
		CreationData: quote.PegoutCreationData{
			FeeRate:       utils.NewBigFloat64(0.001),
			FeePercentage: utils.NewBigFloat64(0.5),
			GasPrice:      entities.NewWei(fixtureGasPrice),
			FixedFee:      entities.NewWei(2000000),
		},
	}
}

func newRetainedPegoutQuote(hash string, state quote.PegoutState) quote.RetainedPegoutQuote {
	hashPrefix := safeSlice(hash, 0, 8)
	return quote.RetainedPegoutQuote{
		QuoteHash:            hash,
		DepositAddress:       fixturePegoutDepositAddr,
		Signature:            "0xsignature_pegout_" + hashPrefix,
		RequiredLiquidity:    entities.NewWei(1200000000000000000),
		State:                state,
		UserRskTxHash:        "0xusertx" + hashPrefix,
		LpBtcTxHash:          "btctx" + hashPrefix,
		RefundPegoutTxHash:   "0xrefundtx" + hashPrefix,
		BridgeRefundTxHash:   "0xbridgetx" + hashPrefix,
		BridgeRefundGasUsed:  60000,
		BridgeRefundGasPrice: entities.NewWei(fixtureGasPrice),
		RefundPegoutGasUsed:  70000,
		RefundPegoutGasPrice: entities.NewWei(fixtureGasPrice),
		SendPegoutBtcFee:     entities.NewWei(15000),
		BtcReleaseTxHash:     "btcrelease" + hashPrefix,
		OwnerAccountAddress:  fixtureLpRskAddress,
	}
}

func newPegoutDeposit(txHash, quoteHash string) quote.PegoutDeposit {
	return quote.PegoutDeposit{
		TxHash:      txHash,
		QuoteHash:   quoteHash,
		Amount:      entities.NewWei(1000000000000000000),
		Timestamp:   baseTime.UTC(),
		BlockNumber: 50000,
		From:        fixtureLpRskAddress,
	}
}

func newPeginConfig() entities.Signed[liquidity_provider.PeginConfiguration] {
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
		Signature: "pegin-config-sig-fixture",
		Hash:      "aabbccdd",
	}
}

func newPegoutConfig() entities.Signed[liquidity_provider.PegoutConfiguration] {
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
		Signature: "pegout-config-sig-fixture",
		Hash:      "eeff0011",
	}
}

func newGeneralConfig() entities.Signed[liquidity_provider.GeneralConfiguration] {
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
		Signature: "general-config-sig-fixture",
		Hash:      "22334455",
	}
}

func newCredentials() entities.Signed[liquidity_provider.HashedCredentials] {
	return entities.Signed[liquidity_provider.HashedCredentials]{
		Value: liquidity_provider.HashedCredentials{
			HashedUsername: "hashed_user_fixture",
			HashedPassword: "hashed_pass_fixture",
			UsernameSalt:   "salt_user_fixture",
			PasswordSalt:   "salt_pass_fixture",
		},
		Signature: "creds-sig-fixture",
		Hash:      "66778899",
	}
}

func newTrustedAccount(address string) entities.Signed[liquidity_provider.TrustedAccountDetails] {
	return entities.Signed[liquidity_provider.TrustedAccountDetails]{
		Value: liquidity_provider.TrustedAccountDetails{
			Address:        address,
			Name:           "Fixture_" + safeSlice(address, 0, 10),
			BtcLockingCap:  entities.NewWei(5000000000000000000),
			RbtcLockingCap: entities.NewUWei(10000000000000000000),
		},
		Signature: "trusted-sig-fixture-" + safeSlice(address, 0, 8),
		Hash:      "aabb" + safeSlice(address, 2, 6),
	}
}

func safeSlice(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	if start > len(s) {
		return ""
	}
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}

func newPenalizedEvent(quoteHash string) penalization.PenalizedEvent {
	return penalization.PenalizedEvent{
		LiquidityProvider: fixtureLpRskAddress,
		Penalty:           entities.NewWei(5000000),
		QuoteHash:         quoteHash,
	}
}

func newBatchPegOut(txHash string) rootstock.BatchPegOut {
	txHashPrefix := safeSlice(txHash, 0, 8)
	return rootstock.BatchPegOut{
		TransactionHash:    txHash,
		BlockHash:          "0xblockhash_fixture_" + txHashPrefix,
		BlockNumber:        99999,
		BtcTxHash:          "btcbatch_fixture_" + txHashPrefix,
		ReleaseRskTxHashes: []string{"0xrelease1_" + txHashPrefix, "0xrelease2_" + txHashPrefix},
	}
}
