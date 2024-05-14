package handlers

import (
	"encoding/hex"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/assets"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

const (
	nonceBytes = 32
)

// NewManagementInterfaceHandler
// @Title Management Interface
// @Description Serves the static site for the Management UI
// @Success 200 object
// @Route /management [get]
func NewManagementInterfaceHandler(store sessions.Store, useCase *liquidity_provider.GetManagementUiDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := store.Get(req, cookies.ManagementSessionCookieName)
		loggedIn := err == nil && !session.IsNew
		result, err := useCase.Run(req.Context(), loggedIn)
		if err != nil {
			sendErrorTemplate(w)
			return
		}

		bytes, err := utils.GetRandomBytes(nonceBytes)
		if err != nil {
			sendErrorTemplate(w)
			return
		}
		nonce := hex.EncodeToString(bytes)

		htmlTemplateSecurityHeaders(w, nonce)
		tmpl := template.Must(template.ParseFS(assets.FileSystem, string(result.Name)))

		err = tmpl.Execute(w, struct {
			liquidity_provider.ManagementTemplateData
			CsrfToken   string
			ScriptNonce string
		}{
			ManagementTemplateData: result.Data,
			CsrfToken:              csrf.Token(req),
			ScriptNonce:            nonce,
		})
		if err != nil {
			log.Errorf("Error sending %s template to client, a partial version of the template might been sent: %s", result.Name, err.Error())
		}
	}
}

func htmlTemplateSecurityHeaders(w http.ResponseWriter, nonce string) {
	cspHeader := fmt.Sprintf("default-src 'self'; font-src 'self'; style-src 'self'; object-src 'none'; frame-src 'self';script-src 'self' 'nonce-%s'; img-src 'self' data:; connect-src 'self';", nonce)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Security-Policy", cspHeader)
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func sendErrorTemplate(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFS(assets.FileSystem, string(liquidity_provider.ManagementErrorTemplate)))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Errorf("Error sending %s template to client, a partial version of the template might been sent: %s", liquidity_provider.ManagementErrorTemplate, err.Error())
	}
}
