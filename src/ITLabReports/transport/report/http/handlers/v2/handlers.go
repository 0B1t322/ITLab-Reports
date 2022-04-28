package handlers

import (
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v2"
	. "github.com/RTUITLab/ITLab-Reports/transport/report/http/endpoints/v2"

	genhttp "github.com/RTUITLab/ITLab-Reports/pkg/transport/http"

	"github.com/RTUITLab/ITLab-Reports/transport/http/options/serverbefore"
	errenc "github.com/RTUITLab/ITLab-Reports/transport/report/http/options/servererrorencoder"
	kithttp "github.com/go-kit/kit/transport/http"
)

// GetReports
//
// @Tags reports
//
// @Summary return reports according to filters
// 
// @Description.markdown get_reports_v2
//
// @Router /reports/v2/reports [get]
//
// @Security ApiKeyAuth
//
// @Param offset query number false "offset"
// 
// @Param limit query number false "limit"
// 
// @Param dateBegin query string false "date in RFC3339"
// 
// @Param dateEnd query string false "date in RFC3339"
// 
// @Param match query string false "match query"
// 
// @Param sortBy query string false "sorting query"
//
// @Produce json
//
// @Success 200 {object} dto.GetReportsResp
func GetReports(
	e Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.GetReports,
		dto.DecodeGetReportsReq,
		dto.EncodeGetReportsResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
		kithttp.ServerErrorEncoder(
			errenc.EncodeError,
		),
	)
}