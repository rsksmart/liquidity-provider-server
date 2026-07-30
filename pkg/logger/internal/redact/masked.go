package redact

import (
	"encoding/json"
	"log/slog"
)

// Masked is a defined string type for values (API key, token) that may be
// referenced in logs but only with their last maskShownChars visible. Prefer
// typing domain fields as Masked so accidental logging of a containing struct
// stays partially masked.
type Masked string

// NewMasked wraps value as a Masked value.
func NewMasked(value string) Masked { return Masked(value) }

// maskedString returns the display form: placeholder for short values, otherwise
// "****" plus the last maskShownChars characters.
func (m Masked) maskedString() string {
	if len(m) <= maskShownChars {
		return DefaultPlaceholder
	}
	return "****" + string(m[len(m)-maskShownChars:])
}

// LogValue implements slog.LogValuer, exposing at most the last four characters.
func (m Masked) LogValue() slog.Value {
	return slog.StringValue(m.maskedString())
}

// String mirrors LogValue for fmt verbs.
func (m Masked) String() string { return m.maskedString() }

// GoString hides the full value from %#v.
func (m Masked) GoString() string {
	return `Masked("` + m.maskedString() + `")`
}

// Reveal returns the underlying value.
func (m Masked) Reveal() string { return string(m) }

// MarshalText implements encoding.TextMarshaler so text encodings never leak
// the full value.
func (m Masked) MarshalText() ([]byte, error) {
	return []byte(m.maskedString()), nil
}

// MarshalJSON implements json.Marshaler so JSON encodings never leak the full
// value.
func (m Masked) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.maskedString())
}
