// PII deny-list redaction through HandlerOptions.ReplaceAttr, including nested
// groups. Uses a redact-not-drop policy per the Base Logging Standards.
package main

import (
	"context"
	"log/slog"
)

func demoPIIRedaction(ctx context.Context) {
	slog.WarnContext(ctx, "attempted to log PII (expect redaction)",
		slog.String("privateKey", "0xdeadbeefcafebabe"),
		slog.String("apiKey", "sk-live-XYZ"),
		slog.String("authorization", "Bearer eyJhbGciOi..."),
		slog.String("safe", "ok"),
	)
}

// piiDenyList is the deny-list from the Base Logging Standards. Keys
// are matched case-insensitively against attribute names.
var piiDenyList = map[string]struct{}{
	"privatekey":    {},
	"seed":          {},
	"mnemonic":      {},
	"apikey":        {},
	"password":      {},
	"authorization": {},
}

// replaceAttr is the single handler interception point for the slog demo.
// It first normalizes stdlib field names to the Base Logging Standards, then
// applies the PII deny-list policy.
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	a = normalizeStandardKey(groups, a)
	return piiRedactor(groups, a)
}

func normalizeStandardKey(groups []string, a slog.Attr) slog.Attr {
	if len(groups) != 0 {
		return a
	}
	switch a.Key {
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.MessageKey:
		a.Key = "message"
	}
	return a
}

// piiRedactor is the ReplaceAttr function plugged into every handler in
// initLogger. This is slog's "deny-list interception point" — the one
// place where the redaction policy lives.
//
// Decision (see logs.md "Sensitive data"): redact rather than drop, so
// the operational signal "someone tried to log a key" survives even
// when the value does not.
func piiRedactor(groups []string, a slog.Attr) slog.Attr {
	if _, deny := piiDenyList[lowerASCII(a.Key)]; deny {
		return slog.String(a.Key, "[REDACTED]")
	}
	// Recurse into groups (e.g. errorContext) so PII keys nested inside
	// structured attrs are also caught.
	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		out := make([]slog.Attr, 0, len(attrs))
		for _, child := range attrs {
			out = append(out, piiRedactor(append(groups, a.Key), child))
		}
		return slog.Attr{Key: a.Key, Value: slog.GroupValue(out...)}
	}
	return a
}

// lowerASCII avoids strings.ToLower so we don't allocate for every key.
func lowerASCII(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			b := []byte(s)
			for j := i; j < len(b); j++ {
				if b[j] >= 'A' && b[j] <= 'Z' {
					b[j] += 'a' - 'A'
				}
			}
			return string(b)
		}
	}
	return s
}
