package grpc

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
)

// DecodeRequestFunc extracts a user-domain request object from a gRPC request.
// It's designed to be used in gRPC servers, for server-side endpoints. One
// straightforward DecodeRequestFunc could be something that decodes from the
// gRPC request message to the concrete request type.
type DecodeRequestFunc[GrpcReq any, Req any] func(context.Context, GrpcReq) (request Req, err error)

func (d DecodeRequestFunc[GrpcReq, Req]) ToGoKit() grpc.DecodeRequestFunc {
	return func(ctx context.Context, i interface{}) (request interface{}, err error) {
		return d(ctx, i.(GrpcReq))
	}
}

// EncodeResponseFunc encodes the passed response object to the gRPC response
// message. It's designed to be used in gRPC servers, for server-side endpoints.
// One straightforward EncodeResponseFunc could be something that encodes the
// object directly to the gRPC response message.
type EncodeResponseFunc[Resp any, GrpcResp any] func(context.Context, Resp) (response GrpcResp, err error)

func (e EncodeResponseFunc[Resp, GrpcResp]) ToGoKit() grpc.EncodeResponseFunc {
	return func(ctx context.Context, i interface{}) (response interface{}, err error) {
		return e(ctx, i.(Resp))
	}
}