package reports

import (
	"fmt"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type ReportsController struct {
	reportsService reports.ReportsService
	auth           user.UserService

	ErrorFormatter
	AuthErrorHandler
	ReportMarshaller
}

func NewReportsController(
	reportsService reports.ReportsService,
	auth user.UserService,
) *ReportsController {
	formatter := ErrorFormatter{}
	c := &ReportsController{
		reportsService:   reportsService,
		auth:             auth,
		ErrorFormatter:   formatter,
		AuthErrorHandler: AuthErrorHandler{Formatter: formatter},
		ReportMarshaller: ReportMarshaller{},
	}

	return c
}

func NewReportsControllerFrom(i *do.Injector) (*ReportsController, error) {
	reportsService := do.MustInvoke[reports.ReportsService](i)
	auth := do.MustInvoke[user.UserService](i)

	return NewReportsController(reportsService, auth), nil
}

func (c *ReportsController) Build(r gin.IRouter) {
	r.GET("", c.GetReports)
	r.POST("", c.CreateReport)
	r.GET("/employee/:employee", c.GetReportsForEmployee)
	r.GET("/:id", c.GetReport)
	r.POST("/v1/report_from_draft/:id", c.CreateReportFromDraft)
}

func (rc *ReportsController) HandlerError(c *gin.Context, err error) {
	if rc.AuthErrorHandler.HandlerError(c, err) {
		return
	}

	switch {
	case err == reports.ErrReportNotFound:
		rc.FormatError(c, err, http.StatusNotFound)
		return
	case errors.Is(err, aggregate.ErrReportValidation):
		rc.FormatError(c, errors.Unwrap(err), http.StatusBadRequest)
		return
	case errors.Is(err, reports.ErrDraftNotFound):
		rc.FormatError(c, err, http.StatusNotFound)
		return
	default:
		logrus.WithFields(
			logrus.Fields{
				"controller": "reports",
			},
		).Error(err)
		rc.FormatError(c, err, http.StatusInternalServerError)
		return
	}
}

// GetReport
//
// @Tags reports
//
// @Summary return report
//
// @Description.markdown get_report
//
// @Router /reports/{id} [get]
//
// @Security ApiKeyAuth
//
// @Param id path string true "id of report"
//
// @Produce json
//
// @Success 200 {object} ReportView
func (rc *ReportsController) GetReport(c *gin.Context) {
	user, err := rc.auth.AuthUser(
		c.Request.Context(),
		c.GetHeader("Authorization"),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	var req GetReportReq
	{
		if err := c.ShouldBindUri(&req); err != nil {
			rc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	report, err := rc.reportsService.GetReport(
		c,
		rc.MarshallGetReportReq(req, user),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, rc.MarshallReportView(report))
}

// GetEmployeeReports
//
// @Tags reports
//
// @Summary get reports for employee
//
// @Description.markdown get_reports_for_employee
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
// @Success 200 {array} ReportView
func (rc *ReportsController) GetReportsForEmployee(c *gin.Context) {
	user, err := rc.auth.AuthUser(
		c.Request.Context(),
		c.GetHeader("Authorization"),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	var req GetEmployeeReportsReq
	{
		if err := c.ShouldBindUri(&req); err != nil {
			rc.FormatError(c, err, http.StatusBadRequest)
			return
		}

		if err := c.ShouldBindQuery(&req); err != nil {
			rc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	if !user.IsAdminOrSuperAdmin() && user.ID != req.EmployeeID {
		rc.FormatError(
			c,
			fmt.Errorf("You are not admin to get reports about this employee"),
			http.StatusForbidden,
		)
		return
	}

	reports, err := rc.reportsService.GetReports(
		c,
		rc.MarshallGetReportsForEmployeeReq(req, user),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, rc.MarshallReportsView(reports))
}

// GetReports
//
// @Tags reports
//
// @Summary get report list
//
// @Description.markdown get_reports_v1
//
// @Router /reports [get]
//
// @Security ApiKeyAuth
//
// @Param sorted_by query string false "sort by this field" Enums(name, date)
//
// @Produce json
//
// @Success 200 {array} ReportView
func (rc *ReportsController) GetReports(c *gin.Context) {
	user, err := rc.auth.AuthUser(
		c.Request.Context(),
		c.GetHeader("Authorization"),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	var req GetReportsReq
	{
		if err := c.ShouldBindQuery(&req); err != nil {
			rc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	reports, err := rc.reportsService.GetReports(
		c,
		rc.MarshallGetReportsReq(req, user),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, rc.MarshallReportsView(reports))
}

// CreateReport
//
// @Tags reports
//
// @Summary create report
//
// @Description create report
//
// @Description.markdown create_report_v1
//
// @Router /reports [post]
//
// @Security ApiKeyAuth
//
// @Param implementer query string false "implemnter user id"
//
// @Param report body CreateReportReq true "body"
//
// @Acceprt json
//
// @Produce json
//
// @Success 200 {object} ReportView
func (rc *ReportsController) CreateReport(c *gin.Context) {
	user, err := rc.auth.AuthUser(
		c.Request.Context(),
		c.GetHeader("Authorization"),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	var req CreateReportReq
	{
		if err := c.ShouldBindJSON(&req); err != nil {
			rc.FormatError(c, err, http.StatusBadRequest)
			return
		}

		if err := c.ShouldBindQuery(&req); err != nil {
			rc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	r, err := rc.MarshallCreateReportReq(req, user)
	if err != nil {
		rc.FormatError(c, err, http.StatusBadRequest)
		return
	}

	report, err := rc.reportsService.CreateReport(
		c,
		r,
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, rc.MarshallReportView(report))
}

// CreateReportFromDraft
//
// @Tags reports
//
// @Summary create report from draft
//
// @Description.markdown create_report_from_draft
//
// @Router /reports/v1/report_from_draft/{id} [post]
//
// @Param id path string true "id of draft"
//
// @Security ApiKeyAuth
//
// @Produce json
//
// @Success 201 {object} ReportView
func (rc *ReportsController) CreateReportFromDraft(c *gin.Context) {
	user, err := rc.auth.AuthUser(
		c.Request.Context(),
		c.GetHeader("Authorization"),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	var req CreateReportFromDraftReq
	{
		if err := c.ShouldBindUri(&req); err != nil {
			rc.FormatError(c, err, http.StatusBadRequest)
			return
		}
	}

	report, err := rc.reportsService.CreateReportFromDraft(
		c,
		rc.MarshallCreateReportFromDraftReq(req, user),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, rc.MarshallReportView(report))
}
