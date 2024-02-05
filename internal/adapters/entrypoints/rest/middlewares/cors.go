package middlewares

import "net/http"

func NewCorsMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Add("Access-Control-Allow-Origin", "*")
			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, X-Captcha-Token")
			headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			next.ServeHTTP(w, r)
		})
	}
}
