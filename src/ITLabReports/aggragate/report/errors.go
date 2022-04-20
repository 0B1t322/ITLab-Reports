package report

import "errors"

var (
	// Validate errors
	ErrNameEmpty = errors.New("Name can't be empty")
	ErrReporterEmpty = errors.New("Reporter can't be empty")
	ErrImplementorEmpty = errors.New("Implementer can't be empty")
	ErrTextEmpty = errors.New("Text can't be empty")
	ErrReportIsNil = errors.New("Report is nil")
)