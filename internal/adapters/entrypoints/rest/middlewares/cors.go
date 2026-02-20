package middlewares

import (
	"net/http"
	"slices"
	"strings"
)

func NewCorsMiddleware(allowedOrigins []string) func(next http.Handler) http.Handler {
	normalizedAllowedOrigins := make([]string, len(allowedOrigins))
	for i, origin := range allowedOrigins {
		normalizedAllowedOrigins[i] = strings.ToLower(strings.TrimSpace(origin))
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			origin := strings.ToLower(r.Header.Get("Origin"))
			normalizedOrigin := strings.ToLower(strings.TrimSpace(origin))

			if slices.Contains(normalizedAllowedOrigins, normalizedOrigin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			headers.Set("Vary", "Origin")
			headers.Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, X-Captcha-Token, X-Csrf-Token")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			next.ServeHTTP(w, r)
		})
	}
}
