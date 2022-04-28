package servererrorencoder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	rerr "github.com/RTUITLab/ITLab-Reports/transport/report/http/errors"
	serr "github.com/RTUITLab/ITLab-Reports/transport/report/http/errors"
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
	case err == middlewares.NotAdmin, err == middlewares.NotSuperAdmin:
		statusCode = http.StatusForbidden
	case err == middlewares.YouAreNotOwner:
		statusCode = http.StatusForbidden
		err = fmt.Errorf("You are not owner of this draft")
	// Unauth
	case 	err == middlewares.FailedToParseToken, err == middlewares.RoleNotFound, 
			err == middlewares.TokenNotValid, err == middlewares.TokenExpired:
		statusCode = http.StatusUnauthorized
	// BadRequest 
	case 	err == reports.ErrReportIDNotValid, errors.Is(err, reports.ErrValidationError), 
			err == rerr.DraftIdNotValud,
			errors.Is(err, serr.ValidationError):
		statusCode = http.StatusBadRequest
	// NotFound
	case err == reports.ErrReportNotFound, err == rerr.DraftNotFound:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
		logrus.WithFields(
			logrus.Fields{
				"from": "reports",
				"err": err,
			},
		).Error("")
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorModel{Error: err.Error()})
}