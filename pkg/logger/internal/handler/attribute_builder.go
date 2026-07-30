package handler

import (
	"log/slog"
	"time"

	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/core"
	"github.com/rsksmart/liquidity-provider-server/pkg/logger/internal/redact"
)

// attrBuilder normalises attributes: it renames the slog built-in keys to the
// canonical base-field names, lowercases the level, formats the timestamp, and
// applies redaction to every attribute.
type attrBuilder struct {
	redactor *redact.Redactor
	tz       *time.Location
}

// replaceAttr is used as slog.HandlerOptions.ReplaceAttr for the JSON and
// logfmt handlers. Built-in slog keys are remapped only at the top level;
// nested attributes go straight to redaction.
func (b *attrBuilder) replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if len(groups) > 0 {
		return b.redactor.Apply(groups, a)
	}
	switch a.Key {
	case slog.TimeKey:
		if a.Value.Kind() == slog.KindTime {
			return slog.String(string(core.FieldTimestamp), a.Value.Time().In(b.tz).Format(core.TimestampLayout))
		}
	case slog.LevelKey:
		if lvl, ok := a.Value.Any().(slog.Level); ok {
			return slog.String(string(core.FieldLevel), core.Level(lvl).String())
		}
	case slog.MessageKey:
		return slog.String(string(core.FieldMessage), a.Value.String())
	}
	return b.redactor.Apply(groups, a)
}
