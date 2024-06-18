package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
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

func NewCaptchaMiddleware(captchaUrl string, captchaThreshold float32, disabled bool, captchaSecretKey string) func(http.Handler) http.Handler {
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

			validation, err := validateCaptcha(r, captchaUrl, captchaSecretKey, token)
			if err != nil {
				unexpectedCaptchaError(w, err)
				return
			} else if len(validation.ErrorCodes) != 0 {
				unexpectedCaptchaError(w, errors.New(strings.Join(validation.ErrorCodes, ", ")))
				return
			}

			validCaptcha := calculateCaptchaValidity(validation, captchaThreshold)
			if !validCaptcha {
				unexpectedCaptchaError(w, errors.New("captcha doesn't meet the required threshold"))
			}
			log.Debugf("Valid captcha solved on %s\n", validation.Hostname)
			next.ServeHTTP(w, r)
		})
	}
}

func calculateCaptchaValidity(validation captchaValidationResponse, captchaThreshold float32) bool {
	validCaptcha := validation.Success
	if validation.Score != nil { // if is v3 we also use the score
		validCaptcha = validCaptcha && *validation.Score >= captchaThreshold
	}
	return validCaptcha
}

func unexpectedCaptchaError(w http.ResponseWriter, err error) {
	details := make(rest.ErrorDetails)
	details["error"] = err.Error()
	jsonErr := rest.NewErrorResponseWithDetails("error validating captcha", details, false)
	rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
}

func validateCaptcha(r *http.Request, captchaUrl, captchaSecretKey, token string) (captchaValidationResponse, error) {
	var validation captchaValidationResponse
	form := make(url.Values)
	form.Set("secret", captchaSecretKey)
	form.Set("response", token)
	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		captchaUrl,
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return captchaValidationResponse{}, err
	}
	req.Header.Set(rest.HeaderContentType, rest.ContentTypeForm)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return captchaValidationResponse{}, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Error("Error closing response body: ", err)
		}
	}()

	err = json.NewDecoder(res.Body).Decode(&validation)
	if err != nil {
		return captchaValidationResponse{}, err
	}
	return validation, nil
}
