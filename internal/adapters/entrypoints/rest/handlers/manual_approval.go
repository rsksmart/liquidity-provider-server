package handlers

import (
	"encoding/hex"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/assets"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

// NewManualApprovalPageHandler serves the manual approval HTML page
// @Title Manual Approval Page
// @Description Serves the manual approval interface for reviewing big transactions
// @Success 200 object
// @Route /management/manual-approval [get]
func NewManualApprovalPageHandler(env environment.ManagementEnv, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const errorGeneratingTemplate = "Error generating manual approval template: %v"
		
		// Generate nonce for CSP
		bytes, err := utils.GetRandomBytes(nonceBytes)
		if err != nil {
			log.Errorf(errorGeneratingTemplate, err)
			sendErrorTemplate(w)
			return
		}
		nonce := hex.EncodeToString(bytes)
		
		w.Header().Set(rest.HeaderContentType, "text/html")
		if env.EnableSecurityHeaders {
			htmlTemplateSecurityHeaders(w, nonce)
		}
		
		// Parse the manual-approval.html template
		tmpl := template.Must(template.ParseFS(assets.TemplateFileSystem, "manual-approval.html"))
		
		// Execute template with CSRF token and nonce
		err = tmpl.Execute(w, struct {
			CSRFToken   string
			ScriptNonce string
		}{
			CSRFToken:   csrf.Token(req),
			ScriptNonce: nonce,
		})
		
		if err != nil {
			log.Errorf("Error sending manual-approval.html template to client: %s", err.Error())
		}
	}
}
