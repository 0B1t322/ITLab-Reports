package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)


func Auth[Req any, Resp any](
	a Auther,
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			var req any = request
			_, err := a.Auth()(Nop[any, any])(ctx, req)
			if err != nil {
				return *new(Resp), err
			}

			return next(ctx, request)
		}
	}
}

func IsAdmin[Req any, Resp any](
	a Auther,
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			var req any = request
			_, err := a.IsAdmin()(Nop[any, any])(ctx, req)
			if err != nil {
				return *new(Resp), err
			}

			return next(ctx, request)
		}
	}
}

func IsSuperAdmin[Req any, Resp any](
	a Auther,
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			var req any = request
			_, err := a.IsSuperAdmin()(Nop[any, any])(ctx, req)
			if err != nil {
				return *new(Resp), err
			}

			return next(ctx, request)
		}
	}
}