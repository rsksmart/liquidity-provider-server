package redact_test

import (
	"encoding"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/redact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyDropAndMaskKeys(t *testing.T) {
	r := redact.New(redact.Config{})

	dropped := r.Apply(nil, slog.String("privateKey", "0xabc123"))
	assert.Equal(t, redact.DefaultPlaceholder, dropped.Value.String())

	masked := r.Apply(nil, slog.String("apiKey", "supersecretab3f"))
	assert.Equal(t, "****ab3f", masked.Value.String())
}

func TestApplyExtraKeysAndShortMask(t *testing.T) {
	r := redact.New(redact.Config{
		ExtraDropKeys: []string{"custom_secret"},
		ExtraMaskKeys: []string{"session_token"},
	})

	assert.Equal(t, redact.DefaultPlaceholder, r.Apply(nil, slog.String("custom_secret", "x")).Value.String())
	assert.Equal(t, "****wxyz", r.Apply(nil, slog.String("session_token", "abcdefghwxyz")).Value.String())
	assert.Equal(t, redact.DefaultPlaceholder, r.Apply(nil, slog.String("apiKey", "ab")).Value.String())
}

func TestApplyValueScan(t *testing.T) {
	r := redact.New(redact.Config{})
	mnemonic := "legal winner thank year wave sausage worth useful legal winner thank yellow"

	assert.Equal(t, redact.DefaultPlaceholder, r.Apply(nil, slog.String("note", mnemonic)).Value.String())
	assert.Equal(t, "Bearer "+redact.DefaultPlaceholder,
		r.Apply(nil, slog.String("header", "Bearer abcdef123456ghijkl")).Value.String())
}

func TestApplyDisableValueScan(t *testing.T) {
	r := redact.New(redact.Config{DisableValueScan: true})
	mnemonic := "legal winner thank year wave sausage worth useful legal winner thank yellow"

	assert.Equal(t, mnemonic, r.Apply(nil, slog.String("note", mnemonic)).Value.String())
}

func TestApplyDisabled(t *testing.T) {
	r := redact.New(redact.Config{Disabled: true})

	got := r.Apply(nil, slog.String("privateKey", "exposed"))
	assert.Equal(t, "exposed", got.Value.String())
}

func TestApplyNestedMapAndSlice(t *testing.T) {
	r := redact.New(redact.Config{})
	body := map[string]any{
		"user": "alice",
		"credentials": map[string]any{
			"privateKey": "0xabc",
			"apiKey":     "supersecretab3f",
			"token":      42, // non-string mask key
		},
		"notes": []any{"legal winner thank year wave sausage worth useful legal winner thank yellow"},
	}

	got := r.Apply(nil, slog.Any("body", body))
	out, ok := got.Value.Any().(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "alice", out["user"])

	creds, ok := out["credentials"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, redact.DefaultPlaceholder, creds["privateKey"])
	assert.Equal(t, "****ab3f", creds["apiKey"])
	assert.Equal(t, redact.DefaultPlaceholder, creds["token"])

	notes, ok := out["notes"].([]any)
	require.True(t, ok)
	assert.Equal(t, redact.DefaultPlaceholder, notes[0])
}

func TestApplyTopLevelSlice(t *testing.T) {
	r := redact.New(redact.Config{})
	items := []any{
		"legal winner thank year wave sausage worth useful legal winner thank yellow",
		"safe",
	}

	got := r.Apply(nil, slog.Any("items", items))
	out, ok := got.Value.Any().([]any)
	require.True(t, ok)
	assert.Equal(t, redact.DefaultPlaceholder, out[0])
	assert.Equal(t, "safe", out[1])
}

func TestApplyLeavesNonSensitiveKinds(t *testing.T) {
	r := redact.New(redact.Config{})

	assert.Equal(t, int64(7), r.Apply(nil, slog.Int64("n", 7)).Value.Int64())
	assert.True(t, r.Apply(nil, slog.Bool("ok", true)).Value.Bool())
}

func TestSecret(t *testing.T) {
	s := redact.NewSecret("top-secret")

	assert.Equal(t, "top-secret", s.Reveal())
	assert.Equal(t, redact.DefaultPlaceholder, s.String())
	assert.Equal(t, redact.DefaultPlaceholder, s.LogValue().String())
	assert.NotContains(t, s.GoString(), "top-secret")

	text, err := encoding.TextMarshaler(s).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, []byte(redact.DefaultPlaceholder), text)

	raw, err := json.Marshal(s)
	require.NoError(t, err)
	assert.JSONEq(t, `"`+redact.DefaultPlaceholder+`"`, string(raw))
}

func TestMasked(t *testing.T) {
	m := redact.NewMasked("abcd1234wxyz")

	assert.Equal(t, "abcd1234wxyz", m.Reveal())
	assert.Equal(t, "****wxyz", m.String())
	assert.Equal(t, "****wxyz", m.LogValue().String())
	assert.NotContains(t, m.GoString(), "abcd1234")

	short := redact.NewMasked("xy")
	assert.Equal(t, redact.DefaultPlaceholder, short.String())

	text, err := encoding.TextMarshaler(m).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, []byte("****wxyz"), text)

	raw, err := json.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, `"****wxyz"`, string(raw))
}
