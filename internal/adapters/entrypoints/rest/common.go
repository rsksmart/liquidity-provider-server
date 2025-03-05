package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

const (
	HeaderContentType = "Content-Type"
)

const (
	ContentTypeJson = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

var RequestValidator = validator.New(validator.WithRequiredStructEnabled())

func PositiveStringValidationRule(value string) bool {
	bigIntValue := new(big.Int)
	bigIntValue.SetString(value, 10)
	return bigIntValue.Cmp(big.NewInt(0)) > 0
}

func ZeroOrPositiveStringValidationRule(value string) bool {
	bigIntValue := new(big.Int)
	bigIntValue.SetString(value, 10)
	return bigIntValue.Cmp(big.NewInt(0)) >= 0
}

func NonNegativeStringValidationRule(value string) bool {
	bigIntValue := new(big.Int)
	bigIntValue.SetString(value, 10)
	return bigIntValue.Cmp(big.NewInt(0)) >= 0
}

func percentageFeeValidator(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	if val < 0 || val > 100 {
		return false
	}
	valTimes100 := val * 100
	return math.Mod(valTimes100, 1) == 0
}

func init() {
	if err := registerValidations(); err != nil {
		log.Fatal("Error registering validations: ", err)
	}
}

func registerValidations() error {
	if err := RequestValidator.RegisterValidation("positive_string", func(field validator.FieldLevel) bool {
		return PositiveStringValidationRule(field.Field().String())
	}); err != nil {
		return fmt.Errorf("registering positive_string validation: %w", err)
	}

	if err := RequestValidator.RegisterValidation("percentage_fee", percentageFeeValidator); err != nil {
		return fmt.Errorf("registering percentage_fee validation: %w", err)
	}

	if err := RequestValidator.RegisterValidation("non_negative_string", func(field validator.FieldLevel) bool {
		return NonNegativeStringValidationRule(field.Field().String())
	}); err != nil {
		return fmt.Errorf("registering non_negative_string validation: %w", err)
	}

	if err := RequestValidator.RegisterValidation("zero_or_positive_string", func(field validator.FieldLevel) bool {
		return ZeroOrPositiveStringValidationRule(field.Field().String())
	}); err != nil {
		return fmt.Errorf("registering zero_or_positive_string validation: %w", err)
	}

	return nil
}

type ErrorDetails = map[string]any

type ErrorResponse struct {
	Message     string       `json:"message"`
	Details     ErrorDetails `json:"details"`
	Timestamp   int64        `json:"timestamp"`
	Recoverable bool         `json:"recoverable"`
}

func NewErrorResponseWithDetails(message string, details ErrorDetails, recoverable bool) *ErrorResponse {
	return &ErrorResponse{Message: message, Details: details, Timestamp: time.Now().Unix(), Recoverable: recoverable}
}

func NewErrorResponse(message string, recoverable bool) *ErrorResponse {
	return NewErrorResponseWithDetails(message, make(ErrorDetails), recoverable)
}

func DetailsFromError(err error) ErrorDetails {
	details := make(ErrorDetails)
	details["error"] = err.Error()
	return details
}

func JsonResponse(w http.ResponseWriter, statusCode int) {
	JsonResponseWithBody[any](w, statusCode, nil)
}

func JsonErrorResponse(w http.ResponseWriter, code int, response *ErrorResponse) {
	JsonResponseWithBody(w, code, response)
}

func JsonResponseWithBody[T any](w http.ResponseWriter, statusCode int, body *T) {
	var err error
	w.Header().Set(HeaderContentType, ContentTypeJson)
	w.WriteHeader(statusCode)
	if body == nil {
		return
	} else if err = json.NewEncoder(w).Encode(body); err != nil {
		responseError := NewErrorResponse("Unable to build response", true)
		if err = json.NewEncoder(w).Encode(responseError); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func DecodeRequestError(w http.ResponseWriter, err error) {
	log.Error("Error decoding request: ", err.Error())
	jsonErr := NewErrorResponse(fmt.Sprintf("Error decoding request: %v", err), true)
	JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
}

func ValidateRequestError(w http.ResponseWriter, err error) {
	log.Error("Error validating request: ", err.Error())
	jsonErr := NewErrorResponse(fmt.Sprintf("Error validating request: %v", err), true)
	JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
}

func DecodeRequest[T any](w http.ResponseWriter, req *http.Request, body *T) error {
	var err error
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err = dec.Decode(body); err != nil {
		DecodeRequestError(w, err)
		return err
	}
	return nil
}

func ValidateRequest[T any](w http.ResponseWriter, body *T) error {
	var validationErrors validator.ValidationErrors
	err := RequestValidator.Struct(body)
	if err == nil {
		return nil
	} else if !errors.As(err, &validationErrors) {
		ValidateRequestError(w, err)
		return err
	}
	details := make(ErrorDetails)
	for _, field := range validationErrors {
		details[field.Field()] = "validation failed: " + field.Tag()
	}
	jsonErr := NewErrorResponseWithDetails("validation error", details, true)
	JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
	return err
}

func RequiredQueryParam(name string) error {
	return fmt.Errorf("required query parameter %s is missing", name)
}
