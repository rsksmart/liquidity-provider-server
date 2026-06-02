// PII deny-list redaction through a logrus Hook that mutates Entry.Data, with
// recursion into nested Fields for errorContext-shaped data.
package main

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

func demoPIIRedaction(ctx context.Context) {
	baseEntry(ctx).WithFields(logrus.Fields{
		"privateKey":    "0xdeadbeefcafebabe",
		"apiKey":        "sk-live-XYZ",
		"authorization": "Bearer eyJhbGciOi...",
		"safe":          "ok",
	}).Warn("attempted to log PII (expect redaction)")
}

// piiDenyList — matches the slog demo so the policy is identical.
var piiDenyList = map[string]struct{}{
	"privatekey":    {},
	"seed":          {},
	"mnemonic":      {},
	"apikey":        {},
	"password":      {},
	"authorization": {},
}

// piiRedactHook implements logrus.Hook. logrus's redaction story is
// hooks that mutate Entry.Data in place. This is the "PII deny-list
// interception" row for logrus in the comparison table — usable, but
// only sees fields after they have been built into the Entry, and does
// not recurse into nested logrus.Fields the way slog's ReplaceAttr does.
type piiRedactHook struct{}

func (h *piiRedactHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *piiRedactHook) Fire(e *logrus.Entry) error {
	for k := range e.Data {
		if _, deny := piiDenyList[strings.ToLower(k)]; deny {
			e.Data[k] = "[REDACTED]"
			continue
		}
		// Recurse into nested logrus.Fields (used by errorContext).
		if nested, ok := e.Data[k].(logrus.Fields); ok {
			for nk := range nested {
				if _, deny := piiDenyList[strings.ToLower(nk)]; deny {
					nested[nk] = "[REDACTED]"
				}
			}
		}
	}
	return nil
}
