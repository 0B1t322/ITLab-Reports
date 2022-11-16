package token

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type ExternalTokenService struct {
	clientID         string
	clientSecret     string
	tokenEndpointUrl string
}

func NewExternalTokenService(
	clientID string,
	clientSecret string,
	tokenEndpointURL string,
) (TokenService, error) {
	return &ExternalTokenService{
		clientID:         clientID,
		clientSecret:     clientSecret,
		tokenEndpointUrl: tokenEndpointURL,
	}, nil
}

func (e *ExternalTokenService) RequestToken(ctx context.Context) (string, error) {
	cfg := clientcredentials.Config{
		ClientID:     e.clientID,
		ClientSecret: e.clientSecret,
		TokenURL:     e.tokenEndpointUrl,
		Scopes:       []string{},
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	t, err := cfg.Token(ctx)
	if err != nil {
		return "", err
	}

	return t.AccessToken, nil
}
