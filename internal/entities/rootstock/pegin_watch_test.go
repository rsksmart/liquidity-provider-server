package rootstock_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPegInWatch(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	watch := rootstock.NewPegInWatch("0xhash", 3, 101, "0xrsk", "0xregistrant", [32]byte{1})
	after := time.Now().UTC().Add(time.Second)

	assert.Equal(t, "0xhash", watch.TxHash)
	assert.Equal(t, uint(3), watch.LogIndex)
	assert.Equal(t, uint64(101), watch.BlockNumber)
	assert.Equal(t, "0xrsk", watch.RskAddress)
	assert.Equal(t, "0xregistrant", watch.Registrant)
	assert.Equal(t, [32]byte{1}, watch.RegistrationRoot)
	assert.Equal(t, rootstock.PegInWatchDiscovered, watch.State)
	assert.True(t, !watch.CreatedAt.Before(before) && !watch.CreatedAt.After(after))
	assert.Equal(t, watch.CreatedAt, watch.UpdatedAt)
}

func TestPegInWatch_SetEncoding_LeavesSupportedEncodingDiscovered(t *testing.T) {
	watch := rootstock.NewPegInWatch("0xhash", 1, 1, "0xrsk", "0xregistrant", [32]byte{})
	watch.LastError = "previous"
	watch.SetEncoding(rootstock.PegInAddressRegistryEncodingBase58)

	assert.Equal(t, uint8(rootstock.PegInAddressRegistryEncodingBase58), watch.Encoding)
	assert.Equal(t, rootstock.PegInWatchDiscovered, watch.State)
	assert.Equal(t, "previous", watch.LastError)
}

func TestPegInWatch_SetEncoding_MarksUnsupportedAndSetsLastError(t *testing.T) {
	watch := rootstock.NewPegInWatch("0xhash", 1, 1, "0xrsk", "0xregistrant", [32]byte{})
	watch.LastError = "previous"
	watch.SetEncoding(rootstock.PegInAddressRegistryEncodingBech32)

	assert.Equal(t, uint8(rootstock.PegInAddressRegistryEncodingBech32), watch.Encoding)
	assert.Equal(t, rootstock.PegInWatchUnsupportedEncoding, watch.State)
	assert.NotEmpty(t, watch.LastError)
	assert.NotEqual(t, "previous", watch.LastError)
	assert.True(t, watch.UpdatedAt.After(watch.CreatedAt) || watch.UpdatedAt.Equal(watch.CreatedAt))
}

func TestIsSupportedPegInEncoding(t *testing.T) {
	assert.True(t, rootstock.IsSupportedPegInEncoding(rootstock.PegInAddressRegistryEncodingBase58))
	assert.False(t, rootstock.IsSupportedPegInEncoding(rootstock.PegInAddressRegistryEncodingBech32))
	assert.False(t, rootstock.IsSupportedPegInEncoding(rootstock.PegInAddressRegistryEncodingBech32M))
}

func TestPegInWatch_SameLog(t *testing.T) {
	watch := rootstock.NewPegInWatch("0xhash", 3, 1, "0xrsk", "0xregistrant", [32]byte{})

	assert.True(t, watch.SameLog("0xhash", 3))
	assert.False(t, watch.SameLog("0xother", 3))
	assert.False(t, watch.SameLog("0xhash", 2))
}

func TestPegInWatch_MarkImported(t *testing.T) {
	watch := rootstock.NewPegInWatch("0xhash", 1, 1, "0xrsk", "0xregistrant", [32]byte{})
	watch.LastError = "previous"
	watch.MarkImported()

	assert.Equal(t, rootstock.PegInWatchImported, watch.State)
	assert.Empty(t, watch.LastError)
}

func TestPegInWatch_RecordError(t *testing.T) {
	watch := rootstock.NewPegInWatch("0xhash", 1, 1, "0xrsk", "0xregistrant", [32]byte{})
	first := errors.New("boom")
	require.True(t, watch.RecordError(first))
	assert.Equal(t, "boom", watch.LastError)
	require.False(t, watch.RecordError(errors.New("boom")))
	require.True(t, watch.RecordError(errors.New("other")))
	assert.Equal(t, "other", watch.LastError)
}
