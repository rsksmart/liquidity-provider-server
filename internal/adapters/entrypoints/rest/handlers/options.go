package handlers

import "net/http"

func NewOptionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
