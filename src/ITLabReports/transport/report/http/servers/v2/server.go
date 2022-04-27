package servers

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/endpoints/v2"
	. "github.com/RTUITLab/ITLab-Reports/transport/report/http/handlers/v2"
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

func MergeServerOptions(opts ...ServerOptions) *serverOptions {
	s := &serverOptions{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

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

	r.Handle(
		"/reports/v2/reports",
		GetReports(e),
	).Methods(http.MethodGet)
}

func BuildMiddlewares(
	e endpoints.Endpoints,
	opt *serverOptions,
) endpoints.Endpoints {
	e.GetReports.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
		middlewares.RunMiddlewareIfAllFail(
			middlewares.SetReporterAndImplementer[*dto.GetReportsReq, *dto.GetReportsResp](),
			// If fail
			middlewares.IsAdmin[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
			middlewares.IsSuperAdmin[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
		),
	)

	return e
}