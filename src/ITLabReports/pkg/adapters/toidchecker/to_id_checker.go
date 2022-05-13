package toidchecker

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	. "github.com/RTUITLab/ITLab-Reports/service/idvalidator"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
)

type toIdCheckerAdapter struct {
	s *Service
}

func ToIdChecker(s *Service) middlewares.IdChecker {
	return &toIdCheckerAdapter{s}
}

type GetTokener interface {
	GetToken() (string, error)
}

func (t *toIdCheckerAdapter) CheckIds(token string, ids []string) error {
	return t.s.CheckUserIds(
		context.Background(),
		token,
		ids,	
	)
}

func (t *toIdCheckerAdapter) GetIncorrectId(err error) string {
	if err == nil || (err != nil && errors.Is(err, ErrFailedToValidateId)){
		return ""
	}

	return errors.Unwrap(err).Error()
}

func (t *toIdCheckerAdapter) IsIncorrectIdError(err error) bool {
	return errors.Is(err, ErrIdInvalid)
}
