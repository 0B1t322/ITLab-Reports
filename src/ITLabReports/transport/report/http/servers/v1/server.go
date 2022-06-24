package servers

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/endpoints/v1"
	serr "github.com/RTUITLab/ITLab-Reports/transport/report/http/errors"
	. "github.com/RTUITLab/ITLab-Reports/transport/report/http/handlers/v1"
	"github.com/gorilla/mux"
)

var (
	EmployeeCantBeEmpty = errors.New("Employee can't be empty")
)

type DraftService = endpoints.DraftService

type serverOptions struct {
	auther       middlewares.Auther
	draftService DraftService
	idCheker     middlewares.IdChecker
}

type ServerOptions func(s *serverOptions)

func WithAuther(a middlewares.Auther) ServerOptions {
	return func(s *serverOptions) {
		s.auther = a
	}
}

func WithDraftService(service DraftService) ServerOptions {
	return func(s *serverOptions) {
		s.draftService = service
	}
}

func WithIdChecker(checker middlewares.IdChecker) ServerOptions {
	return func(s *serverOptions) {
		s.idCheker = checker
	}
}

func MergeServerOptions(opts ...ServerOptions) *serverOptions {
	s := &serverOptions{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// NewServer copy given endpoints and build middlewares for http
func NewServer(
	ctx context.Context,
	r *mux.Router,
	ends report.Endpoints,
	opts ...ServerOptions,
) endpoints.Endpoints {
	s := &serverOptions{}

	for _, opt := range opts {
		opt(s)
	}

	e := endpoints.NewEndpoints(ends)

	if s.draftService != nil {
		r.Handle(
			"/reports/v1/report_from_draft/{id}",
			CreateReportFromDraftHandler(
				endpoints.NewDraftServiceEndpoints(s.draftService, e),
			),
		)
	}

	e = BuildMiddlewares(e, s)

	r.Handle(
		"/reports",
		GetReports(e),
	).Methods(http.MethodGet)

	r.Handle(
		"/reports/employee/{employee}",
		GetReportsForEmployee(e),
	).Methods(http.MethodGet)

	r.Handle(
		"/reports/{id}",
		GetReport(e),
	).Methods(http.MethodGet)

	r.Handle(
		"/reports",
		CreateReport(e),
	).Methods(http.MethodPost)

	return e
}

func BuildMiddlewares(
	e endpoints.Endpoints,
	opt *serverOptions,
) endpoints.Endpoints {
	e.GetReport.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
		middlewares.RunMiddlewareIfAllFail(
			middlewares.CheckUserIsReporterOrImplementer[*dto.GetReportReq, *dto.GetReportResp](),
			middlewares.IsSuperAdmin[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
			middlewares.IsAdmin[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
		),
	)

	e.CreateReport.AddCustomMiddlewares(
		middlewares.Auth[*dto.CreateReportReq, *dto.CreateReportResp](opt.auther),
		middlewares.SetReporter[*dto.CreateReportReq, *dto.CreateReportResp](),
		middlewares.SetImplementerIfEmpty[*dto.CreateReportReq, *dto.CreateReportResp](),
		middlewares.CheckIds[*dto.CreateReportReq, *dto.CreateReportResp](opt.idCheker),
	)

	e.GetReportsForEmployee.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](opt.auther),
		middlewares.EmployeeIsNotEmpty[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](
			func() error {
				return errors.Wrap(EmployeeCantBeEmpty, serr.ValidationError)
			},
		),
		middlewares.MergeMiddlewaresIntoOr(
			middlewares.IsSuperAdmin[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](opt.auther),
			middlewares.IsAdmin[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](opt.auther),
			middlewares.UserIsEmployee[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](),
		),
	)

	e.GetReports.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
		middlewares.RunMiddlewareIfAllFail(
			middlewares.SetReporterAndImplementer[*dto.GetReportsReq, *dto.GetReportsResp](),
			middlewares.IsAdmin[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
			middlewares.IsSuperAdmin[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
		),
	)

	return e
}
