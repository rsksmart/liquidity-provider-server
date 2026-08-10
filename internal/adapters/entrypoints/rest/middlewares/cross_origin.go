package middlewares

import (
	"net/http"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
)

// NewCrossOriginProtectionMiddleware builds the management API's CSRF protection.
//
// How it works: instead of a CSRF token round-trip, it uses the stdlib
// http.CrossOriginProtection (Go 1.25+). For every state-changing (unsafe) request
// -- POST/PUT/DELETE/PATCH -- the browser is required to send either a same-origin
// Sec-Fetch-Site header or an Origin header matching the request's host. A request
// carrying a cross-site Sec-Fetch-Site, or an Origin that does not match the host, is
// rejected with 403 before reaching the handler. Safe methods (GET/HEAD/OPTIONS) and
// requests with no Origin/Sec-Fetch-Site header at all (same-origin browsers that omit
// them, or non-browser clients) are allowed through.
//
// The management UI is served from the same origin as the API, so no trusted origins are
// added: a page hosted on any other origin cannot forge an unsafe request. This matches the
// prior gorilla/csrf boundary, where the Strict CSRF cookie made cross-origin state-changing
// requests (including login) unforgeable -- but without any token plumbing.
func NewCrossOriginProtectionMiddleware() func(http.Handler) http.Handler {
	protection := http.NewCrossOriginProtection()
	protection.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		jsonErr := rest.NewErrorResponse("cross-origin request rejected", true)
		rest.JsonErrorResponse(w, http.StatusForbidden, jsonErr)
	}))
	return protection.Handler
}
