package http

import (
	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewServer[Req any, Resp any](
	e endpoint.Endpoint[Req, Resp],
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	options ...kithttp.ServerOption,
) *kithttp.Server {
	return kithttp.NewServer(
		e.ToGoKitEndpoint(),
		dec.ToGoKit(),
		enc.ToGoKit(),
		options...,
	)
}

