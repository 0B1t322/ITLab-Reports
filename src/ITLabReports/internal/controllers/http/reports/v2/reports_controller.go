package reports

import (
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v2/dto"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus"
)

type ReportsController struct {
	reportsService   reports.ReportsService
	auth             user.UserService
	marshaller       ReportsMarshaller
	approvedStrategy ApprovedStateStrategy

	ErrorFormatter
	AuthErrorHandler
}

func NewReportsController(
	reportsService reports.ReportsService,
	auth user.UserService,
	approvedStrategy ApprovedStateStrategy,
) *ReportsController {
	formatter := ErrorFormatter{}
	c := &ReportsController{
		reportsService:   reportsService,
		auth:             auth,
		ErrorFormatter:   formatter,
		AuthErrorHandler: AuthErrorHandler{Formatter: formatter},
		approvedStrategy: approvedStrategy,
	}

	return c
}

func NewReportsControllerFrom(i *do.Injector) (*ReportsController, error) {
	reportsService := do.MustInvoke[reports.ReportsService](i)
	auth := do.MustInvoke[user.UserService](i)

	var approvedStrategy ApprovedStateStrategy = NewInternalApprovedStrategy()
	{
		ovverideApprovedStrategy, err := do.Invoke[ApprovedStateStrategy](i)
		if err == nil {
			approvedStrategy = ovverideApprovedStrategy
		}
	}

	return NewReportsController(reportsService, auth, approvedStrategy), nil
}

func (c *ReportsController) Build(r gin.IRouter) {
	reports := r.Group("/v2/reports")
	{
		reports.GET("", c.GetReports)
	}
}

func (rc *ReportsController) HandlerError(c *gin.Context, err error) {
	if rc.AuthErrorHandler.HandlerError(c, err) {
		return
	}

	switch {
	default:
		logrus.WithFields(
			logrus.Fields{
				"controller": "reports",
				"version":    "v2",
				"transport":  "http",
				"handler":    c.HandlerName(),
			},
		).Error(err)
		rc.FormatError(c, err, http.StatusInternalServerError)
		return
	}
}

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
// @Param approvedState query string false "filtering on approved state" Enums(approved,notApproved)
//
// @Produce json
//
// @Success 200 {object} GetReportsResp
func (rc *ReportsController) GetReports(c *gin.Context) {
	user, err := rc.auth.AuthUser(
		c,
		c.GetHeader("Authorization"),
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	var req GetReportsReq
	{
		if err := c.ShouldBindQuery(&req); err != nil {
			rc.HandlerError(c, err)
			return
		}
	}

	r, err := rc.marshaller.MarshallGetReportsReq(req, user)
	if err != nil {
		rc.FormatError(c, err, http.StatusBadRequest)
		return
	}

	approvedStateQuery, err := rc.approvedStrategy.SetApproved(
		ApprovedStateReq{
			State: lo.Switch[dto.ApprovedState, aggregate.ReportState](req.ApprovedState).
				Case(dto.ApprovedState_Approved, aggregate.ReportStatePaid).
				Case(dto.ApprovedState_Not, aggregate.ReportStateCreated).
				Default("unknown"),
			Token: c.GetHeader("Authorization"),
			UserID: lo.Ternary(
				user.IsAdminOrSuperAdmin(),
				mo.None[string](),
				mo.Some(user.ID),
			),
		},
	)
	if err != nil {
		rc.FormatError(c, err, http.StatusConflict)
		return
	}

	r.Query.Filter.And = append(
		r.Query.Filter.And,
		approvedStateQuery,
	)

	reports, err := rc.reportsService.GetReports(
		c,
		r,
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	count, err := rc.reportsService.CountReports(
		c,
		r,
	)
	if err != nil {
		rc.HandlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, rc.marshaller.MarshallGetReportsResp(req, reports, count))
}
