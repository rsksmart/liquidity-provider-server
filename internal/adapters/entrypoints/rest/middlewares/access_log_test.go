package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type messageOnlyFormatter struct{}

func (messageOnlyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func newCapturingLogger() (*logrus.Logger, *bytes.Buffer) {
	logger := logrus.New()
	buffer := &bytes.Buffer{}
	logger.SetOutput(buffer)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(messageOnlyFormatter{})
	return logger, buffer
}

func TestNewAccessLogMiddleware_LogsCommonLogFormat(t *testing.T) {
	logger, buffer := newCapturingLogger()
	handler := middlewares.NewAccessLogMiddleware(logger, logrus.InfoLevel)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("hello"))
		assert.NoError(t, err)
	}))

	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	line := buffer.String()
	assert.Regexp(t, `^192\.0\.2\.1 - - \[\d{2}/[A-Za-z]{3}/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4}\] "GET /path HTTP/1\.1" 200 5\n$`, line)
}

func TestNewAccessLogMiddleware_CapturesStatusAndSize(t *testing.T) {
	logger, buffer := newCapturingLogger()
	handler := middlewares.NewAccessLogMiddleware(logger, logrus.InfoLevel)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte("nope"))
		assert.NoError(t, err)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/missing", nil))

	assert.Contains(t, buffer.String(), `"GET /missing HTTP/1.1" 404 4`)
}

func TestNewAccessLogMiddleware_EscapesInjectedRequestTarget(t *testing.T) {
	logger, buffer := newCapturingLogger()
	handler := middlewares.NewAccessLogMiddleware(logger, logrus.InfoLevel)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	req.RequestURI = "/evil\"\nforged log line"
	handler.ServeHTTP(httptest.NewRecorder(), req)

	line := buffer.String()
	logged := strings.TrimSuffix(line, "\n")
	assert.NotContains(t, logged, "\n")
	assert.Contains(t, line, strconv.QuoteToASCII("GET /evil\"\nforged log line HTTP/1.1"))
}

func TestNewAccessLogMiddleware_EscapesInjectedUsername(t *testing.T) {
	logger, buffer := newCapturingLogger()
	handler := middlewares.NewAccessLogMiddleware(logger, logrus.InfoLevel)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	req.URL.User = url.User("evil\"\nforged log line")
	handler.ServeHTTP(httptest.NewRecorder(), req)

	line := buffer.String()
	logged := strings.TrimSuffix(line, "\n")
	assert.NotContains(t, logged, "\n")
	assert.Contains(t, line, strconv.Quote("evil\"\nforged log line"))
}

func TestNewAccessLogMiddleware_UnwrapExposesUnderlyingCapabilities(t *testing.T) {
	logger, _ := newCapturingLogger()
	var flushErr error
	handler := middlewares.NewAccessLogMiddleware(logger, logrus.InfoLevel)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		flushErr = http.NewResponseController(w).Flush()
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/path", nil))

	require.NoError(t, flushErr)
}
