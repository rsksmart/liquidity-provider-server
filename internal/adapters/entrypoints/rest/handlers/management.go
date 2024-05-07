package handlers

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/server/cookies"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

const templatePath = "internal/adapters/entrypoints/rest/assets/"

// NewManagementInterfaceHandler
// @Title Management Interface
// @Description Serves the static site for the Management UI
// @Success 200 object
// @Route /management [get]
func NewManagementInterfaceHandler(store sessions.Store, useCase *liquidity_provider.GetManagementUiDataUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var templateName liquidity_provider.ManagementTemplateId
		session, err := store.Get(req, cookies.ManagementSessionCookieName)
		loggedIn := err == nil && !session.IsNew
		result, err := useCase.Run(req.Context(), loggedIn)
		if err != nil {
			templateName = liquidity_provider.ManagementErrorTemplate
		} else {
			templateName = result.Name
		}
		tmpl := template.Must(template.ParseFiles(templatePath + string(templateName)))
		err = tmpl.Execute(w, struct {
			liquidity_provider.ManagementTemplateData
			CsrfToken string
		}{
			ManagementTemplateData: result.Data,
			CsrfToken:              csrf.Token(req),
		})
		if err != nil {
			log.Errorf("Error sending %s template to client, a partial version of the template might been sent: %s", templateName, err.Error())
		}
	}
}
