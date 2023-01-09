package response

import "encoding/json"

type Data = interface{}
type Meta = interface{}

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
