package middlewares

import "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"

type ReqWithSetUserID interface {
	SetUserID(string)
}

func SetUserID[Req ReqWithSetUserID, Resp any]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			userId, err := ctx.GetUserID()
			if err != nil {
				return *new(Resp), err
			}
			request.SetUserID(userId)

			return next(ctx, request)
		}
	}
}