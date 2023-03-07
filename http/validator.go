package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
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
			BuildErrorResponse(w, http.StatusUnprocessableEntity, message)
			return false
		}
		return true
	}
}
func BuildErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
