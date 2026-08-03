package logging

// slog handler wrappers used by Init: inject traceId from context and redact PII deny-list
// keys (value replaced, key kept). Lives in its own file so logging.go stays focused on
// the public API and startup configuration.

import (
	"context"
	"log/slog"
	"strings"
)

const redactedValue = "[REDACTED]"

// Canonical PII attribute keys for redaction and call-site migration (Phase 7).
// Matching is case-insensitive on the leaf key.
const (
	KeyPrivateKey    = "privateKey"
	KeySeed          = "seed"
	KeyMnemonic      = "mnemonic"
	KeyAPIKey        = "apiKey"
	KeyPassword      = "password"
	KeyAuthorization = "authorization"
)

var piiDenyList = map[string]struct{}{
	strings.ToLower(KeyPrivateKey):    {},
	strings.ToLower(KeySeed):          {},
	strings.ToLower(KeyMnemonic):      {},
	strings.ToLower(KeyAPIKey):        {},
	strings.ToLower(KeyPassword):      {},
	strings.ToLower(KeyAuthorization): {},
}

type contextHandler struct {
	handler slog.Handler
}

func newContextHandler(handler slog.Handler) *contextHandler {
	return &contextHandler{handler: handler}
}

func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *contextHandler) Handle(ctx context.Context, record slog.Record) error {
	if traceID := TraceIDFromContext(ctx); traceID != "" {
		record.AddAttrs(slog.String("traceId", traceID))
	}
	return h.handler.Handle(ctx, record)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newContextHandler(h.handler.WithAttrs(attrs))
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return newContextHandler(h.handler.WithGroup(name))
}

func buildHandlerOptions(level *slog.LevelVar) slog.HandlerOptions {
	return slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			return redactAttr(groups, attr)
		},
	}
}

func redactAttr(groups []string, attr slog.Attr) slog.Attr {
	if attr.Equal(slog.Attr{}) {
		return attr
	}
	if attr.Value.Kind() == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		args := make([]any, len(groupAttrs))
		for i, a := range groupAttrs {
			args[i] = redactAttr(append(groups, attr.Key), a)
		}
		return slog.Group(attr.Key, args...)
	}
	if isPIIKey(attr.Key) {
		return slog.String(attr.Key, redactedValue)
	}
	return attr
}

func isPIIKey(key string) bool {
	_, ok := piiDenyList[strings.ToLower(key)]
	return ok
}
