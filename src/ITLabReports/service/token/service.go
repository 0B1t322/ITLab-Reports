package token

import "context"

type TokenService interface {
	RequestToken(ctx context.Context) (string, error)
}
