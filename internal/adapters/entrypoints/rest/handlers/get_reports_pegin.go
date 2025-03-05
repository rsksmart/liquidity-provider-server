package handlers

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"net/http"
)

// NewGetReportsPeginHandler
// @Title Get Pegin Reports
// @Description Get the last pegins on the API. Included in the management API.
// @Param PeginReportsRequest  body pkg.PeginReportsRequest true "Date range for the report, case not provided will be last 30 days"
// @Success 200 object
// @Route /reports/pegin [get]
func NewGetReportsPeginHandler(useCase *pegin.GetPeginReportUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		var response pegin.GetPeginReportResult

		response, err = useCase.Run(req.Context())
		if err != nil {
			jsonErr := rest.NewErrorResponseWithDetails(UnknownErrorMessage, rest.DetailsFromError(err), false)
			rest.JsonErrorResponse(w, http.StatusInternalServerError, jsonErr)
			return
		}
		rest.JsonResponseWithBody(w, http.StatusOK, &response)
	}
}
