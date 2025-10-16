package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/singleflight"
)

const (
	PegInReportSingleFlightKey   = "pegin-report-singleflight"
	PegOutReportSingleFlightKey  = "pegout-report-singleflight"
	RevenueReportSingleFlightKey = "revenue-report-singleflight"
)

var SingleFlightGroup = new(singleflight.Group)

func CalculateSingleFlightKey(baseKey string, req *http.Request) string {
	return fmt.Sprintf("%s_%s", baseKey, req.URL.String())
}
