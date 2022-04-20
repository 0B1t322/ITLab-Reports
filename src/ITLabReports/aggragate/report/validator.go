package report

type IReport interface{
	GetName() string
	GetText() string
	GetImplementer() string
	GetReporter() string
}

type Validator interface {
	Validate() error
}

type ReportValidateOptions func(r IReport) error

// Try to cast IReport to Validator interface
// 
// If can use it for validate
// 
// If not do nothing
// 
// Validator interface:
// 	type Validator interface {
// 		Validate() error
// 	}
func WithSelfValidator() ReportValidateOptions {
	return func(r IReport) error {
		if rv, ok := r.(Validator); ok {
			return rv.Validate()
		}
		return nil
	}
}

func WithValidateAllFields() ReportValidateOptions {
	return MergeValidateOptions(
		WithValidateName(),
		WithValidateText(),
		WithValidateImplementor(),
		WithValidateReporter(),
	)
}

// Merge into And conditions
func MergeValidateOptions(opts ...ReportValidateOptions) ReportValidateOptions {
	return func(r IReport) error {
		for _, opt := range opts {
			if err := opt(r); err != nil {
				return err
			}
		}

		return nil
	}
}

func WithValidateName() ReportValidateOptions {
	return func(r IReport) error {
		if r.GetName() == "" {
			return ErrNameEmpty
		}
		return nil
	}
}

func WithValidateText() ReportValidateOptions {
	return func(r IReport) error {
		if r.GetText() == "" {
			return ErrTextEmpty
		}
		return nil
	}
}

func WithValidateImplementor() ReportValidateOptions {
	return func(r IReport) error {
		if r.GetImplementer() == "" {
			return ErrImplementorEmpty
		}
		return nil
	}
}

func WithValidateReporter() ReportValidateOptions {
	return func(r IReport) error {
		if r.GetReporter() == "" {
			return ErrReporterEmpty
		}
		return nil
	}
}

// NewReportValidator return validator with accroding opts
// 
// If options is nil use option
// 
// 	WithValidateAllFields()
func NewReportValidator(
	opts ...ReportValidateOptions,
) ReportValidator {
	r := reportValidator{}

	if len(opts) == 0 {
		r.validator = WithValidateAllFields()
	} else {
		r.validator = MergeValidateOptions(opts...)
	}

	return r
}

type ReportValidator interface {
	Validate(IReport) error
}

type reportValidator struct{
	validator ReportValidateOptions
}

func (rv reportValidator) Validate(r IReport) error {
	if r == nil {
		return ErrReportIsNil
	}

	return rv.validator(r)
}