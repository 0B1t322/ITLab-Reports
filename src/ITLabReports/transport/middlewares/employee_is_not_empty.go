package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

func EmployeeIsNotEmpty[Req ReqWithEmployee, Resp any](
	errBuilder func() error,
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {

			if request.GetEmployee() == "" {
				return *new(Resp), errBuilder()
			}

			return next(ctx, request)
		}
	}
}