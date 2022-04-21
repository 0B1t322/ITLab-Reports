package middlewares

import (
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

type ReqWithGetSetImplementer interface {
	SetImplementor(i string)
	GetImplementor() string 
}

func SetImplementerIfEmpty[Req ReqWithGetSetImplementer, Resp any]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			id, err := ctx.GetUserID()
			if err != nil {
				return *new(Resp), err
			}

			if request.GetImplementor() == "" {
				request.SetImplementor(id)
			}

			return next(ctx, request)
		}
	}
}