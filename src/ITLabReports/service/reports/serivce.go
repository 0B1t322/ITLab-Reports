package reports

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/pkg/errors"
)

var (
	ErrReportNotFound = errors.New("Report not found")
	ErrReportIDNotValid = errors.New("Report id is not valid")
	ErrValidationError = errors.New("Report is not valid")
)

// Service interface represet ReportService methods
type Service interface {
	// GetReport return report by id
	// 	catchable errors:
	// 		ErrReportIDNotValid
	// 		ErrReportNotFound
	GetReport(
		ctx context.Context,
		id	string,
	) (*report.Report, error)
	
	// DeleteReport delete report by id
	// 	catchable errors:
	// 		ErrReportIDNotValid
	// 		ErrReportNotFound
	DeleteReport(
		ctx	context.Context,
		id string,
	) error
	
	// UpdateReport update reports by id and not nil optionals
	// Name, Text, Implemtner fields can't be empty
	// 	catchable errors:
	// 		ErrReportIDNotValid
	// 		ErrReportNotFound
	// 		ErrValidationError as target
	// Target errors catch by:
	// 		errors.Is(err, ErrValidationError)
	UpdateReport(
		ctx context.Context,
		id string,
		params reportdomain.UpdateReportParams,
	) (*report.Report, error)
	
	// GetReports return reports acording to filters
	// 
	// don't have catchable errors
	GetReports(
		ctx	context.Context,
		params *reportdomain.GetReportsParams,
	) ([]*report.Report, error)
	
	// CreateReport create report and return it
	// 	catchable errors:
	// 		ErrValidationError as target
	// Target errors catch by:
	// 		errors.Is(err, ErrValidationError)
	CreateReport(
		ctx context.Context,
		report *report.Report,
	) (*report.Report, error)

	// CountReport count report according to filter and return count
	// 
	// don't have catchable errors
	CountReports(
		ctx	context.Context,
		filter *reportdomain.GetReportsFilterFieldsWithOrAnd,
	) (int64, error)
}