package errors

import (
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/gin-gonic/gin"
)

type AuthErrorHandler struct {
	Formatter ErrorFormatter
}

func NewAuthErrorHandler(formatter ErrorFormatter) AuthErrorHandler {
	return AuthErrorHandler{
		Formatter: formatter,
	}
}

// HandleError handle the auth error
//
// If it was auth error return true
func (a AuthErrorHandler) HandlerError(c *gin.Context, err error) bool {
	switch {
	case errors.Is(err, user.ErrTokenNotValid),
		err == user.ErrDontHaveUserID,
		err == user.ErrDontHaveRole,
		err == user.ErrDontHaveScope:
		a.Formatter.FormatError(c, err, http.StatusUnauthorized)
		return true
	}

	return false
}
