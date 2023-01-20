package response

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Data = interface{}
type Meta = interface{}
type Details = interface{}

type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
	Meta    Meta `json:"meta"`
}

func New(s bool, d Data, m Meta) Response {
	return Response{s, d, m}
}

func (r *Response) JsonMarshal() []byte {
	j, err := json.Marshal(r)

	if err != nil {
		return []byte(err.Error())
	}

	return j
}

type ErrorBody struct {
	//Code        string  `json:"code"`
	Message     string  `json:"message"`
	Details     Details `json:"details"`
	Timestamp   string  `json:"timestamp"`
	Recoverable bool    `json:"recoverable"`
}

func newError(m string, d Details, t string, r bool) ErrorBody {
	return ErrorBody{
		Message:     m,
		Details:     d,
		Timestamp:   t,
		Recoverable: r,
	}
}

func httpError(w http.ResponseWriter, er ErrorBody, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)

	err := enc.Encode(er)
	if err != nil {
		log.Fatal("[response package] error encoding response: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	fmt.Fprintln(w, er)
}
