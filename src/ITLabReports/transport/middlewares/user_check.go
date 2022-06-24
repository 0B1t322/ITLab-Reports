package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

type RespWithReporter interface {
	GetReporter() string
}

type RespWithImplementer interface {
	GetImplementer() string
}

type RespWithReporterAndImplementer interface {
	RespWithReporter
	RespWithImplementer
}

func CheckUserIsReporter[Req any, Resp RespWithReporter]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			resp, err := next(ctx, request)
			if err != nil {
				return resp, err
			}

			userId, err := ctx.GetUserID()
			if err != nil {
				return resp, err
			}

			if resp.GetReporter() != userId {
				return resp, NotAdmin
			}

			return resp, err
		}
	}
}

func CheckUserIsImplementer[Req any, Resp RespWithImplementer]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			resp, err := next(ctx, request)
			if err != nil {
				return resp, err
			}

			userId, err := ctx.GetUserID()
			if err != nil {
				return resp, err
			}

			if resp.GetImplementer() != userId {
				return resp, NotAdmin
			}

			return resp, err
		}
	}
}

func CheckUserIsReporterOrImplementer[Req any, Resp RespWithReporterAndImplementer]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			resp, err := next(ctx, request)
			if err != nil {
				return resp, err
			}

			userId, err := ctx.GetUserID()
			if err != nil {
				return resp, err
			}

			if resp.GetImplementer() == userId || resp.GetReporter() == userId {
				return resp, err
			}
			return resp, NotAdmin
		}
	}
}