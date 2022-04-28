package servers

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/handlers/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/gorilla/mux"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/endpoints/v1"
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
) endpoints.Endpoints {
	e := endpoints.NewEndpoints(ends)

	s := &serverOptions{}

	for _, opt := range opts {
		opt(s)
	}

	e = BuildMiddlewares(e, s)

	r.Handle(
		"/reports/v1/draft",
		handlers.GetDrafts(
			e,
		),
	).Methods(http.MethodGet)

	r.Handle(
		"/reports/v1/draft/{id}",
		handlers.GetDraft(
			e,
		),
	).Methods(http.MethodGet)

	r.Handle(
		"/reports/v1/draft/{id}",
		handlers.DeleteDraft(
			e,
		),
	).Methods(http.MethodDelete)

	r.Handle(
		"/reports/v1/draft/{id}",
		handlers.UpdateDraft(
			e,
		),
	).Methods(http.MethodPut)

	return e
}

func BuildMiddlewares(
	ends endpoints.Endpoints,
	opt *serverOptions,
) endpoints.Endpoints {
	e := ends

	e.CreateDraft.AddCustomMiddlewares(
		middlewares.Auth[*dto.CreateDraftReq, *dto.CreateDraftResp](opt.auther),
		middlewares.SetReporter[*dto.CreateDraftReq, *dto.CreateDraftResp](),
	)

	e.DeleteDraft.AddCustomMiddlewares(
		middlewares.Auth[*dto.DeleteDraftReq, *dto.DeleteDraftResp](opt.auther),
		middlewares.UserIsOwner[*dto.DeleteDraftReq, *dto.DeleteDraftResp](e),
	)

	e.GetDraft.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetDraftReq, *dto.GetDraftResp](opt.auther),
		middlewares.UserIsOwner[*dto.GetDraftReq, *dto.GetDraftResp](e),
	)

	e.UpdateDraft.AddCustomMiddlewares(
		middlewares.Auth[*dto.UpdateDraftReq, *dto.UpdateDraftResp](opt.auther),
		middlewares.UserIsOwner[*dto.UpdateDraftReq, *dto.UpdateDraftResp](e),
	)

	e.GetDrafts.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetDraftsReq, *dto.GetDraftsResp](opt.auther),
		middlewares.SetUserID[*dto.GetDraftsReq, *dto.GetDraftsResp](),
	)

	return e
}
