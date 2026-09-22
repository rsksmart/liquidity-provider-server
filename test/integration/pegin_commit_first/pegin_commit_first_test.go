package pegin_commit_first

import (
	"encoding/hex"
	"math/big"
	"sort"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/integration"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommitFirstSuite struct {
	suite.Suite
	btc         *regtest.BitcoinStack
	rsk         *regtest.RootstockStack
	mongo       *regtest.MongoStack
	user        common.Address
	btcAddr     string
	depositTxID string
}

type importedDeposit struct {
	user        common.Address
	btcAddr     string
	depositTxID string
	amount      *entities.Wei
	fee         *entities.Wei
	net         *entities.Wei
	before      *big.Int
}

func TestCommitFirstSuite(t *testing.T) {
	suite.Run(t, new(CommitFirstSuite))
}

func (s *CommitFirstSuite) SetupSuite() {
	t := s.T()
	cfg := integration.ReadTestConfig(t)
	s.btc = regtest.OpenBitcoin(t)
	s.rsk = regtest.OpenRootstock(t, cfg)
	s.mongo = regtest.OpenMongo(t)
	s.rsk.WaitBridgeBtcAtLeast(t, s.btc.Height(t))
}

func (s *CommitFirstSuite) mineUntilBridgeSeesTip() {
	t := s.T()
	s.btc.Mine(t, 1)
	s.rsk.WaitBridgeBtcAtLeast(t, s.btc.Height(t))
}

func (s *CommitFirstSuite) importAtOneConfirmation() importedDeposit {
	t := s.T()
	user := s.rsk.NewUser(t)
	btcAddr := s.rsk.GetPegInAddress(t, user)
	amount := entities.NewBigWei(regtest.DepositAmountWei())
	require.Equal(t, uint64(2), s.rsk.RequiredConfirmations(t, amount))
	fee := s.rsk.CalculatePegInFee(t, amount)
	net := regtest.PayableNet(t, amount, fee)
	before := s.rsk.Balance(t, user)
	depositTxID := s.btc.PayAndMine(t, btcAddr, 1)
	s.rsk.RegisterAddress(t, user, s.btc.Rpc, depositTxID)
	s.rsk.Advance(t, regtest.FinalityDepth)
	watch := s.mongo.WaitWatchImported(t, user)
	require.Equal(t, btcAddr, watch.BtcAddress)
	return importedDeposit{
		user:        user,
		btcAddr:     btcAddr,
		depositTxID: depositTxID,
		amount:      amount,
		fee:         fee,
		net:         net,
		before:      before,
	}
}

func (s *CommitFirstSuite) assertClaimOnChain(t *testing.T, deposit importedDeposit, claim *rootstock.PegInClaim) {
	t.Helper()
	event := s.rsk.UnpackPegInRequested(t, claim.TxHash)
	require.NotEqual(t, [32]byte{}, event.PegInId)
	require.Equal(t, hex.EncodeToString(event.PegInId[:]), claim.PegInID)
	require.Equal(t, common.HexToAddress(regtest.LPSHotWalletHex), common.HexToAddress(event.Claimer))
	require.Equal(t, deposit.user, common.HexToAddress(event.RskAddress))
	require.Equal(t, 0, deposit.amount.Cmp(event.Amount))
	require.Equal(t, 0, deposit.net.Cmp(event.NetToUser))
	require.True(t, event.CallSuccess)
}

func (s *CommitFirstSuite) assertRegistrationRoot(t *testing.T) {
	t.Helper()
	head := s.rsk.Head(t)
	require.GreaterOrEqual(t, head, uint64(regtest.FinalityDepth))
	probe := head - regtest.FinalityDepth
	chain := s.rsk.GetRegistrationRoot(t, probe)
	watches := s.mongo.ListWatches(t)
	type row struct {
		blockNumber uint64
		logIndex    uint
		txHash      string
		rskAddress  string
	}
	rows := make([]row, 0, len(watches))
	for _, watch := range watches {
		if watch.BlockNumber <= probe {
			rows = append(rows, row{
				blockNumber: watch.BlockNumber,
				logIndex:    watch.LogIndex,
				txHash:      watch.TxHash,
				rskAddress:  watch.RskAddress,
			})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].blockNumber != rows[j].blockNumber {
			return rows[i].blockNumber < rows[j].blockNumber
		}
		if rows[i].logIndex != rows[j].logIndex {
			return rows[i].logIndex < rows[j].logIndex
		}
		// TODO: evaluate the necessity of txHash and rskAddress on both test and prod.
		if rows[i].txHash != rows[j].txHash {
			return rows[i].txHash < rows[j].txHash
		}
		return rows[i].rskAddress < rows[j].rskAddress
	})
	var local [32]byte
	for _, item := range rows {
		folded, err := blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, local, item.rskAddress)
		require.NoError(t, err)
		local = folded
	}
	require.Equal(t, chain, local, "local watch fold != getRegistrationRoot at %d", probe)
}

