package watcher_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	timelineAddressA = "0x00000000000000000000000000000000000000a1"
	timelineAddressB = "0x00000000000000000000000000000000000000b2"
	timelineAddressC = "0x00000000000000000000000000000000000000c3"
	timelineAddressD = "0x00000000000000000000000000000000000000d4"
)

func timelineRow(
	blockNumber uint64,
	logIndex uint,
	txHash string,
	rskAddress string,
) rootstock.PegInWatch {
	return rootstock.PegInWatch{
		BlockNumber: blockNumber,
		LogIndex:    logIndex,
		TxHash:      txHash,
		RskAddress:  rskAddress,
		State:       rootstock.PegInWatchImported,
	}
}

// foldedRoot recomputes the expected root from the addresses in the given order.
func foldedRoot(t *testing.T, addresses ...string) [32]byte {
	t.Helper()
	root := [32]byte{}
	for _, address := range addresses {
		var err error
		root, err = blockchain.FoldPegInAddressRegistryRoot(crypto.Keccak256, root, address)
		require.NoError(t, err)
	}
	return root
}

func TestNewLocalRootTimeline_EmptyRowsBuildAnEmptyTimeline(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		entries []rootstock.PegInWatch
	}{
		{name: "nil rows", entries: nil},
		{name: "no rows", entries: []rootstock.PegInWatch{}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			timeline, err := watcher.NewLocalRootTimeline(testCase.entries, 100, 104, crypto.Keccak256)

			require.NoError(t, err)
			assert.Equal(t, [32]byte{}, timeline.RootAt(104))
		})
	}
}

func TestNewLocalRootTimeline_RejectsRowBelowConfiguredStartBlock(t *testing.T) {
	t.Parallel()
	entries := []rootstock.PegInWatch{
		timelineRow(100, 0, "tx-a", timelineAddressA),
		timelineRow(99, 0, "tx-b", timelineAddressB),
	}

	timeline, err := watcher.NewLocalRootTimeline(entries, 100, 104, crypto.Keccak256)

	require.EqualError(t, err, "PegIn watch exists below configured start block 100")
	assert.Equal(t, [32]byte{}, timeline.RootAt(104))
}

func TestNewLocalRootTimeline_RejectsRowAboveCurrentHead(t *testing.T) {
	t.Parallel()
	entries := []rootstock.PegInWatch{
		timelineRow(100, 0, "tx-a", timelineAddressA),
		timelineRow(105, 0, "tx-b", timelineAddressB),
	}

	timeline, err := watcher.NewLocalRootTimeline(entries, 100, 104, crypto.Keccak256)

	require.EqualError(t, err, "PegIn watch exists above current head 104")
	assert.Equal(t, [32]byte{}, timeline.RootAt(104))
}

func TestNewLocalRootTimeline_FoldsRowsInDeterministicOrder(t *testing.T) {
	t.Parallel()
	entries := []rootstock.PegInWatch{
		timelineRow(102, 1, "tx-b", timelineAddressD),
		timelineRow(100, 0, "tx-a", timelineAddressA),
		timelineRow(102, 1, "tx-a", timelineAddressC),
		timelineRow(102, 0, "tx-z", timelineAddressB),
	}
	for index := range entries {
		entries[index].RegistrationRoot = [32]byte{byte(index + 1)}
	}

	timeline, err := watcher.NewLocalRootTimeline(entries, 100, 104, crypto.Keccak256)

	require.NoError(t, err)
	rootAtFirstBlock := foldedRoot(t, timelineAddressA)
	rootAtSecondBlock := foldedRoot(t,
		timelineAddressA,
		timelineAddressB,
		timelineAddressC,
		timelineAddressD,
	)
	assert.Equal(t, [32]byte{}, timeline.RootAt(99))
	assert.Equal(t, rootAtFirstBlock, timeline.RootAt(100))
	assert.Equal(t, rootAtFirstBlock, timeline.RootAt(101))
	assert.Equal(t, rootAtSecondBlock, timeline.RootAt(102))
	assert.Equal(t, rootAtSecondBlock, timeline.RootAt(104))
}
