package middlewares

import (
	"encoding/json"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type captchaValidationResponse struct {
	Success     bool      `json:"success"`
	Score       *float32  `json:"score"`
	Action      *string   `json:"action"`
	ChallengeTs time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func NewCaptchaMiddleware(captchaThreshold float32, disabled bool, captchaSecretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if captchaThreshold < 0.5 {
			log.Warn("Too low captcha threshold value!")
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Captcha-Token")
			if disabled {
				log.Warn("IMPORTANT! Handling request with captcha validation disabled")
				next.ServeHTTP(w, r)
				return
			} else if token == "" {
				jsonErr := rest.NewErrorResponse("missing X-Captcha-Token header", true)
				rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
				return
			}

			form := make(url.Values)
			form.Set("secret", captchaSecretKey)
			form.Set("response", token)
			res, err := http.DefaultClient.PostForm("https://www.google.com/recaptcha/api/siteverify", form)
			if err != nil {
				details := make(rest.ErrorDetails)
				details["error"] = err.Error()
				jsonErr := rest.NewErrorResponseWithDetails("error validating captcha", details, false)
				rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
				return
			}

			defer func() {
				if err = res.Body.Close(); err != nil {
					log.Error("Error closing response body: ", err)
				}
			}()

			var validation captchaValidationResponse
			err = json.NewDecoder(res.Body).Decode(&validation)
			if err != nil {
				jsonErr := rest.NewErrorResponse("error validating captcha", false)
				rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
				return
			}

			validCaptcha := validation.Success
			if validation.Score != nil { // if is v3 we also use the score
				validCaptcha = validCaptcha && *validation.Score >= captchaThreshold
			}

			if validCaptcha {
				log.Debugf("Valid captcha solved on %s\n", validation.Hostname)
				next.ServeHTTP(w, r)
			} else {
				details := make(rest.ErrorDetails)
				details["errors"] = validation.ErrorCodes
				jsonErr := rest.NewErrorResponseWithDetails("error validating captcha", details, true)
				rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			}
		})
	}
}
