package reports

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/pkg/errors"
)

var (
	ErrReportNotFound            = errors.New("Report not found")
	ErrDraftNotFound             = errors.New("Draft not found")
	ErrReportsIsAlreadyPaid      = errors.New("Reports is already paid")
	ErrCantCreateReportFromDraft = errors.New("Cant create report from draft")
	ErrCantSetReportPaid         = errors.New("Cant set report paid")
	ErrCantGetReport             = errors.New("Cant get report")

	ErrFailedToCreateReport          = errors.New("Failed to create report")
	ErrFailedToGetReports            = errors.New("Failed to get reports")
	ErrFailedToGetReport             = errors.New("Failed to get report")
	ErrFailedToCountReports          = errors.New("Failed to count reports")
	ErrFailedToCreateReportFromDraft = errors.New("Failed to create report from draft")
	ErrFailedToSetReportPaid         = errors.New("Failed to set report paid")
)

// ReportsService interface represent ReportService methods
type ReportsService interface {
	// GetReport return report by id
	// 	catchable errors:
	// 		ErrReportNotFound
	GetReport(
		ctx context.Context,
		req GetReportReq,
	) (aggregate.Report, error)

	// GetReports return reports acording to filters
	//
	// don't have catchable errors
	GetReports(
		ctx context.Context,
		req GetReportsReq,
	) ([]aggregate.Report, error)

	// CreateReport create report and return it
	// 	catchable errors:
	// 		ErrValidationError as target
	// Target errors catch by:
	// 		errors.Is(err, ErrValidationError)
	CreateReport(
		ctx context.Context,
		req CreateReportReq,
	) (aggregate.Report, error)

	// CountReport count report according to filter and return count
	//
	// don't have catchable errors
	CountReports(
		ctx context.Context,
		req GetReportsReq,
	) (int64, error)

	// CreateReportFrom draft get draft and create report from it
	//
	// throws errors:
	//
	// 	1. ErrFailedCreateReportFromDraft
	//
	//  2. wrapped ErrReportValidation from aggregate package
	//
	//  3. wrapped ErrCantCreateReportFromDraft
	//
	// If some uncatchable error occured, it will be wrapped by ErrFailedCreateReportFromDraft
	CreateReportFromDraft(
		ctx context.Context,
		req CreateReportFromDraftReq,
	) (aggregate.Report, error)

	// SetReportPaid
	//
	// throws errors:
	//
	// 	1. ErrReportsIsAlreadyPaid
	//
	//  2. ErrReportNotFound
	//
	//  3. wrapped ErrCantSetReportPaid
	//
	// If some uncatchable error occured, it will be wrapped by ErrFailedToSetReportPaid
	SetReportPaid(
		ctx context.Context,
		req SetReportPaidReq,
	) error
}
