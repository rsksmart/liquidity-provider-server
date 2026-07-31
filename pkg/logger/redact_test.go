package logger_test

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDenylistedKeyIsRedacted(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "leak attempt", logger.String("privateKey", "0xabc123def456"))

	assert.Equal(t, "[REDACTED]", decodeLine(t, &buf)["privateKey"])
}

func TestMaskedKeyShowsLastFour(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "api call", logger.String("apiKey", "supersecretab3f"))

	assert.Equal(t, "****ab3f", decodeLine(t, &buf)["apiKey"])
}

func TestSecretWrapperAlwaysRedacts(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	// Even under a non-denylisted key, the typed wrapper hides the value.
	log.Info(context.Background(), "wrapped", logger.Any("customField", logger.Secret("top-secret")))

	assert.Equal(t, "[REDACTED]", decodeLine(t, &buf)["customField"])
}

func TestMaskedWrapperShowsLastFour(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "wrapped", logger.Any("customToken", logger.Masked("abcd1234wxyz")))

	assert.Equal(t, "****wxyz", decodeLine(t, &buf)["customToken"])
}

func TestMnemonicValueIsRedacted(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	mnemonic := "legal winner thank year wave sausage worth useful legal winner thank yellow"
	log.Info(context.Background(), "note", logger.String("note", mnemonic))

	assert.Equal(t, "[REDACTED]", decodeLine(t, &buf)["note"])
}

func TestBearerTokenIsStripped(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "auth", logger.String("header", "Bearer abcdef123456ghijkl"))

	assert.Equal(t, "Bearer [REDACTED]", decodeLine(t, &buf)["header"])
}

// TestDomainHashIsNotRedacted guards against over-eager censorship: a 64-hex
// transaction hash is a legitimate domain identifier and must survive.
func TestDomainHashIsNotRedacted(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	txHash := "0x" + "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"
	log.Info(context.Background(), "tx", logger.String("txHash", txHash))

	assert.Equal(t, txHash, decodeLine(t, &buf)["txHash"])
}

func TestRedactionCanBeDisabled(t *testing.T) {
	var buf bytes.Buffer
	cfg := logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Output:      &buf,
		Redaction:   logger.RedactionConfig{Disabled: true},
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)

	log.Info(context.Background(), "raw", logger.String("privateKey", "exposed"))

	assert.Equal(t, "exposed", decodeLine(t, &buf)["privateKey"])
}

// TestNestedMapValuesAreRedacted covers the M3 recursion: logging a decoded
// map (e.g. a JSON body) must still censor denylisted keys and secret-looking
// values nested inside it.
func TestNestedMapValuesAreRedacted(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	body := map[string]any{
		"user": "alice",
		"credentials": map[string]any{
			"privateKey": "0xabc123def456",
			"apiKey":     "supersecretab3f",
		},
		"notes": []any{"legal winner thank year wave sausage worth useful legal winner thank yellow"},
	}
	log.Info(context.Background(), "decoded body", logger.Any("body", body))

	m := decodeLine(t, &buf)
	got, ok := m["body"].(map[string]any)
	require.True(t, ok, "body should be a nested object")
	assert.Equal(t, "alice", got["user"])

	creds, ok := got["credentials"].(map[string]any)
	require.True(t, ok, "credentials should be a nested object")
	assert.Equal(t, "[REDACTED]", creds["privateKey"])
	assert.Equal(t, "****ab3f", creds["apiKey"])

	notes, ok := got["notes"].([]any)
	require.True(t, ok, "notes should be a nested array")
	assert.Equal(t, "[REDACTED]", notes[0])
}

