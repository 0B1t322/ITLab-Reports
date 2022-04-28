package errors

import "github.com/RTUITLab/ITLab-Reports/pkg/errors"

var (
	DraftNotFound = errors.New("Draft not found")

	DraftIDIsInvalid = errors.New("Draft id is invalid")
)