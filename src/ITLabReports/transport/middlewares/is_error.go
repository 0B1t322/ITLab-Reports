package middlewares

import "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"

type RespWithIsError interface {
	IsError() bool
}

// Middleware need to not laucnh other middlewares that check responce
// In gRPC we can handle error in body with oneof
// If not error return "No Errors"
// Use with only RunMiddlewareIfAllFailed
func CheckIsError[Req any, Resp RespWithIsError]() MiddlewareWithContext[Req, Resp] {
	return func(next EndpointWithContext[Req, Resp]) EndpointWithContext[Req, Resp] {
		return func(ctx context.MiddlewareContext, request Req) (Resp, error) {
			resp, err := next(ctx, request)
			if err != nil {
				return resp, err
			}

			if resp.IsError() {
				return resp, nil
			}

			return resp, nil
		}
	}
}
