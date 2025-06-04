package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
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
	_, ok := bigIntValue.SetString(value, 10)
	if !ok {
		return false
	}
	return bigIntValue.Cmp(big.NewInt(0)) > 0
}

func decimalPlacesValidator(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	param := fl.Param()
	maxDecimals, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	factor := math.Pow10(maxDecimals)
	valTimesFactor := val * factor
	diff := math.Abs(valTimesFactor - math.Round(valTimesFactor))
	return diff < 1e-9
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

	if err := RequestValidator.RegisterValidation("max_decimal_places", decimalPlacesValidator); err != nil {
		return fmt.Errorf("registering max_decimal_places validation: %w", err)
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

func getValidationMessage(field validator.FieldError) string {
	if field.Field() == "FeePercentage" {
		switch field.Tag() + ":" + field.Param() {
		case "gte:0":
			return "Fee percentage cannot be negative. Please enter a value between 0% and 100%."
		case "lte:100":
			return "Fee percentage cannot exceed 100%. Please enter a value between 0% and 100%."
		}
	}
	switch field.Tag() {
	case "required":
		return "is required"
	case "numeric":
		return "must be numeric"
	case "positive_string":
		return "must be a positive number"
	case "gte", "lte":
		op := map[string]string{
			"gte": "greater than or equal to ",
			"lte": "less than or equal to ",
		}
		return "must be " + op[field.Tag()] + field.Param()
	case "max_decimal_places":
		return fmt.Sprintf("must have at most %s decimal places", field.Param())
	default:
		return "validation failed: " + field.Tag()
	}
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
		details[field.Field()] = getValidationMessage(field)
	}
	jsonErr := NewErrorResponseWithDetails("validation error", details, true)
	JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
	return err
}

func RequiredQueryParam(name string) error {
	return fmt.Errorf("required query parameter %s is missing", name)
}
