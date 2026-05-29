package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/assets"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
)

const (
	nonceBytes = 32
)

type GetManagementUiDataUseCase interface {
	Run(ctx context.Context, loggedIn bool) (*liquidity_provider.ManagementTemplate, error)
}

// NewManagementInterfaceHandler
// @Title Management Interface
// @Description Serves the static site for the Management UI
// @Success 200 object
// @Route /management [get]
func NewManagementInterfaceHandler(env environment.ManagementEnv, store sessions.Store, useCase GetManagementUiDataUseCase, templateOverride ...liquidity_provider.ManagementTemplateId) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const errorGeneratingTemplate = "Error generating template: %v"
		session, err := store.Get(req, cookies.ManagementSessionCookieName)
		loggedIn := err == nil && !session.IsNew
		result, err := useCase.Run(req.Context(), loggedIn)
		if err != nil {
			log.Errorf(errorGeneratingTemplate, err)
			sendErrorTemplate(w)
			return
		}

		templateName := result.Name
		if len(templateOverride) > 0 && loggedIn {
			templateName = templateOverride[0]
		}

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
		tmpl := template.Must(template.ParseFS(assets.TemplateFileSystem, string(templateName)))

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
			log.Errorf("Error sending %s template to client, a partial version of the template might been sent: %s", templateName, err.Error())
		}
	}
}

func htmlTemplateSecurityHeaders(w http.ResponseWriter, nonce string) {
	cspHeader := fmt.Sprintf("default-src 'self'; font-src 'self' data:; style-src 'self' 'sha256-yr5DcAJJmu0m4Rv1KfUyA8AJj1t0kAJ1D2JuSBIT1DU='; object-src 'none'; frame-src 'self'; script-src 'self' 'nonce-%s'; img-src 'self' data:; connect-src 'self';", nonce)
	w.Header().Set("Content-Security-Policy", cspHeader)
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func sendErrorTemplate(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFS(assets.TemplateFileSystem, string(liquidity_provider.ManagementErrorTemplate)))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Errorf("Error sending %s template to client, a partial version of the template might been sent: %s", liquidity_provider.ManagementErrorTemplate, err.Error())
	}
}
