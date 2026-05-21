// PII deny-list redaction through a zapcore.Core wrapper (piiCore). Heavier than
// slog ReplaceAttr but workable for arbitrary key denial.
package main

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func demoPIIRedaction(ctx context.Context) {
	loggerWithTrace(ctx).Warn("attempted to log PII (expect redaction)",
		zap.String("privateKey", "0xdeadbeefcafebabe"),
		zap.String("apiKey", "sk-live-XYZ"),
		zap.String("authorization", "Bearer eyJhbGciOi..."),
		zap.String("safe", "ok"),
	)
}

var piiDenyList = map[string]struct{}{
	"privatekey":    {},
	"seed":          {},
	"mnemonic":      {},
	"apikey":        {},
	"password":      {},
	"authorization": {},
}

// piiCore wraps a zapcore.Core. We intercept Write so we can mutate
// fields before they reach the encoder. This is the "PII deny-list
// interception" row for zap — workable, but heavier than slog's
// ReplaceAttr: zap's fields are typed, and we must reset each denied
// field through zapcore.Field assignment rather than a simple value
// swap.
type piiCore struct{ zapcore.Core }

func (c piiCore) With(fields []zapcore.Field) zapcore.Core {
	return piiCore{Core: c.Core.With(redactFields(fields))}
}

func (c piiCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c piiCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.Core.Write(ent, redactFields(fields))
}

func redactFields(fields []zapcore.Field) []zapcore.Field {
	out := make([]zapcore.Field, len(fields))
	for i, f := range fields {
		if _, deny := piiDenyList[strings.ToLower(f.Key)]; deny {
			out[i] = zapcore.Field{Key: f.Key, Type: zapcore.StringType, String: "[REDACTED]"}
			continue
		}
		out[i] = f
	}
	return out
}
