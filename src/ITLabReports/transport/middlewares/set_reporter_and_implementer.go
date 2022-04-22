package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

type ReqWithSetImplementorAndReporter interface{
	SetImplementerAndReporter(implementer, reporter string)
}

// Deprecated
func SetReporterAndImplementerIfFailed[Req ReqWithSetImplementorAndReporter, Resp any](
	m MiddlewareWithContext[Req, Resp],
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			_, err := m(Nop[Req, Resp])(ctx, request)
			if err != nil {
				id, err := ctx.GetUserID()
				if err != nil {
					return *new(Resp), err
				}

				request.SetImplementerAndReporter(id, id)
			}

			return next(ctx, request)
		}
	}
}

func SetReporterAndImplementer[Req ReqWithSetImplementorAndReporter, Resp any]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			id, err := ctx.GetUserID()
			if err != nil {
				return *new(Resp), err
			}
			request.SetImplementerAndReporter(id, id)

			return next(ctx, request)
		}
	}
}