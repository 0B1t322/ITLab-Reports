package report

import (
	"context"
	"errors"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
)

var (
	ErrIDIsNotValid = errors.New("ID is not valid")
	ErrReportNotFound = errors.New("Report not found")
)

type GetReportsParams struct {
	Filter *GetReportsFilter `swaggerignore:"true"`

	Limit  optional.Optional[int64] `swaggerignore:"true"`

	Offset optional.Optional[int64] `swaggerignore:"true"`
}

type GetReportsFilter struct {
	GetReportsFilterFieldsWithOrAnd `swaggerignore:"true"`

	GetReportsSort `swaggerignore:"true"`
}

type GetReportsFilterFieldsWithOrAnd struct {
	GetReportsFilterFields

	Or	[]*GetReportsFilterFieldsWithOrAnd
	And []*GetReportsFilterFieldsWithOrAnd
}

type GetReportsFilterFields struct {
	ReportID	*filter.FilterField[string]

	Name		*filter.FilterField[string]

	Date		*filter.FilterField[string]

	Implementer	*filter.FilterField[string]

	Reporter	*filter.FilterField[string]
}

type GetReportsSort struct {
	NameSort	optional.Optional[ordertype.OrderType]

	DateSort	optional.Optional[ordertype.OrderType]
}

type UpdateReportParams struct {
	Name	optional.Optional[string]

	Text	optional.Optional[string]

	Implementer optional.Optional[string]
}

// ReportRepository is interface of all ReportRepositorys
type ReportRepository interface {
	// GetReport return report by id
	// 	catchable errors:
	// 		ErrIDIsNotValid
	// 		ErrReportNotFound
	GetReport(
		ctx context.Context,
		id string,
	) (*report.Report, error)
	
	// CreateReport create report and return it
	// 	don't have catchable errors
	CreateReport(
		ctx context.Context,
		report *report.Report,
	) (*report.Report, error)
	
	// DeleteReport delete report by id
	// 	catchable errors:
	// 		ErrIDIsNotValid
	// 		ErrReportNotFound
	DeleteReport(
		ctx	context.Context,
		id string,
	) error
	
	// GetReports return reports acording to filters
	// 	don't have catchable errors
	GetReports(
		ctx context.Context,
		params *GetReportsParams,
	) ([]*report.Report, error)
	
	// UpdateReport update reports by id and not nil optionals
	// 	catchable errors:
	// 		ErrIDIsNotValid
	// 		ErrReportNotFound
	UpdateReport(
		ctx	context.Context,
		id	string,
		params UpdateReportParams,
	) (*report.Report, error)
	
	// CountByFilter count reports accroding to filter
	// 	don't have catchable errors
	CountByFilter(
		ctx context.Context,
		params *GetReportsFilterFieldsWithOrAnd,
	) (int64, error)
}
