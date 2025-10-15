package handlers

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type GetPeginReportUseCase interface {
	Run(ctx context.Context, startDate, endDate time.Time) (reports.GetPeginReportResult, error)
}

// NewGetReportsPeginHandler
// @Title Get Pegin Reports
// @Description Get the last pegins on the API. Included in the management API.
// @Param startDate query string true "Start date for the report. Supports YYYY-MM-DD (expands to full day) or ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Param endDate query string true "End date for the report. Supports YYYY-MM-DD (expands to end of day) or ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Success 200 pkg.GetPeginReportResponse
// @Route /reports/pegin [get]
func NewGetReportsPeginHandler(
	singleFlightGroup *singleflight.Group,
	singleFlightKey string,
	useCase GetPeginReportUseCase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var requestParams pkg.GetReportsByPeriodRequest
		var err error
		requestParams.StartDate = req.URL.Query().Get("startDate")
		requestParams.EndDate = req.URL.Query().Get("endDate")

		if err = rest.ValidateRequest(w, &requestParams); err != nil {
			return
		}

		if err = requestParams.ValidateGetReportsByPeriodRequest(); err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Validation error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		startTime, endTime, err := requestParams.GetTimestamps()
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails("Date conversion error", rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusBadRequest, jsonErr)
			return
		}

		rawPeginReport, err, shared := singleFlightGroup.Do(
			CalculateSingleFlightKey(singleFlightKey, req),
			// callback function signature comes from the std lib we can't modify it
			// nolint:contextcheck
			func() (any, error) {
				return useCase.Run(req.Context(), startTime, endTime)
			})
		peginReport, ok := rawPeginReport.(reports.GetPeginReportResult)
		if !ok {
			jsonErr := rest.NewErrorResponse("Internal error parsing result", false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		} else if shared {
			log.Info("GetPeginReport result was shared with multiple requests")
		}

		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		response := pkg.GetPeginReportResponse{
			NumberOfQuotes:     peginReport.NumberOfQuotes,
			MinimumQuoteValue:  peginReport.MinimumQuoteValue.AsBigInt(),
			MaximumQuoteValue:  peginReport.MaximumQuoteValue.AsBigInt(),
			AverageQuoteValue:  peginReport.AverageQuoteValue.AsBigInt(),
			TotalFeesCollected: peginReport.TotalFeesCollected.AsBigInt(),
			AverageFeePerQuote: peginReport.AverageFeePerQuote.AsBigInt(),
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
