package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type CustomEndpoint[Req any, Resp any] interface {
	ToEndpoint() Endpoint[Req, Resp]
}

type Endpoint[ReqType any, RespType any] func(ctx context.Context, req ReqType) (responce RespType, err error)

func (e Endpoint[ReqType, RespType]) ToEndpoint() Endpoint[ReqType, RespType] {
	return e
}

func(e Endpoint[ReqType, RespType]) ToGoKitEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return e(ctx, request.(ReqType))
	}
}

func (e *Endpoint[ReqType, RespType]) AddMiddleware(m Middleware[ReqType, RespType]) Endpoint[ReqType, RespType] {
	*e = m(*e)
	return *e
}

func (e *Endpoint[ReqType, RespType]) AddMiddlwares(outer Middleware[ReqType, RespType], others ...Middleware[ReqType, RespType]) Endpoint[ReqType, RespType] {
	ms := Chain(outer, others...)
	*e = ms(*e)
	return *e
}

func (e *Endpoint[ReqType, RespType]) AddCustomMiddlewares(outer CustomMiddleware[ReqType, RespType], others ...CustomMiddleware[ReqType, RespType]) Endpoint[ReqType, RespType] {
	ms := CustomChain(outer, others...)
	*e = ms.ToMiddleware()(*e)
	return *e
}

func (e *Endpoint[ReqType, RespType]) AddCustomMiddleware(m CustomMiddleware[ReqType, RespType]) Endpoint[ReqType, RespType] {
	e.AddCustomMiddlewares(m)
	return *e
}

func (e *Endpoint[ReqType, RespType]) AddGoKitMiddleware(m endpoint.Middleware) Endpoint[ReqType, RespType] {
	kitEndpoint := e.ToGoKitEndpoint()
	kitEndpoint = m(kitEndpoint)
	*e = EndpointFromGoKitEndpoint[ReqType, RespType](kitEndpoint)
	return *e
}

func (e *Endpoint[ReqType, RespType]) AddGoKitMiddlewares(
	outer endpoint.Middleware, 
	others ...endpoint.Middleware,
) Endpoint[ReqType, RespType] {
	m := endpoint.Chain(
		outer,
		others...,
	)
	*e = EndpointFromGoKitEndpoint[ReqType, RespType](m(e.ToGoKitEndpoint()))

	return *e
}

func EndpointFromGoKitEndpoint[ReqType any, RespType any](e endpoint.Endpoint) Endpoint[ReqType, RespType] {
	return func(ctx context.Context, req ReqType) (responce RespType, err error) {
		kitResp, err := e(ctx, req)
		if err != nil {
			return *new(RespType), err
		}

		responce = kitResp.(RespType)

		return responce, err
	}
}

type Middleware[ReqType any, RespType any] func(Endpoint[ReqType, RespType]) Endpoint[ReqType, RespType]

func (m Middleware[ReqType, RespType]) ToMiddleware() Middleware[ReqType, RespType] {
	return m
}

func (m Middleware[ReqType, RespType]) ToGoKitMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return m(EndpointFromGoKitEndpoint[ReqType, RespType](e)).ToGoKitEndpoint()
	}
}

func MiddlewareFromGoKitMiddleware[ReqType any, RespType any](m endpoint.Middleware) Middleware[ReqType, RespType] {
	return func(e Endpoint[ReqType, RespType]) Endpoint[ReqType, RespType] {
		return EndpointFromGoKitEndpoint[ReqType, RespType](
			m(
				e.ToGoKitEndpoint(),
			),
		)
	}
}

type CustomMiddleware[Req any, Resp any] interface {
	ToMiddleware() Middleware[Req, Resp]
}

func Chain[ReqType any, RespType any](outer Middleware[ReqType, RespType], others ...Middleware[ReqType, RespType]) Middleware[ReqType, RespType] {
	return MiddlewareFromGoKitMiddleware[ReqType, RespType](
		endpoint.Chain(
			outer.ToGoKitMiddleware(),

			func (others ...Middleware[ReqType, RespType]) (slice []endpoint.Middleware) {
				for _, other := range others {
					slice = append(slice, other.ToGoKitMiddleware())
				}
				return slice
			}(others...)...,
		),
	)
}

func CustomChain[Req any, Resp any](outer CustomMiddleware[Req, Resp], others ...CustomMiddleware[Req, Resp]) CustomMiddleware[Req, Resp] {
	newOthers := func(others ...CustomMiddleware[Req, Resp]) (slice []Middleware[Req, Resp]) {
		for _, other := range others {
			slice = append(slice, other.ToMiddleware())
		}
		return slice
	}(others...)

	return Chain(outer.ToMiddleware(), newOthers...)
}