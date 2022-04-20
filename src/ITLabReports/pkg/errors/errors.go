package errors

import (
	"errors"
	"fmt"
)

type MyError struct {
	base error
	with error
}

func (m MyError) Error() string {
	return fmt.Sprintf("%s: %s", m.with, m.base)
}

func Wrap(err error, with error) error {
	return &MyError{
		base: err,
		with: with,
	}
}

func Unwrap(err error) error {
	if err, ok := err.(*MyError); ok {
		return err.base
	} else {
		return errors.Unwrap(err)
	}
}

// Check all wrapings errors to this target
func Is(err error, target error) bool {
	if err, ok := err.(*MyError); ok {
		if err.base == target || err.with == target {
			return true
		} else if Is(err.base, target) {
			return true
		} else if Is(err.with, target) {
			return true
		} else {
			return false
		}
	} else if err == target {
		return true
	} else {
		return errors.Is(err, target)
	}
}

func New(text string) error {
	return errors.New(text)
}