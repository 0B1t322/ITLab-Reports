package report

import (
	"context"
	"errors"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	rm "github.com/RTUITLab/ITLab-Reports/entity/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/samber/mo"
)

var (
	ErrIDIsNotValid   = errors.New("ID is not valid")
	ErrReportNotFound = errors.New("Report not found")
)

type GetReportsParams struct {
	Filter *GetReportsFilter `swaggerignore:"true"`

	Limit mo.Option[int64] `swaggerignore:"true"`

	Offset mo.Option[int64] `swaggerignore:"true"`
}

type GetReportsFilter struct {
	GetReportsFilterFieldsWithOrAnd `swaggerignore:"true"`

	SortParams []GetReportsSort `swaggerignore:"true"`
}

type GetReportsFilterFieldsWithOrAnd struct {
	GetReportsFilterFields

	Or  []*GetReportsFilterFieldsWithOrAnd
	And []*GetReportsFilterFieldsWithOrAnd
}

type GetReportsFilterFields struct {
	ReportID *filter.FilterField[string]

	ReportsId *filter.FilterField[[]string]

	Name *filter.FilterField[string]

	Date *filter.FilterField[string]

	Implementer *filter.FilterField[string]

	Reporter *filter.FilterField[string]

	State *filter.FilterField[rm.ReportState]
}

type GetReportsSort struct {
	NameSort mo.Option[ordertype.OrderType]

	DateSort mo.Option[ordertype.OrderType]
}

type UpdateReportParams struct {
	Name mo.Option[string]

	Text mo.Option[string]

	Implementer mo.Option[string]
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
		ctx context.Context,
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
		ctx context.Context,
		id string,
		params UpdateReportParams,
	) (*report.Report, error)

	// CountByFilter count reports accroding to filter
	// 	don't have catchable errors
	CountByFilter(
		ctx context.Context,
		params *GetReportsFilterFieldsWithOrAnd,
	) (int64, error)
}
