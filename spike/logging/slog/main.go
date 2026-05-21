// Package main is the FLY-2308 spike demo for log/slog. See
// ../README.md and docs/spikes/logs.md for context.
//
// Run: LOG_FORMAT=json LOG_LEVEL=debug go run ./spike/logging/slog
package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
)

func main() {
	cfg = readConfig()
	initLogger(os.Stdout, cfg)
	ctx := contextWithTrace(context.Background(), newTraceContext())

	// 1. LOG_FORMAT supporting json, logfmt, and otel where supported.
	demoLogFormatSupport()

	// 2. LOG_LEVEL configurable at startup.
	demoLogLevelConfiguredAtStartup(cfg)

	// 3. No plain unstructured output outside local development.
	demoStructuredOutputOnly(ctx, cfg)

	// 4. Required fields on every log line:
	//    timestamp, level, service, environment, version, traceId, message.
	demoRequiredFields(ctx)

	// Extra LPS example: imports project packages and logs real domain fields.
	demoRealLPSUseCaseLogging(ctx)

	// 5-6 and 10 use the structured request middleware below.
	inbound := httptest.NewServer(traceMiddleware(http.HandlerFunc(sampleHandler)))
	defer inbound.Close()

	// 5. Inbound W3C traceparent extraction, with a generated trace if absent.
	demoInboundTraceExtraction(inbound.URL + "/quote")

	// 6. traceId carried through context.Context and emitted on request logs.
	demoRequestScopedTrace(inbound.URL + "/quote")

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Received-Traceparent", r.Header.Get("traceparent"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer upstream.Close()

	// 7. Outbound HTTP/RPC traceparent injection with a fresh span-id.
	demoOutboundTraceInjection(upstream.URL + "/lbc/getQuote")

	// 8. Error logs containing errorCode, errorMessage, errorStack,
	//    and errorContext.
	demoErrorFields(ctx)

	// Extra LPS example: sentinel + WrapUseCaseError migrated to the
	// structured error shape, with errors.Is still matching.
	demoRealLPSUseCaseError(ctx)

	// 9. PII drop or redaction by key for the required deny-list.
	demoPIIRedaction(ctx)

	// 10. Structured replacement for gorilla/handlers.LoggingHandler.
	demoStructuredRequestLoggingReplacement(inbound.URL + "/quote")

	// Extra task deliverable: fatal-equivalent behavior.
	demoFatal(ctx)
}
