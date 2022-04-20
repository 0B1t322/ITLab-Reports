package middlewares

import "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"

func Nop[Req any, Resp any](ctx context.MiddlewareContext, req Req) (Resp, error) {
	return *new(Resp), nil
}

// One of middlewares should be without error
//
// It's be good if all this middlewares return similat error because if all of them return error, will return only last failed middleware
func MergeMiddlewaresIntoOr[Req any, Resp any](
	ms ...MiddlewareWithContext[Req, Resp],
) MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(
			ctx context.MiddlewareContext, 
			request Req,
		) (Resp, error) {
			var err error
			var success bool = false

			for _, m := range ms {
				_, err = m(Nop[Req,Resp])(ctx, request)
				if err == nil {
					success = true
					break
				}
			}

			if success {
				return next(ctx, request)
			}

			return *new(Resp), err
		}
	}
}