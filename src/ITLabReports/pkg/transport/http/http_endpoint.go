package http

import (
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
)

type HTTPEndpoint[
	HTTPReq any, 
	Req any, 
	HTTPResp, 
	Resp any,
] endpoint.Endpoint[Req, Resp]

func (h HTTPEndpoint[HTTPReq, Req, HTTPResp, Resp]) ToGoKit() kitendpoint.Endpoint {
	return h.ToGoKit()
}

func (e HTTPEndpoint[HTTPReq, Req, HTTPResp, Resp]) NewHandler(
	dec DecodeRequestFunc[Req],
	enc EncodeResponseFunc[Resp],
	opts ...kithttp.ServerOption,
) http.Handler {
	return kithttp.NewServer(
		e.ToGoKit(),
		dec.ToGoKit(),
		enc.ToGoKit(),
		opts...,
	)
}