package redact

import (
	"log/slog"
	"regexp"
	"strings"
)

// DefaultPlaceholder replaces the value of a denied field or a value that
// matched a sensitive pattern.
const (
	DefaultPlaceholder = "[REDACTED]"
	maskShownChars     = 4
)

// dropKeyDefaults are field names whose values must never appear in logs, not
// even partially. Matching is case-insensitive and ignores separators, so
// "privateKey", "private_key" and "private-key" all match.
var dropKeyDefaults = []string{
	"privatekey", "privkey", "seed", "seedphrase", "mnemonic",
	"recoveryphrase", "walletseed", "password", "passwd", "pwd",
	"secretkey", "xprv", "signingkey",
}

// maskKeyDefaults are field names whose values may be referenced but only with
// the last four characters visible (e.g. "****ab3f").
var maskKeyDefaults = []string{
	"apikey", "apisecret", "clientsecret", "token", "accesstoken",
	"refreshtoken", "secret", "authorization", "auth", "sessionid",
	"cookie", "bearer",
}

// mnemonicPattern matches a BIP-39-looking phrase: twelve or more lowercase
// words separated by single spaces.
var mnemonicPattern = regexp.MustCompile(`^([a-z]+ ){11,}[a-z]+$`)

// bearerPattern matches an HTTP bearer credential so the token can be stripped
// while keeping the "Bearer" marker.
var bearerPattern = regexp.MustCompile(`(?i)(bearer\s+)[A-Za-z0-9._~+/=-]{8,}`)

// Config controls censorship of sensitive data. The zero value keeps redaction
// enabled with the built-in key denylists and value scanning.
type Config struct {
	// Disabled turns off all censorship.
	Disabled bool
	// DisableValueScan turns off value-pattern scanning (mnemonics, bearer
	// tokens) while keeping the key denylists.
	DisableValueScan bool
	// ExtraDropKeys are additional field names to fully redact, on top of the
	// built-in defaults.
	ExtraDropKeys []string
	// ExtraMaskKeys are additional field names to mask (last four characters),
	// on top of the built-in defaults.
	ExtraMaskKeys []string
	// Placeholder overrides the redaction marker. Defaults to "[REDACTED]".
	Placeholder string
}

// withDefaults returns a copy of the config with the placeholder defaulted.
func (c Config) withDefaults() Config {
	if c.Placeholder == "" {
		c.Placeholder = DefaultPlaceholder
	}
	return c
}

// Redactor applies the configured censorship rules to individual attributes.
type Redactor struct {
	cfg      Config
	dropKeys map[string]struct{}
	maskKeys map[string]struct{}
}

// New builds a Redactor, merging the default denylists with any extra keys
// supplied in the config.
func New(cfg Config) *Redactor {
	cfg = cfg.withDefaults()
	r := &Redactor{
		cfg:      cfg,
		dropKeys: make(map[string]struct{}),
		maskKeys: make(map[string]struct{}),
	}
	for _, k := range dropKeyDefaults {
		r.dropKeys[k] = struct{}{}
	}
	for _, k := range maskKeyDefaults {
		r.maskKeys[k] = struct{}{}
	}
	for _, k := range cfg.ExtraDropKeys {
		r.dropKeys[normalizeKey(k)] = struct{}{}
	}
	for _, k := range cfg.ExtraMaskKeys {
		r.maskKeys[normalizeKey(k)] = struct{}{}
	}
	return r
}

// Apply censors a single attribute. It first resolves any LogValuer (so typed
// wrappers like Secret take effect), then applies key-based rules and value
// scanning. Map and slice values are censored recursively so that logging a
// decoded map does not leak nested secrets; arbitrary structs passed as a
// single value are not field-walked — type sensitive fields as Secret/Masked
// (or log flat attributes) so encoding cannot leak them.
func (r *Redactor) Apply(_ []string, attr slog.Attr) slog.Attr {
	attr.Value = attr.Value.Resolve()
	if r.cfg.Disabled {
		return attr
	}

	key := normalizeKey(attr.Key)
	if _, drop := r.dropKeys[key]; drop {
		return slog.String(attr.Key, r.cfg.Placeholder)
	}
	if _, mask := r.maskKeys[key]; mask {
		return slog.String(attr.Key, r.maskValue(attr.Value))
	}
	return r.scanAttr(attr)
}

// scanAttr censors an attribute by inspecting its value: strings are scanned
// for sensitive patterns, and map/slice values are censored recursively. Other
// kinds carry no textual secret and pass through unchanged.
func (r *Redactor) scanAttr(attr slog.Attr) slog.Attr {
	switch attr.Value.Kind() {
	case slog.KindString:
		if !r.cfg.DisableValueScan {
			if scanned, changed := r.scanString(attr.Value.String()); changed {
				return slog.String(attr.Key, scanned)
			}
		}
	case slog.KindAny:
		return r.recurseMapSlice(attr)
	default:
		// Numeric, bool, time and duration values carry no textual secrets.
		return attr
	}
	return attr
}

func (r *Redactor) recurseMapSlice(attr slog.Attr) slog.Attr {
	switch inner := attr.Value.Any().(type) {
	case map[string]any:
		return slog.Any(attr.Key, r.censorValue(inner))
	case []any:
		return slog.Any(attr.Key, r.censorValue(inner))
	default:
		return attr
	}
}

// censorKeyed applies the key-based rules (drop, mask) and then value scanning
// to a key/value pair drawn from a nested map.
func (r *Redactor) censorKeyed(key string, value any) any {
	normalizedKey := normalizeKey(key)
	if _, drop := r.dropKeys[normalizedKey]; drop {
		return r.cfg.Placeholder
	}
	if _, mask := r.maskKeys[normalizedKey]; mask {
		if s, ok := value.(string); ok {
			return r.maskString(s)
		}
		return r.cfg.Placeholder
	}
	return r.censorValue(value)
}

// censorValue recursively censors a decoded value (string, map or slice),
// returning any other kind unchanged.
func (r *Redactor) censorValue(value any) any {
	switch t := value.(type) {
	case string:
		if !r.cfg.DisableValueScan {
			if scanned, changed := r.scanString(t); changed {
				return scanned
			}
		}
		return t
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[k] = r.censorKeyed(k, val)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, val := range t {
			out[i] = r.censorValue(val)
		}
		return out
	default:
		return value
	}
}

// maskValue renders a slog value showing at most its last four characters.
func (r *Redactor) maskValue(v slog.Value) string {
	return r.maskString(v.String())
}

func (r *Redactor) maskString(s string) string {
	if len(s) <= maskShownChars {
		return r.cfg.Placeholder
	}
	return "****" + s[len(s)-maskShownChars:]
}

// scanString reports whether s matched a sensitive pattern, returning the
// censored replacement when it did.
func (r *Redactor) scanString(s string) (string, bool) {
	if mnemonicPattern.MatchString(strings.TrimSpace(s)) {
		return r.cfg.Placeholder, true
	}
	if bearerPattern.MatchString(s) {
		return bearerPattern.ReplaceAllString(s, "${1}"+r.cfg.Placeholder), true
	}
	return s, false
}

// normalizeKey lowercases a field name and strips non-alphanumeric characters
// so that different casings and separators map to the same denylist entry.
func normalizeKey(key string) string {
	var b strings.Builder
	b.Grow(len(key))
	for _, c := range key {
		switch {
		case c >= 'A' && c <= 'Z':
			b.WriteRune(c + ('a' - 'A'))
		case (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			b.WriteRune(c)
		}
	}
	return b.String()
}