func (s *CommitFirstSuite) waitClaimSuccess(deposit importedDeposit) {
	t := s.T()
	s.rsk.WaitBalanceIncreasedBy(t, deposit.user, deposit.before, deposit.net)
	s.mongo.WaitClaimStates(
		t,
		deposit.user,
		deposit.depositTxID,
		rootstock.PegInClaimSubmitting,
		rootstock.PegInClaimClaimed,
	)
	s.rsk.Advance(t, regtest.FinalityDepth)
	claim := s.mongo.WaitClaimed(t, deposit.user, deposit.depositTxID)
	s.assertClaimOnChain(t, deposit, claim)
	s.assertRegistrationRoot(t)
}

func (s *CommitFirstSuite) TestFirstDepositRegisterImportClaim() {
	t := s.T()
	deposit := s.importAtOneConfirmation()
	head := s.rsk.Head(t)
	require.GreaterOrEqual(t, head, uint64(regtest.FinalityDepth))
	root := s.rsk.GetRegistrationRoot(t, head-regtest.FinalityDepth)
	require.NotEqual(t, [32]byte{}, root)

	s.mineUntilBridgeSeesTip()
	s.waitClaimSuccess(deposit)

	s.user = deposit.user
	s.btcAddr = deposit.btcAddr
	s.depositTxID = deposit.depositTxID
}

func (s *CommitFirstSuite) TestSecondDepositSameWatch() {
	t := s.T()
	require.NotEqual(t, common.Address{}, s.user, "TestFirstDepositRegisterImportClaim must pass first")
	require.NotEmpty(t, s.btcAddr)
	require.NotEmpty(t, s.depositTxID)

	watch := s.mongo.GetWatch(t, s.user)
	require.NotNil(t, watch)
	require.Equal(t, rootstock.PegInWatchImported, watch.State)
	first := s.mongo.GetClaim(t, s.user, s.depositTxID)
	require.NotNil(t, first)
	require.Equal(t, rootstock.PegInClaimClaimed, first.State)
	require.NotEmpty(t, first.TxHash)

	before := s.rsk.Balance(t, s.user)
	amount := entities.NewBigWei(regtest.DepositAmountWei())
	fee := s.rsk.CalculatePegInFee(t, amount)
	net := regtest.PayableNet(t, amount, fee)
	depositTxID2 := s.btc.PayAndMine(t, s.btcAddr, 2)
	s.rsk.WaitBridgeBtcAtLeast(t, s.btc.Height(t))

	second := importedDeposit{
		user:        s.user,
		btcAddr:     s.btcAddr,
		depositTxID: depositTxID2,
		amount:      amount,
		fee:         fee,
		net:         net,
		before:      before,
	}
	s.waitClaimSuccess(second)
	require.Equal(t, 1, s.mongo.CountWatches(t, s.user))
}

func (s *CommitFirstSuite) TestBelowConfirmationsThenClaim() {
	t := s.T()
	deposit := s.importAtOneConfirmation()
	s.mongo.WaitNoPaidClaim(t, deposit.user, deposit.depositTxID)
	require.Zero(t, deposit.before.Cmp(s.rsk.Balance(t, deposit.user)))

	s.mineUntilBridgeSeesTip()
	s.waitClaimSuccess(deposit)
}

func (s *CommitFirstSuite) TestRaceLostCompetingClaimer() {
	t := s.T()
	deposit := s.importAtOneConfirmation()
	s.mineUntilBridgeSeesTip()
	s.rsk.RequestPegInAsCompetitor(t, deposit.user, s.btc.Rpc, deposit.depositTxID, deposit.net)
	s.rsk.WaitBalanceIncreasedBy(t, deposit.user, deposit.before, deposit.net)
	paid := s.rsk.Balance(t, deposit.user)
	s.mongo.WaitClaimStates(t, deposit.user, deposit.depositTxID, rootstock.PegInClaimRaceLost)
	require.Equal(t, 0, paid.Cmp(s.rsk.Balance(t, deposit.user)))
}

