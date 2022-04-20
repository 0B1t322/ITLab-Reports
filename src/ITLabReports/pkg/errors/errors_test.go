package errors_test

import (
	stderror "errors"
	"fmt"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestFunc(t *testing.T) {
	err1 := stderror.New("Err1")
	err2 := stderror.New("Err2")

	err := errors.Wrap(err1, err2)
	t.Run(
		"Wrap",
		func(t *testing.T) {
			require.Equal(
				t,
				fmt.Sprintf("%s: %v", err2, err1),
				err.Error(),
			)
		},
	)

	t.Run(
		"Unwrap",
		func(t *testing.T) {
			require.Equal(
				t,
				err1,
				errors.Unwrap(err),
			)
		},
	)

	t.Run(
		"Is",
		func(t *testing.T) {
			require.Equal(
				t,
				true,
				errors.Is(err, err1),
			)

			require.Equal(
				t,
				true,
				errors.Is(err, err2),
			)

			require.Equal(
				t,
				false,
				errors.Is(err, fmt.Errorf("AnotherError")),
			)
		},
	)
}