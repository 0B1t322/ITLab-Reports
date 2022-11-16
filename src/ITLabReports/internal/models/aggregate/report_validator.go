package aggregate

import (
	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/common/utils"
	"github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"
	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	ErrReportValidation = errors.New("Report validation error")
)

type ReportValidator struct {
	utils.ValidationErrorsMerger
}

func NewReportValidator() ReportValidator {
	return ReportValidator{}
}

func (rv ReportValidator) Validate(report Report) error {
	var errors []error
	{
		errors = append(errors, rv.ValidateName(report.Name))
		errors = append(errors, rv.ValidateText(report.Text))
		errors = append(errors, rv.ValidateAssignees(report.Assignees))
	}

	return rv.Merge(errors...)
}

func (ReportValidator) ValidateName(name string) error {
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

func (ReportValidator) ValidateText(text string) error {
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

func (ReportValidator) ValidateAssignees(assignees valueobject.Assignees) error {
	err := validation.ValidateStruct(&assignees,
		validation.Field(&assignees.Reporter, validation.Required),
		validation.Field(&assignees.Implementer, validation.Required),
	)

	if err != nil {
		return validation.Errors{
			"assignees": err,
		}
	}

	return nil
}

func (r ReportValidator) Merge(errs ...error) error {
	if err := r.ValidationErrorsMerger.Merge(errs...); err != nil {
		return errors.Wrap(err, ErrReportValidation)
	}

	return nil
}
