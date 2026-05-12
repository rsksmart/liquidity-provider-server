package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

const (
	HeaderContentType = "Content-Type"

	ContentTypeJson = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"

	StartDateParam = "startDate"
	EndDateParam   = "endDate"
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

func ConfirmationsMapValidator(fl validator.FieldLevel) bool {
	kind := fl.Field().Kind()
	if kind != reflect.Map {
		return false
	}
	confirmations, ok := fl.Field().Interface().(map[string]uint16)
	if !ok {
		return false
	}
	if len(confirmations) == 0 {
		return false
	}
	for key := range confirmations {
		bigInt, valid := new(big.Int).SetString(key, 10)
		if !valid || bigInt.Sign() <= 0 {
			return false
		}
	}
	return true
}

func init() {
	if err := registerValidations(); err != nil {
		log.Fatal("Error registering validations: ", err)
	}
}

func positiveIntegerBigintValidator(field validator.FieldLevel) bool {
	fieldValue := field.Field().Interface()

	// Handle both *big.Int and big.Int
	var bigIntVal *big.Int
	switch v := fieldValue.(type) {
	case *big.Int:
		bigIntVal = v
	case big.Int:
		bigIntVal = &v
	default:
		return false // Not a big.Int type
	}

	if bigIntVal == nil {
		return false
	}

	return bigIntVal.Sign() > 0 // Only positive values (> 0)
}

func notBlankValidator(field validator.FieldLevel) bool {
	str, ok := field.Field().Interface().(string)
	if !ok {
		return false // Not a string type
	}

	// Check if string is not empty after trimming whitespace
	return len(strings.TrimSpace(str)) > 0
}

func positiveStringValidator(field validator.FieldLevel) bool {
	return PositiveStringValidationRule(field.Field().String())
}

func registerValidator(tag string, fn validator.Func) error {
	if err := RequestValidator.RegisterValidation(tag, fn); err != nil {
		return fmt.Errorf("registering %s validation: %w", tag, err)
	}
	return nil
}

func registerValidations() error {
	validators := map[string]validator.Func{
		"positive_string":         positiveStringValidator,
		"max_decimal_places":      decimalPlacesValidator,
		"confirmations_map":       ConfirmationsMapValidator,
		"positive_integer_bigint": positiveIntegerBigintValidator,
		"not_blank":               notBlankValidator,
	}

	for tag, fn := range validators {
		if err := registerValidator(tag, fn); err != nil {
			return err
		}
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
	jsonErr := NewErrorResponseWithDetails("Error decoding request", DetailsFromError(err), true)
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

// nolint: cyclop
func getValidationMessage(field validator.FieldError) string {
	switch field.Tag() {
	case "required":
		return "is required"
	case "numeric":
		return "must be numeric"
	case "positive_string":
		return "must be a positive number"
	case "gte":
		return "must be greater than or equal to " + field.Param()
	case "lte":
		return "must be less than or equal to " + field.Param()
	case "max_decimal_places":
		return fmt.Sprintf("must have at most %s decimal places", field.Param())
	case "positive_integer_bigint":
		return "must be a positive integer"
	case "not_blank":
		return "cannot be blank"
	case "confirmations_map":
		return "must not be empty"
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

func ParseDateRange(req *http.Request, dateFormat string) (time.Time, time.Time, error) {
	start := req.URL.Query().Get(StartDateParam)
	end := req.URL.Query().Get(EndDateParam)
	if start == "" || end == "" {
		missing := []string{}
		if start == "" {
			missing = append(missing, StartDateParam)
		}
		if end == "" {
			missing = append(missing, EndDateParam)
		}
		return time.Time{}, time.Time{}, fmt.Errorf("missing required parameters: %v", missing)
	}
	startDate, err := time.Parse(dateFormat, start)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date format: %w", err)
	}
	endDate, err := time.Parse(dateFormat, end)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date format: %w", err)
	}
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, time.UTC)
	return startDate, endDate, nil
}

func ValidateDateRange(startDate, endDate time.Time, dateFormat string) error {
	if endDate.Before(startDate) {
		return fmt.Errorf("invalid date range: end date %s is before start date %s",
			endDate.Format(dateFormat), startDate.Format(dateFormat))
	}
	return nil
}
