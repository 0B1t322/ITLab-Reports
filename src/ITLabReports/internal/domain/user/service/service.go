package user

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
)

var (
	ErrTokenExpired       = errors.New("Token expired")
	ErrFailedToParseToken = errors.New("Failed to parse token")
	ErrTokenNotValid      = errors.New("Token not valid")
	ErrDontHaveScope      = errors.New("Don't have scope")
	ErrDontHaveRole       = errors.New("Don't have role")
	ErrDontHaveUserID     = errors.New("Don't have user id")
)

type UserService interface {
	/*
		AuthUser checks user's token and returns user's info

		bearerToken - user's token with bearer prefix

		throws errors:
			1. wrapped ErrTokenNotValid with reasons:

				1.1 ErrTokenExpired - if token expired

				1.2 ErrFailedToParseToken - if token not valid

			2. ErrDontHaveScope - if token don't have needed scope
			3. ErrDontHaveRole - if token don't have needed role
			4. ErrDontHaveUserID - if token don't have user id in specifed claim or in client_id
	*/
	AuthUser(
		ctx context.Context,
		// JWT token with bearer prefix
		bearerToken string,
	) (aggregate.User, error)
}
