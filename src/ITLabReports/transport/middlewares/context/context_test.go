package context_test

import (
	"context"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
	"github.com/stretchr/testify/require"
)

type someEnd struct {
	end endpoint.Endpoint[string, string]
}

func TestFunc_Context(t *testing.T) {

	ends := someEnd{
		end: func(ctx context.Context, s string) (string, error) {
			return s, nil
		},
	}

	t.Run(
		"Cast",
		func(t *testing.T) {
			f := func(ctx context.Context) bool {
				_, ok := ctx.(mcontext.MiddlewareContext)
				return ok
			}

			mctx := mcontext.New(context.Background())
			require.True(t, f(mctx))
			require.False(t, f(context.Background()))
		},
	)

	t.Run(
		"AddMiddlewareToEndpoint",
		func(t *testing.T) {
			copyEnd := ends

			copyEnd.end.AddMiddleware(
				endpoint.Chain(
					SetToken[string, string](t).ToMiddleware(),
					GetToken[string, string](t).ToMiddleware(),
				),
			)

			resp, err := copyEnd.end(context.Background(), "string")
			require.NoError(t, err)
			t.Log(resp)
		},
	)
}


func GetToken[Req any, Resp any](t *testing.T) middlewares.MiddlewareWithContext[Req, Resp] {
	return func(next middlewares.EndpointWithContext[Req, Resp]) middlewares.EndpointWithContext[Req, Resp] {
		return func(ctx mcontext.MiddlewareContext, req Req) (Resp, error) {
			token, err := ctx.GetToken()
			require.NoError(t, err)
			require.Equal(t, "token", token)
			
			return next(ctx, req)
		}
	}
}

func SetToken[Req any, Resp any](t *testing.T) middlewares.MiddlewareWithContext[Req, Resp] {
	return func(next middlewares.EndpointWithContext[Req, Resp]) middlewares.EndpointWithContext[Req, Resp] {
		return func(ctx mcontext.MiddlewareContext, req Req) (Resp, error) {
			ctx.SetToken("token")
			return next(ctx, req)
		}
	}
}
