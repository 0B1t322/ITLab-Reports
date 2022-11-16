package token

import "context"

type TestTokenService struct {
}

func NewTestTokenService() TokenService {
	return &TestTokenService{}
}

func (t *TestTokenService) RequestToken(ctx context.Context) (string, error) {
	return "test_token", nil
}
