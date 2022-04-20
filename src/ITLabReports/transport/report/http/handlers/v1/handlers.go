package handlers

import (
	// "context"
	"net/http"

	genhttp "github.com/RTUITLab/ITLab-Reports/pkg/transport/http"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
	. "github.com/RTUITLab/ITLab-Reports/transport/report/http/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/options/serverbefore"
	kithttp "github.com/go-kit/kit/transport/http"
)

// GetReport
//
// @Tags reports
//
// @Summary return report
//
// @Router /reports/{id} [get]
//
// @Security ApiKeyAuth
//
// @Param id path string true "id of report"
//
// @Produce json
//
// @Success 200 {object} dto.GetReportResp
func GetReport(
	e Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.GetReport,
		dto.DecodeGetReportReq,
		dto.EncodeGetReportResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}

// GetEmployeeReports
//
// @Tags reports
//
// @Summary get reports for employee
//
// @Description get reports for current employee
//
// @Router /reports/employee/{employee} [get]
//
// @Security ApiKeyAuth
//
// @Param dateBegin query string false "begin date of reports"
// 
// @Param dateEnd query string false "end date of reports"
// 
// @Param employee path string true "employee user id"
//
// @Produce json
//
// @Success 200 {array} dto.GetReportResp
func GetReportsForEmployee(
	e Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.GetReportsForEmployee,
		dto.DecodeGetReportsForEmployeeReq,
		dto.EncodeGetReportsResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}

// GetReports
//
// @Tags reports
//
// @Summary get report list
//
// @Router /reports [get]
//
// @Security ApiKeyAuth
//
// @Param sorted_by query string false "sort by this field" Enums(name, date)
//
// @Produce json
//
// @Success 200 {array} dto.GetReportResp
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
	)
}

// CreateReport
//
// @Tags reports
//
// @Summary create report
//
// @Description create report
//
// @Description query value implementor indicate who make things described in report
// 
// @Description name is optional field, if use it text should be in standart md format
// 
// @Description if name is not given text should be in format "report_name@\n\t\n@report_rext"
//
// @Description if implementor is not specified report-maker is implementor
//
// @Router /reports [post]
//
// @Security ApiKeyAuth
//
// @Param implementor query string false "implemntor user id"
//
// @Param report body dto.CreateReportReq true "body"
//
// @Acceprt json
//
// @Produce json
//
// @Success 200 {object} dto.CreateReportResp
func CreateReports(
	e Endpoints,
) http.Handler {
	return genhttp.NewServer(
		e.CreateReport,
		dto.DecodeCreateReportReq,
		dto.EncodeCreateReportResp,
		kithttp.ServerBefore(
			serverbefore.TokenFromReq,
		),
	)
}