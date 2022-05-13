package idvalidator

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
)

var (
	ErrIdInvalid          = errors.New("Id is invalid")
	ErrFailedToValidateId = errors.New("Failed to validate")
)

// Struct because it a helper service
//
type Service struct {
	idValidator IdsValidator
}

type IdsValidator interface {
	// If id is ivalid error should be wrap with ErrIdInvalid and invalid id in
	ValidateIds(ctx context.Context, token string, ids []string) error
}

func New(
	validator IdsValidator,
) *Service {
	s := &Service{
		idValidator: validator,
	}

	return s
}

func (s *Service) CheckUserIds(
	ctx context.Context,
	token string,
	ids []string,
) error {
	err := s.idValidator.ValidateIds(ctx, token, ids)
	if errors.Is(err, ErrIdInvalid) {
		return err
	} else if err != nil {
		return errors.Wrap(err, ErrFailedToValidateId)
	}

	return nil
}
