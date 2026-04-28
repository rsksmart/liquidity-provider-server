package integration

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
	"math/big"
	"net/http"
)

type Execution struct {
	Body   any
	Method string
	URL    string
}

type Result[responseType any] struct {
	Response    responseType
	RawResponse []byte
	StatusCode  int
}

func ExecuteHttpRequest[responseType any](execution Execution) (Result[responseType], error) {
	payload, err := json.Marshal(execution.Body)
	if err != nil {
		return Result[responseType]{}, err
	}
	req, err := http.NewRequestWithContext(context.Background(), execution.Method, execution.URL, bytes.NewBuffer(payload))
	if err != nil {
		return Result[responseType]{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return Result[responseType]{}, err
	}
	defer func(res *http.Response) {
		closingErr := res.Body.Close()
		if closingErr != nil {
			log.Debug("Error closing response body: ", closingErr)
		}
	}(res)

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return Result[responseType]{}, err
	}

	var response responseType
	if len(bodyBytes) > 0 {
		if err = json.Unmarshal(bodyBytes, &response); err != nil {
			return Result[responseType]{}, err
		}
	}

	result := Result[responseType]{
		Response:    response,
		StatusCode:  res.StatusCode,
		RawResponse: bodyBytes,
	}
	return result, nil
}

func AssertFields(s *suite.Suite, expectedFields []string, object map[string]any) {
	for _, field := range expectedFields {
		_, exists := object[field]
		s.Require().True(exists, "Field %v is missing", field)
	}
}

func SumAll(numbers ...*big.Int) *big.Int {
	sum := new(big.Int)
	for _, number := range numbers {
		sum.Add(sum, number)
	}
	return sum
}
