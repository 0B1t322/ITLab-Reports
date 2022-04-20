package endpoint_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/stretchr/testify/require"
)

type CreateSomeReq struct {
	name string
}

func (c CreateSomeReq) getName() string {
	return c.name
}

func (c *CreateSomeReq) SetName(name string) {
	c.name = name
}

type CreateSomeResp struct {
	newName string
}

type MapEndpoint struct {
	CreateSome endpoint.Endpoint[*CreateSomeReq, *CreateSomeResp]
}

type NameChangerReq interface {
	SetName(string)
}

func ChangeNameMiddleware[Req NameChangerReq, Resp any]() endpoint.Middleware[Req, Resp] {
		return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
			return func(
				ctx		context.Context,
				request		Req,
			) (response Resp, err error) {
				request.SetName("new_developer")
				return next(ctx, request)
			}
		}
}

func ChangeNameMiddlewareTo[Req NameChangerReq, Resp any](to string) endpoint.Middleware[Req, Resp] {
	return func(next endpoint.Endpoint[Req, Resp]) endpoint.Endpoint[Req, Resp] {
		return func(
			ctx		context.Context,
			request		Req,
		) (response Resp, err error) {
			request.SetName(to)
			fmt.Println(to)
			return next(ctx, request)
		}
	}
}

func TestFunc_Endpoint(t *testing.T) {

	mapEndpoints := MapEndpoint{
		CreateSome: func(ctx context.Context, req *CreateSomeReq) (*CreateSomeResp, error){
			return &CreateSomeResp{
				newName: fmt.Sprintf("hello %s", req.name),
			}, nil
		},
	}

	t.Run(
		"RunEndpoint",
		func(t *testing.T) {
			resp, err := mapEndpoints.CreateSome(
				context.Background(),
				&CreateSomeReq{
					name: "developer",
				},
			)
			require.NoError(t, err)

			require.Equal(
				t,
				&CreateSomeResp{
					newName: "hello developer",
				},
				resp,
			)
		},
	)

	t.Run(
		"AddMiddleware",
		func(t *testing.T) {
			copyMapEndpoints := mapEndpoints
			copyMapEndpoints.CreateSome.AddMiddleware(
				ChangeNameMiddleware[*CreateSomeReq, *CreateSomeResp](),
			)

			resp, _ := copyMapEndpoints.CreateSome(
				context.Background(),
				&CreateSomeReq{
					name: "developer",
				},
			)

			require.Equal(
				t,
				&CreateSomeResp{
					newName: "hello new_developer",
				},
				resp,
			)
		},
	)

	t.Run(
		"AddChainMiddleware",
		func(t *testing.T) {
			copyMapEndpoints := mapEndpoints
			copyMapEndpoints.CreateSome.AddMiddleware(
				endpoint.Chain(
					ChangeNameMiddlewareTo[*CreateSomeReq, *CreateSomeResp]("develop_1"),
					ChangeNameMiddlewareTo[*CreateSomeReq, *CreateSomeResp]("develop_2"),
					ChangeNameMiddlewareTo[*CreateSomeReq, *CreateSomeResp]("develop_3"),
				),
			)

			resp, _ := copyMapEndpoints.CreateSome(
				context.Background(),
				&CreateSomeReq{
					name: "developer",
				},
			)
			t.Log(resp.newName)
			require.Equal(
				t,
				&CreateSomeResp{
					newName: "hello develop_3",
				},
				resp,
			)
		},
	)

	t.Run(
		"ToGoKitEndpoint",
		func(t *testing.T) {
			resp, err := mapEndpoints.CreateSome.ToGoKitEndpoint()(context.Background(), &CreateSomeReq{
				name: "developer",
			},)
			require.NoError(t, err)

			require.IsType(t, &CreateSomeResp{}, resp)
		},
	)
}