// PII deny-list redaction through piiWriter, a JSON regex rewrite on the
// io.Writer. Demonstrates the weak arbitrary-key story called out in logs.md.
package main

import (
	"context"
	"io"
	"regexp"
)

func demoPIIRedaction(ctx context.Context) {
	loggerWithTrace(ctx).Warn().
		Str("privateKey", "0xdeadbeefcafebabe").
		Str("apiKey", "sk-live-XYZ").
		Str("authorization", "Bearer eyJhbGciOi...").
		Str("safe", "ok").
		Msg("attempted to log PII (expect redaction)")
}

var piiDenyList = []string{
	"privateKey",
	"seed",
	"mnemonic",
	"apiKey",
	"password",
	"authorization",
}

// piiPattern matches one denied key in zerolog's JSON output:
//
//	"privateKey":"0xdeadbeef"
//
// We replace the value with "[REDACTED]". This is the structural
// downside of zerolog called out in logs.md "PII deny-list interception":
// zerolog has no per-field hook, so redaction either has to live inside
// the writer (here) or every call site must remember to use a wrapper.
// Both are weaker than slog's ReplaceAttr.
var piiPattern = func() *regexp.Regexp {
	alt := piiDenyList[0]
	for _, k := range piiDenyList[1:] {
		alt += "|" + k
	}
	return regexp.MustCompile(`"(` + alt + `)":"[^"]*"`)
}()

// piiWriter wraps zerolog's destination io.Writer and rewrites denied
// keys on the way out. JSON-aware enough for flat values; nested objects
// would need a real JSON walker.
type piiWriter struct{ Out io.Writer }

func (w piiWriter) Write(p []byte) (int, error) {
	redacted := piiPattern.ReplaceAll(p, []byte(`"$1":"[REDACTED]"`))
	if _, err := w.Out.Write(redacted); err != nil {
		return 0, err
	}
	// Report the original length so zerolog does not think it wrote less.
	return len(p), nil
}
