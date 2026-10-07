package watcher

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type matcherContextKey struct{}

// fakeBlockMatcher answers root comparison probes from memory so the binary
// search can be tested without any RPC.
type fakeBlockMatcher struct {
	matchesBelow uint64
	errorAt      map[uint64]error
	probes       []uint64
	contextTags  []any
}

func (matcher *fakeBlockMatcher) matchesAt(ctx context.Context, height uint64) (bool, error) {
	matcher.probes = append(matcher.probes, height)
	matcher.contextTags = append(matcher.contextTags, ctx.Value(matcherContextKey{}))
	if err := matcher.errorAt[height]; err != nil {
		return false, err
	}
	return height < matcher.matchesBelow, nil
}

func TestFindFirstMismatchingBlock_AllProbesMatch(t *testing.T) {
	t.Parallel()
	matcher := &fakeBlockMatcher{matchesBelow: math.MaxUint64}

	firstMismatch, err := findFirstMismatchingBlock(
		context.Background(),
		matcher.matchesAt,
		10,
		20,
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(20), firstMismatch)
}

func TestFindFirstMismatchingBlock_AllProbesMismatch(t *testing.T) {
	t.Parallel()
	matcher := &fakeBlockMatcher{matchesBelow: 0}

	firstMismatch, err := findFirstMismatchingBlock(
		context.Background(),
		matcher.matchesAt,
		10,
		20,
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(10), firstMismatch)
}

func TestFindFirstMismatchingBlock_FindsBoundary(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		startBlock    uint64
		head          uint64
		matchesBelow  uint64
		firstMismatch uint64
	}{
		{name: "middle of range", startBlock: 10, head: 20, matchesBelow: 17, firstMismatch: 17},
		{name: "second block", startBlock: 10, head: 20, matchesBelow: 11, firstMismatch: 11},
		{name: "last block", startBlock: 10, head: 20, matchesBelow: 20, firstMismatch: 20},
		{name: "single block range", startBlock: 5, head: 6, matchesBelow: 0, firstMismatch: 5},
		{name: "wide range", startBlock: 0, head: 1_000_000, matchesBelow: 654_321, firstMismatch: 654_321},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			matcher := &fakeBlockMatcher{matchesBelow: testCase.matchesBelow}

			firstMismatch, err := findFirstMismatchingBlock(
				context.Background(),
				matcher.matchesAt,
				testCase.startBlock,
				testCase.head,
			)

			require.NoError(t, err)
			assert.Equal(t, testCase.firstMismatch, firstMismatch)
		})
	}
}

func TestFindFirstMismatchingBlock_ProbesExactMidpoints(t *testing.T) {
	t.Parallel()
	matcher := &fakeBlockMatcher{matchesBelow: 17}

	firstMismatch, err := findFirstMismatchingBlock(
		context.Background(),
		matcher.matchesAt,
		10,
		20,
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(17), firstMismatch)
	assert.Equal(t, []uint64{15, 18, 17, 16}, matcher.probes)
}

func TestFindFirstMismatchingBlock_EmptyRangeDoesNotProbe(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		startBlock uint64
		head       uint64
	}{
		{name: "start equals head", startBlock: 12, head: 12},
		{name: "start above head", startBlock: 30, head: 12},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			matcher := &fakeBlockMatcher{matchesBelow: 0}

			firstMismatch, err := findFirstMismatchingBlock(
				context.Background(),
				matcher.matchesAt,
				testCase.startBlock,
				testCase.head,
			)

			require.NoError(t, err)
			assert.Equal(t, testCase.startBlock, firstMismatch)
			assert.Empty(t, matcher.probes)
		})
	}
}

func TestFindFirstMismatchingBlock_MidpointDoesNotOverflow(t *testing.T) {
	t.Parallel()
	startBlock := uint64(math.MaxUint64 - 4)
	matcher := &fakeBlockMatcher{matchesBelow: math.MaxUint64}

	firstMismatch, err := findFirstMismatchingBlock(
		context.Background(),
		matcher.matchesAt,
		startBlock,
		math.MaxUint64,
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(math.MaxUint64), firstMismatch)
	assert.Equal(t, []uint64{math.MaxUint64 - 2, math.MaxUint64 - 1}, matcher.probes)
}

func TestFindFirstMismatchingBlock_PropagatesMatcherError(t *testing.T) {
	t.Parallel()
	probeError := errors.New("probe failed")
	matcher := &fakeBlockMatcher{
		matchesBelow: math.MaxUint64,
		errorAt:      map[uint64]error{15: probeError},
	}

	firstMismatch, err := findFirstMismatchingBlock(
		context.Background(),
		matcher.matchesAt,
		10,
		20,
	)

	require.ErrorIs(t, err, probeError)
	assert.Equal(t, uint64(0), firstMismatch)
	assert.Equal(t, []uint64{15}, matcher.probes)
}

func TestFindFirstMismatchingBlock_PassesContextToMatcher(t *testing.T) {
	t.Parallel()
	ctx := context.WithValue(context.Background(), matcherContextKey{}, "tag")
	matcher := &fakeBlockMatcher{matchesBelow: 17}

	_, err := findFirstMismatchingBlock(ctx, matcher.matchesAt, 10, 20)

	require.NoError(t, err)
	require.NotEmpty(t, matcher.contextTags)
	for _, tag := range matcher.contextTags {
		assert.Equal(t, "tag", tag)
	}
}
