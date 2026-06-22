package middlewares

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const commonLogTimeFormat = "02/Jan/2006:15:04:05 -0700"

func NewAccessLogMiddleware(logger *log.Logger, level log.Level) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &accessLogResponseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)
			logger.Log(level, buildAccessLogLine(r, start, recorder.status, recorder.size))
		})
	}
}

type accessLogResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *accessLogResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *accessLogResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *accessLogResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func buildAccessLogLine(r *http.Request, ts time.Time, status, size int) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	username := "-"
	if r.URL.User != nil {
		if name := r.URL.User.Username(); name != "" {
			username = name
		}
	}
	uri := r.RequestURI
	// HTTP/2 CONNECT requests carry the authority in r.Host rather than the request target.
	if r.ProtoMajor == 2 && r.Method == http.MethodConnect {
		uri = r.Host
	}
	if uri == "" {
		uri = r.URL.RequestURI()
	}
	// QuoteToASCII both wraps the request line in the double quotes CLF expects and escapes any
	// control characters or non-ASCII bytes a client smuggles into the URI, so a crafted request
	// target cannot forge log lines or inject terminal escape sequences.
	// TODO(slog): once the logrus->slog migration lands, drop this manual escaping and emit the
	// fields as structured attributes (e.g. slog.String("uri", uri)); the slog TextHandler quotes
	// values with the same strconv.Quote primitive.
	requestLine := strconv.QuoteToASCII(r.Method + " " + uri + " " + r.Proto)

	var builder strings.Builder
	builder.Grow(len(host) + len(username) + len(requestLine) + 40)
	builder.WriteString(host)
	builder.WriteString(" - ")
	builder.WriteString(username)
	builder.WriteString(" [")
	builder.WriteString(ts.Format(commonLogTimeFormat))
	builder.WriteString("] ")
	builder.WriteString(requestLine)
	builder.WriteByte(' ')
	builder.WriteString(strconv.Itoa(status))
	builder.WriteByte(' ')
	builder.WriteString(strconv.Itoa(size))
	return builder.String()
}
