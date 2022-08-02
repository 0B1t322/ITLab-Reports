package servers

import (
	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleErrors(err error) error {

	var code codes.Code
	switch {
	// INVALID_ARGUMENT
	case errors.Is(err, reports.ErrGetReportsBadParams):
		code = codes.InvalidArgument
	// Permission
	case err == middlewares.NotAdmin, err == middlewares.NotSuperAdmin:
		code = codes.PermissionDenied
	// Auth
	case err == middlewares.FailedToParseToken, err == middlewares.RoleNotFound,
		err == middlewares.TokenNotValid, err == middlewares.TokenExpired:
		code = codes.Unauthenticated
	case err != nil:
		code = codes.Internal
		logrus.WithFields(
			logrus.Fields{
				"err": err,
				"service": "reports",
			},
		).Error()
	}
	return status.Error(code, err.Error())
}