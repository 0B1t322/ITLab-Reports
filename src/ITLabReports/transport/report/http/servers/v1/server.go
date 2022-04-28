package servers

import (
	"context"
	"net/http"

	agragate "github.com/RTUITLab/ITLab-Reports/aggragate/report"
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

type DraftService interface {
	IsDraftNotFoundErr(error) bool
	IsDraftIdNotValidErr(error) bool
	// Should throws errors
	// 	Draft not found
	// 	Draft id is invalid
	GetDraft(ctx context.Context, id string) (*agragate.Report, error)

	// Should throws errors
	// 	Draft not found
	// 	Draft id is invalid
	DeleteDraft(ctx context.Context, id string) error
}

type serverOptions struct {
	auther       middlewares.Auther
	draftService DraftService
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
) (endpoints.Endpoints) {
	e := endpoints.NewEndpoints(ends)

	s := &serverOptions{}

	for _, opt := range opts {
		opt(s)
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
		CreateReports(e),
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
			middlewares.CheckUserIsReporter[*dto.GetReportReq, *dto.GetReportResp](),
			middlewares.IsSuperAdmin[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
			middlewares.IsAdmin[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
		),
	)

	e.CreateReport.AddCustomMiddlewares(
		middlewares.Auth[*dto.CreateReportReq, *dto.CreateReportResp](opt.auther),
		middlewares.SetReporter[*dto.CreateReportReq, *dto.CreateReportResp](),
		middlewares.SetImplementerIfEmpty[*dto.CreateReportReq, *dto.CreateReportResp](),
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
