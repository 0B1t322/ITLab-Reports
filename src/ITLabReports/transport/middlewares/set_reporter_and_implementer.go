package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
	"github.com/go-kit/kit/endpoint"
)

type ReqWithSetImplementorAndReporter interface{
	SetImplementerAndReporter(implementer, reporter string)
}

func SetReporterAndImplementerIfFailed[Req ReqWithSetImplementorAndReporter, Resp any](
	m endpoint.Middleware,
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			var req any = request

			_, err := m(endpoint.Nop)(ctx, req)
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