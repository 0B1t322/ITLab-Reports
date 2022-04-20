package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

type RespWithReporterAndImplementor interface {
	GetReporter() string

	GetImplementer() string
}

func CheckUserIsReporter[Req any, Resp RespWithReporterAndImplementor]() MiddlewareWithContext[Req, Resp] {
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