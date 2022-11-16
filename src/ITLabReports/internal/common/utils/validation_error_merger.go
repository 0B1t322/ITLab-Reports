package utils

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ValidationErrorsMerger struct {
}

func (ValidationErrorsMerger) Merge(errs ...error) error {
	errors := validation.Errors{}
	{
		for _, err := range errs {
			if err != nil {
				if valid, ok := err.(validation.Errors); ok {
					for k, v := range valid {
						errors[k] = v
					}
				}
			}
		}
	}

	return errors.Filter()
}
