package logger

import "github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/redact"

// RedactionConfig controls censorship of sensitive data. The zero value keeps
// redaction enabled with the built-in key denylists and value scanning.
//
// Censorship is applied per attribute: string values are scanned for known
// sensitive patterns, and map[string]any / []any values are censored
// recursively. Arbitrary structs passed as a single value are not field-walked;
// type sensitive fields as Secret/Masked (or log flat attributes) so encoding
// cannot leak them.
type RedactionConfig = redact.Config

// Secret is a defined string type for highly sensitive values (private key,
// seed, mnemonic). It always renders as the redaction placeholder when logged,
// formatted, or JSON-marshaled. Retrieve the raw value with Reveal (or an
// explicit string conversion) and never pass that result to a logger. Prefer
// typing domain fields as Secret so accidental logging of a containing struct
// stays safe.
type Secret = redact.Secret

// Masked is a defined string type for values (API key, token) that may be
// referenced in logs but only with their last four characters visible. Prefer
// typing domain fields as Masked so accidental logging of a containing struct
// stays partially masked.
type Masked = redact.Masked

// NewSecret wraps value as a Secret.
func NewSecret(value string) Secret { return redact.NewSecret(value) }

// NewMasked wraps value as a Masked value.
func NewMasked(value string) Masked { return redact.NewMasked(value) }
