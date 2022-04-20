package http

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

// DecodeRequestFunc extracts a user-domain request object from an HTTP
// request object. It's designed to be used in HTTP servers, for server-side
// endpoints. One straightforward DecodeRequestFunc could be something that
// JSON decodes from the request body to the concrete request type.
type DecodeRequestFunc[Req any] func(context.Context, *http.Request) (request Req, err error)

type RequestDecoder[Req any] interface {
	DecodeRequest(context.Context, *http.Request) (request Req, err error)
}

func (d DecodeRequestFunc[Req]) ToGoKit() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		return d(ctx, r)
	}
}

// EncodeRequestFunc encodes the passed request object into the HTTP request
// object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward EncodeRequestFunc could be something that JSON
// encodes the object directly to the request body.
type EncodeRequestFunc[Req any] func(context.Context, *http.Request, Req) error

func (e EncodeRequestFunc[Req]) ToGoKit() kithttp.EncodeRequestFunc {
	return func(ctx context.Context, r *http.Request, i interface{}) error {
		return e(ctx,r, i.(Req))
	}
}

// CreateRequestFunc creates an outgoing HTTP request based on the passed
// request object. It's designed to be used in HTTP clients, for client-side
// endpoints. It's a more powerful version of EncodeRequestFunc, and can be used
// if more fine-grained control of the HTTP request is required.
type CreateRequestFunc[Req any] func(context.Context, Req) (*http.Request, error)

func (c CreateRequestFunc[Req]) ToGoKit() kithttp.CreateRequestFunc {
	return func(ctx context.Context, i interface{}) (*http.Request, error) {
		return c(ctx, i.(Req))
	}
}

// EncodeResponseFunc encodes the passed response object to the HTTP response
// writer. It's designed to be used in HTTP servers, for server-side
// endpoints. One straightforward EncodeResponseFunc could be something that
// JSON encodes the object directly to the response body.
type EncodeResponseFunc[Resp any] func(context.Context, http.ResponseWriter, Resp) error

func EncodeResponseFuncFromEncoder[Resp any] (enc ResponceEncoder[Resp]) EncodeResponseFunc[Resp] {
	return enc.EncodeResponce
}

type ResponceEncoder[Resp any] interface {
	EncodeResponce(context.Context, http.ResponseWriter, Resp) error
}

func (e EncodeResponseFunc[Resp]) ToGoKit() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, i interface{}) error {
		return e(ctx, w, i.(Resp))
	}
}

// DecodeResponseFunc extracts a user-domain response object from an HTTP
// response object. It's designed to be used in HTTP clients, for client-side
// endpoints. One straightforward DecodeResponseFunc could be something that
// JSON decodes from the response body to the concrete response type.
type DecodeResponseFunc[Resp any] func(context.Context, *http.Response) (response Resp, err error)

func (d DecodeResponseFunc[Resp]) ToGoKit() kithttp.DecodeResponseFunc {
	return func(ctx context.Context, r *http.Response) (response interface{}, err error) {
		return d(ctx, r)
	}
}

type MarshalHTTPReqToEndpointReq[HTTPReq any, Req any] func(ctx context.Context, req HTTPReq) (Req, error)

func (m MarshalHTTPReqToEndpointReq[HTTPReq, Req]) MarshalHttpReq(ctx context.Context, req HTTPReq) (Req, error) {
	return m(ctx, req)
}

type HTTPReqMarshaller[HTTPReq any, Req any] interface {
	MarshalHttpReq(ctx context.Context, req HTTPReq) (Req, error) 
}

type MarshallEndpointRespToHTTPResponce[Resp any, HTTPResp any] func(ctx context.Context, resp Resp) (HTTPResp, error)

func (m MarshallEndpointRespToHTTPResponce[Resp, HTTPResp]) MarshalHttpResp(ctx context.Context, resp Resp) (HTTPResp, error) {
	return m(ctx, resp)
}

type RespToHTTPRespMarhaller[Resp any, HTTPResp any] interface {
	MarshalHttpResp(ctx context.Context, resp Resp) (HTTPResp, error)
}

type RequestEncoder[HTTPReq any, Req any] struct{
	httpMarshaller HTTPReqMarshaller[HTTPReq, Req]
}