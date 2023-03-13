package http

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Validate(schema interface{}) func(w http.ResponseWriter) bool {
	return func(w http.ResponseWriter) bool {
		err := validator.New().Struct(schema)
		if err != nil {
			var message string
			for i, err := range err.(validator.ValidationErrors) {
				if i > 0 {
					message += ", "
				}
				message += fmt.Sprintf("%s is %s", err.Field(), err.Tag())
			}
			customError := NewServerError(message, make(Details), true)
			ResponseError(w, customError, http.StatusBadRequest)
			return false
		}
		return true
	}
}
