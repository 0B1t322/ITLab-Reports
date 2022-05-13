package idvalidator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/service/idvalidator"
	"github.com/stretchr/testify/require"
)

type MockIdValidatorWithoutErrors struct {}

func (MockIdValidatorWithoutErrors) ValidateIds(
	ctx context.Context,
	token string,
	ids []string,
) (error) {
	return nil
}

type MockIdValidatorUnwrapError struct {}

func (MockIdValidatorUnwrapError) ValidateIds(
	ctx context.Context,
	token string,
	ids []string,
) (error) {
	return fmt.Errorf("Some error")
}

type MockIdValidatorInvalidIdError struct {}

func (MockIdValidatorInvalidIdError) ValidateIds(
	ctx context.Context,
	token string,
	ids []string,
) (error) {
	return errors.Wrap(fmt.Errorf("invalid_id"), idvalidator.ErrIdInvalid)
}

func TestFunc_IdValidator(t *testing.T) {
	t.Run(
		"WithoutErrors",
		func(t *testing.T) {
			validator := idvalidator.New(&MockIdValidatorWithoutErrors{})
			err := validator.CheckUserIds(context.Background(), "", []string{})
			require.NoError(t, err)
		},
	)

	t.Run(
		"WithUnwrapError",
		func(t *testing.T) {
			validator := idvalidator.New(&MockIdValidatorUnwrapError{})
			err := validator.CheckUserIds(context.Background(), "", []string{})
			require.Condition(
				t,
				func() (success bool) {
					return errors.Is(err, idvalidator.ErrFailedToValidateId) && errors.Unwrap(err).Error() == "Some error"
				},
			)
		},
	)

	t.Run(
		"WithInvalidIdError",
		func(t *testing.T) {
			validator := idvalidator.New(&MockIdValidatorInvalidIdError{})
			err := validator.CheckUserIds(context.Background(), "", []string{})
			require.Condition(
				t,
				func() (success bool) {
					return errors.Is(err, idvalidator.ErrIdInvalid) && errors.Unwrap(err).Error() == "invalid_id"
				},
			)
		},
	)
}