func TestSecretAndMaskedHideOutsideLogger(t *testing.T) {
	secret := logger.Secret("top-secret-value")
	masked := logger.Masked("abcd1234wxyz")

	assert.Equal(t, "top-secret-value", secret.Reveal())
	assert.Equal(t, "abcd1234wxyz", masked.Reveal())
	assert.Equal(t, "[REDACTED]", secret.String())
	assert.Equal(t, "****wxyz", masked.String())
	assert.NotContains(t, secret.GoString(), "top-secret")
	assert.NotContains(t, masked.GoString(), "abcd1234")

	secretJSON, err := json.Marshal(secret)
	require.NoError(t, err)
	assert.JSONEq(t, `"[REDACTED]"`, string(secretJSON))

	maskedJSON, err := json.Marshal(masked)
	require.NoError(t, err)
	assert.JSONEq(t, `"****wxyz"`, string(maskedJSON))

	payload := struct {
		Key   logger.Secret `json:"key"`
		Token logger.Masked `json:"token"`
	}{
		Key:   logger.NewSecret("nested-secret"),
		Token: logger.NewMasked("tokenvalue99"),
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	assert.JSONEq(t, `{"key":"[REDACTED]","token":"****ue99"}`, string(body))
}

func TestStructWithSecretFieldsIsSafeWhenLogged(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	wallet := struct {
		Address    string
		PrivateKey logger.Secret
		APIKey     logger.Masked
	}{
		Address:    "0xabc",
		PrivateKey: "0xdeadbeef",
		APIKey:     "supersecretab3f",
	}
	log.Info(context.Background(), "wallet", logger.Any("wallet", wallet))

	m := decodeLine(t, &buf)
	got, ok := m["wallet"].(map[string]any)
	require.True(t, ok, "wallet should be a nested object")
	assert.Equal(t, "0xabc", got["Address"])
	assert.Equal(t, "[REDACTED]", got["PrivateKey"])
	assert.Equal(t, "****ab3f", got["APIKey"])
}

func TestExtraDropAndMaskKeys(t *testing.T) {
	var buf bytes.Buffer
	cfg := logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Output:      &buf,
		Redaction: logger.RedactionConfig{
			ExtraDropKeys: []string{"custom_secret"},
			ExtraMaskKeys: []string{"session_token"},
		},
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)

	log.Info(context.Background(), "custom keys",
		logger.String("custom_secret", "drop-me"),
		logger.String("session_token", "abcdefghwxyz"),
	)

	m := decodeLine(t, &buf)
	assert.Equal(t, "[REDACTED]", m["custom_secret"])
	assert.Equal(t, "****wxyz", m["session_token"])
}

func TestDisableValueScanKeepsMnemonicStrings(t *testing.T) {
	var buf bytes.Buffer
	mnemonic := "legal winner thank year wave sausage worth useful legal winner thank yellow"
	cfg := logger.Config{
		Service:     "svc",
		Environment: "production",
		Version:     "v1",
		Output:      &buf,
		Redaction:   logger.RedactionConfig{DisableValueScan: true},
	}.WithClock(fixedClock())
	log, err := logger.New(cfg)
	require.NoError(t, err)

	log.Info(context.Background(), "note", logger.String("note", mnemonic))

	assert.Equal(t, mnemonic, decodeLine(t, &buf)["note"])
}

func TestShortMaskedValueUsesPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	log.Info(context.Background(), "short",
		logger.String("apiKey", "ab"),
		// Non-denylisted key so only Masked's own short-value rule applies.
		logger.Any("custom", logger.Masked("xy")),
	)

	m := decodeLine(t, &buf)
	assert.Equal(t, "[REDACTED]", m["apiKey"])
	assert.Equal(t, "[REDACTED]", m["custom"])
}

func TestNestedMaskKeyNonStringUsesPlaceholder(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	body := map[string]any{
		"apiKey": 12345, // mask key with non-string value
	}
	log.Info(context.Background(), "body", logger.Any("body", body))

	m := decodeLine(t, &buf)
	got, ok := m["body"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "[REDACTED]", got["apiKey"])
}

func TestTopLevelSliceValuesAreRedacted(t *testing.T) {
	var buf bytes.Buffer
	log := newTestLogger(t, &buf, logger.FormatJSON)

	items := []any{
		"legal winner thank year wave sausage worth useful legal winner thank yellow",
		"safe",
	}
	log.Info(context.Background(), "list", logger.Any("items", items))

	m := decodeLine(t, &buf)
	got, ok := m["items"].([]any)
	require.True(t, ok)
	assert.Equal(t, "[REDACTED]", got[0])
	assert.Equal(t, "safe", got[1])
}

func TestSecretAndMaskedMarshalText(t *testing.T) {
	secret := logger.NewSecret("top-secret-value")
	masked := logger.NewMasked("abcd1234wxyz")

	var secretTM encoding.TextMarshaler = secret
	var maskedTM encoding.TextMarshaler = masked

	secretText, err := secretTM.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, []byte("[REDACTED]"), secretText)

	maskedText, err := maskedTM.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, []byte("****wxyz"), maskedText)
}
