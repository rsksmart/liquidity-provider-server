package redact

import (
	"encoding/json"
	"log/slog"
)

// Secret is a defined string type for highly sensitive values (private key,
// seed, mnemonic). It always renders as the redaction placeholder when logged,
// formatted, or JSON-marshaled. The real value is only accessible via Reveal
// (or an explicit string conversion). Prefer typing domain fields as Secret so
// accidental logging of a containing struct stays safe.
type Secret string

// NewSecret wraps value as a Secret.
func NewSecret(value string) Secret { return Secret(value) }

// LogValue implements slog.LogValuer and always returns the redaction marker.
func (s Secret) LogValue() slog.Value { return slog.StringValue(DefaultPlaceholder) }

// String ensures the value is also hidden from fmt verbs.
func (s Secret) String() string { return DefaultPlaceholder }

// GoString hides the value from %#v.
func (s Secret) GoString() string { return `Secret("` + DefaultPlaceholder + `")` }

// Reveal returns the underlying secret value. Use it only where the raw value
// is genuinely needed; never pass the result to a logger.
func (s Secret) Reveal() string { return string(s) }

// MarshalText implements encoding.TextMarshaler so text encodings never leak.
func (s Secret) MarshalText() ([]byte, error) {
	return []byte(DefaultPlaceholder), nil
}

// MarshalJSON implements json.Marshaler so JSON encodings never leak.
func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(DefaultPlaceholder)
}
