package servers

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/endpoints/v1"
	. "github.com/RTUITLab/ITLab-Reports/transport/report/http/handlers/v1"
	"github.com/gorilla/mux"
)


type serverOptions struct{
	auther middlewares.Auther
}

type ServerOptions func(s *serverOptions)

func WithAuther(a middlewares.Auther) ServerOptions {
	return func(s *serverOptions) {
		s.auther = a
	}
}

// NewServer copy given endpoints and build middlewares for http
func NewServer(
	ctx context.Context,
	r *mux.Router,
	ends report.Endpoints,
	opts ...ServerOptions,
) {
	e := endpoints.NewEndpoints(ends)
	
	s := &serverOptions{}

	for _, opt := range opts {
		opt(s)
	}

	e = buildMiddlewares(e, s)

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
}

func buildMiddlewares(
	e endpoints.Endpoints,
	opt *serverOptions,
) endpoints.Endpoints {
	e.GetReport.AddCustomMiddlewares(
		middlewares.MiddlewareWithCTXFrom(
			endpoint.MiddlewareFromGoKitMiddleware[*dto.GetReportReq, *dto.GetReportResp](
				opt.auther.Auth().ToMiddleware().ToGoKitMiddleware(),
			),
		),
		middlewares.CheckUserIsReporter[*dto.GetReportReq, *dto.GetReportResp](),
	)

	e.CreateReport.AddCustomMiddlewares(
		middlewares.MiddlewareWithCTXFrom(
			endpoint.MiddlewareFromGoKitMiddleware[*dto.CreateReportReq, *dto.CreateReportResp](
				opt.auther.Auth().ToMiddleware().ToGoKitMiddleware(),
			),
		),
		middlewares.SetReporter[*dto.CreateReportReq, *dto.CreateReportResp](),
	)

	e.GetReportsForEmployee.AddCustomMiddlewares(
		middlewares.MiddlewareWithCTXFrom(
			endpoint.MiddlewareFromGoKitMiddleware[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](
				opt.auther.Auth().ToMiddleware().ToGoKitMiddleware(),
			),
		),
		middlewares.MergeMiddlewaresIntoOr(
			middlewares.MiddlewareWithCTXFrom(
				endpoint.MiddlewareFromGoKitMiddleware[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](
					opt.auther.IsAdmin().ToMiddleware().ToGoKitMiddleware(),
				),
			),
			middlewares.UserIsEmployee[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp](),
		),
	)

	e.GetReports.AddCustomMiddlewares(
		middlewares.MiddlewareWithCTXFrom(
			endpoint.MiddlewareFromGoKitMiddleware[*dto.GetReportsReq, *dto.GetReportsResp](
				opt.auther.Auth().ToMiddleware().ToGoKitMiddleware(),
			),
		),
		middlewares.SetReporterAndImplementerIfFailed[*dto.GetReportsReq, *dto.GetReportsResp](
			opt.auther.IsAdmin().ToMiddleware().ToGoKitMiddleware(),
		),
	)

	return e
}