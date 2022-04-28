package servererrorencoder

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	derr "github.com/RTUITLab/ITLab-Reports/transport/draft/http/errors"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/clarketm/json"
	"github.com/sirupsen/logrus"
)

type ErrorModel struct {
	Error string	`json:"error"`
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	var statusCode int

	switch {
	// Forbidden errors
	case err == middlewares.NotAdmin, err == middlewares.NotSuperAdmin, err == middlewares.YouAreNotOwner:
		statusCode = http.StatusForbidden
	// Unauth
	case 	err == middlewares.FailedToParseToken, err == middlewares.RoleNotFound, 
			err == middlewares.TokenNotValid, err == middlewares.TokenExpired:
		statusCode = http.StatusUnauthorized
	// BadRequest 
	case 	err == derr.DraftIDIsInvalid, errors.Is(err, reports.ErrValidationError):
		statusCode = http.StatusBadRequest
	// NotFound
	case err == derr.DraftNotFound:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
		logrus.WithFields(
			logrus.Fields{
				"from": "draft",
				"err": err,
			},
		).Error("")
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorModel{Error: err.Error()})
}