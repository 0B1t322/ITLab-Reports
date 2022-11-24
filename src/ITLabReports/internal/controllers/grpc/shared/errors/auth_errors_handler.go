package errors

import (
	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthErrorsHandler struct{}

func (h *AuthErrorsHandler) Handle(err error) error {
	switch {
	case errors.Is(err, user.ErrTokenNotValid),
		err == user.ErrDontHaveUserID,
		err == user.ErrDontHaveRole,
		err == user.ErrDontHaveScope:
		return status.New(codes.Unauthenticated, err.Error()).Err()
	}

	return nil
}
