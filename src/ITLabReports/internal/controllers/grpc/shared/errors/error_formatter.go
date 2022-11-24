package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorFormatter struct{}

func (ErrorFormatter) FormatError(code codes.Code, err error) error {
	return status.New(code, err.Error()).Err()
}
