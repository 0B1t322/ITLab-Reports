package middlewares

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

type EndpointWithContext[Req any, Resp any] func(ctx mcontext.MiddlewareContext, req Req) (Resp, error)

func (e EndpointWithContext[Req, Resp]) ToEndpoint() endpoint.Endpoint[Req, Resp] {
	return func(ctx context.Context, req Req) (Resp, error) {
		mctx := mcontext.CreateOrGetFrom(ctx)
		return e(mctx, req)
	}
}

func EndpointWithCtxFromEndpoint[Req any, Resp any](e endpoint.Endpoint[Req, Resp]) EndpointWithContext[Req, Resp] {
	return func(ctx mcontext.MiddlewareContext, req Req) (Resp, error) {
		return e(ctx, req)
	}
}

type MiddlewareWithContext[Req any, Resp any] func(EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp]

func (m MiddlewareWithContext[Req, Resp]) ToMiddleware() endpoint.Middleware[Req, Resp] {
	return func(e endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return m(EndpointWithCtxFromEndpoint(e)).ToEndpoint()
	}
}

func MiddlewareWithCTXFrom[Req any, Resp any](m endpoint.Middleware[Req, Resp]) MiddlewareWithContext[Req, Resp] {
	return func(e EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return EndpointWithCtxFromEndpoint(m(e.ToEndpoint()))
	}
}