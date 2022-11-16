package aggregate

import (
	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/common/utils"
	"github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrDraftValidation = errors.New("Draft validation error")
)

type DraftValidator struct {
	utils.ValidationErrorsMerger
}

func NewDraftValidator() DraftValidator {
	return DraftValidator{}
}

func (dv DraftValidator) Validate(draft Draft) error {
	var errors []error
	{
		errors = append(errors, dv.ValidateName(draft.Name))
		errors = append(errors, dv.ValidateText(draft.Text))
		errors = append(errors, dv.ValidateAssignees(draft.Assignees))
	}

	return dv.Merge(errors...)
}

func (DraftValidator) ValidateName(name string) error {
	err := validation.Validate(
		name,
		validation.Required,
	)
	if err != nil {
		return validation.Errors{
			"name": err,
		}
	}

	return nil
}

func (DraftValidator) ValidateText(text string) error {
	err := validation.Validate(
		text,
		validation.Required,
	)
	if err != nil {
		return validation.Errors{
			"text": err,
		}
	}

	return nil
}

func (DraftValidator) ValidateAssignees(assignees valueobject.Assignees) error {
	err := validation.ValidateStruct(&assignees,
		validation.Field(&assignees.Reporter, validation.Required),
	)

	if err != nil {
		return validation.Errors{
			"assignees": err,
		}
	}

	return nil
}

func (d DraftValidator) Merge(errs ...error) error {
	if err := d.ValidationErrorsMerger.Merge(errs...); err != nil {
		return errors.Wrap(err, ErrDraftValidation)
	}

	return nil
}
