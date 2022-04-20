package middlewares

import "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"

type ReqWithEmployee interface {
	GetEmployee() string
}

func UserIsEmployee[Req ReqWithEmployee, Resp any]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			userId, err := ctx.GetUserID()
			if err != nil {
				return *new(Resp), err
			}

			e := request.GetEmployee()

			if e != userId {
				return *new(Resp), NotAdmin
			}

			return next(ctx, request)
		}
	}
}