func (s *CommitFirstSuite) TestInsufficientLiquidityThenFund() {
	t := s.T()
	deposit := s.importAtOneConfirmation()
	s.rsk.DrainLPS(t)
	t.Cleanup(func() {
		s.rsk.FundLPS(t, entities.EtherToWei(1).AsBigInt())
	})

	s.mineUntilBridgeSeesTip()
	s.mongo.WaitNoPaidClaim(t, deposit.user, deposit.depositTxID)
	require.Zero(t, deposit.before.Cmp(s.rsk.Balance(t, deposit.user)))

	s.rsk.FundLPS(t, entities.EtherToWei(1).AsBigInt())
	s.waitClaimSuccess(deposit)
}

func (s *CommitFirstSuite) TestSubmittingEmptyTxHashDoesNotResubmit() {
	t := s.T()
	deposit := s.importAtOneConfirmation()
	logBefore := regtest.LPSLogLen(t)
	s.mongo.InsertSubmittingEmptyTxHash(t, deposit.user, deposit.depositTxID, deposit.btcAddr)
	s.mineUntilBridgeSeesTip()
	s.rsk.WaitBalanceUnchanged(t, deposit.user, deposit.before)
	claim := s.mongo.GetClaim(t, deposit.user, deposit.depositTxID)
	require.NotNil(t, claim)
	require.Equal(t, rootstock.PegInClaimSubmitting, claim.State)
	require.Empty(t, claim.TxHash)
	regtest.WaitSubmittingEmptyTxHashLog(t, logBefore, deposit.user, deposit.depositTxID)
}

func (s *CommitFirstSuite) TestBelowFlyoverMinDoesNotClaim() {
	t := s.T()
	user := s.rsk.NewUser(t)
	btcAddr := s.rsk.GetPegInAddress(t, user)
	amount := entities.SatoshiToWei(100_000)
	min := s.rsk.MinAmount(t)
	require.Negative(t, amount.Cmp(min), "fixture 0.001 BTC must sit below on-chain minAmount %s", min.AsBigInt())
	before := s.rsk.Balance(t, user)
	params, err := s.btc.Env.GetNetworkParams()
	require.NoError(t, err)
	decoded, err := btcutil.DecodeAddress(btcAddr, params)
	require.NoError(t, err)
	pay, err := btcutil.NewAmount(0.001)
	require.NoError(t, err)
	txid, err := s.btc.Funding.SendToAddress(decoded, pay)
	require.NoError(t, err)
	s.btc.Mine(t, 1)
	depositTxID := txid.String()
	s.rsk.RegisterAddress(t, user, s.btc.Rpc, depositTxID)
	s.rsk.Advance(t, regtest.FinalityDepth)
	watch := s.mongo.WaitWatchImported(t, user)
	require.Equal(t, btcAddr, watch.BtcAddress)
	s.mineUntilBridgeSeesTip()
	s.mongo.WaitNoPaidClaim(t, user, depositTxID)
	require.Zero(t, before.Cmp(s.rsk.Balance(t, user)))
	require.Nil(t, s.mongo.GetClaim(t, user, depositTxID))
}

func (s *CommitFirstSuite) TestHardPauseDoesNotClaim() {
	t := s.T()
	deposit := s.importAtOneConfirmation()
	s.rsk.SetPauseLevel(t, blockchain.PauseLevelHard, "commit-first-integration-test")
	t.Cleanup(func() {
		s.rsk.SetPauseLevel(t, blockchain.PauseLevelNone, "")
	})
	s.mineUntilBridgeSeesTip()
	s.mongo.WaitNoPaidClaim(t, deposit.user, deposit.depositTxID)
	require.Zero(t, deposit.before.Cmp(s.rsk.Balance(t, deposit.user)))
	s.rsk.SetPauseLevel(t, blockchain.PauseLevelNone, "")
	s.waitClaimSuccess(deposit)
}

func (s *CommitFirstSuite) TestNeverRegisteredDoesNotCreateWatch() {
	t := s.T()
	user := s.rsk.NewUser(t)
	btcAddr := s.rsk.GetPegInAddress(t, user)
	before := s.rsk.Balance(t, user)
	s.btc.PayAndMine(t, btcAddr, 2)
	s.rsk.WaitBridgeBtcAtLeast(t, s.btc.Height(t))
	s.mongo.WaitNoWatch(t, user)
	require.Zero(t, before.Cmp(s.rsk.Balance(t, user)))
}
