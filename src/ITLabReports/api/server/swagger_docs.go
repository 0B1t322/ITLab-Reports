package server

import "ITLabReports/models"

func init() {
	_ = models.Report{}
}

type swaggerDocs struct{}

// GetAllReport
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
// @Success 200 {array} models.Report
func (swaggerDocs) getAllReprots() {}

type CreateReportRequest struct {
	Text string `json:"text"`
}

type CreateReportResponce struct {
	ID           string           `json:"id"`
	Assignees    models.Assignees `json:"assignees"`
	Date         string           `json:"date"`
	Text         string           `json:"text"`
	Archived     bool             `json:"archived"`
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
// @Description if implementor is not specified report-maker is implementor
//
// @Router /reports [post]
//
// @Security ApiKeyAuth
//
// @Param implementor query string false "implemntor user id"
//
// @Param report body server.CreateReportRequest true "body"
//
// @Acceprt json
//
// @Produce json
//
// @Success 200 {object} server.CreateReportResponce
func (swaggerDocs) createReport() {}

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
// @Success 200 {array} models.Report
func(swaggerDocs) getReportsForEmployee() {}

// GetArchivedReports
//
// @Tags reports
//
// @Summary get archived reports
//
// @Router /reports/archived [get]
//
// @Security ApiKeyAuth
//
// @Produce json
//
// @Success 200 {array} models.Report
func(swaggerDocs) getArchivedReports() {}

// GetReport
//
// @Tags reports
//
// @Summary get report
//
// @Router /reports/{id} [get]
//
// @Security ApiKeyAuth
// 
// @Param id path string true "id of report"
//
// @Produce json
//
// @Success 200 {object} models.Report
func(swaggerDocs) getReport() {}
