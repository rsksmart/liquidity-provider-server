package handlers

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"net/http"
)

const (
	PegInReportSingleFlightKey  = "pegin-report-singleflight"
	PegOutReportSingleFlightKey = "pegout-report-singleflight"
)

var SingleFlightGroup = new(singleflight.Group)

func CalculateSingleFlightKey(baseKey string, req *http.Request) string {
	return fmt.Sprintf("%s_%s", baseKey, req.URL.String())
}
