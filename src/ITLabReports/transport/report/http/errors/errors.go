package errors

import "errors"

var (
	ValidationError = errors.New("Validation error")
	DraftNotFound = errors.New("Draft not found")
	DraftIdNotValud = errors.New("Draft id is not valid")
)