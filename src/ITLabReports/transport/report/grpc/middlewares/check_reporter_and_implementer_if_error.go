package middlewares

import (
	. "github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
)

type RespWithReporterAndImplementerError interface {
	IsError() bool
	RespWithReporterAndImplementer
}

// Middleware need to not laucnh other middlewares that check responce
// In gRPC we can handle error in body with oneof
// If not error return "No Errors"
// Use with only RunMiddlewareIfAllFailed
func CheckUserIsReporterOrImplementerIfNotError[Req any, Resp RespWithReporterAndImplementerError]() MiddlewareWithContext[Req, Resp] {
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

			if resp.IsError() {
				return resp, nil
			} else if resp.GetImplementer() == userId || resp.GetReporter() == userId {
				return resp, nil
			}

			return resp, NotAdmin
		}
	}
}