package middlewares

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/go-kit/kit/log"
)

func LoggingMiddleware[ReqType any, RespType any](
	log log.Logger,
) endpoint.Middleware[ReqType, RespType] {
	return func(next endpoint.Endpoint[ReqType, RespType]) endpoint.Endpoint[ReqType, RespType] {
		return func(ctx context.Context, req ReqType) (resp RespType, err error) {
			resp, err = next(ctx, req)
			if err != nil {
				log.Log("Error", err)
			}
			return resp, err
		}
	}
}
