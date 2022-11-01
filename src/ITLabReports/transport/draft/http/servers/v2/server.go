package servers

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/endpoints/v2"
	. "github.com/RTUITLab/ITLab-Reports/transport/draft/http/handlers/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/gorilla/mux"
)

type serverOptions struct {
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

	e = BuildMiddlewares(e, s)

	r.Handle(
		"/reports/v2/draft",
		GetDrafts(e),
	).Methods(http.MethodGet)
}

func BuildMiddlewares(
	e endpoints.Endpoints,
	opt *serverOptions,
) endpoints.Endpoints {
	e.GetDrafts.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetDraftsReq, *dto.GetDraftsResp](opt.auther),
		middlewares.SetUserID[*dto.GetDraftsReq, *dto.GetDraftsResp](),
	)

	return e
}
