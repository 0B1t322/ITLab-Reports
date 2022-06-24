package grpc

import (
	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/go-kit/kit/transport/grpc"
)

func NewServer[GrpcReq any, GrpcResp any, Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[GrpcReq, Req],
	enc EncodeResponseFunc[Resp, GrpcResp],
	opts ...grpc.ServerOption,
) *grpc.Server {
	return grpc.NewServer(
		e.ToGoKitEndpoint(),
		dec.ToGoKit(),
		enc.ToGoKit(),
		opts...,
	)
}