package reports

import (
	"context"
	"errors"
	"time"

	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/0B1t322/RepoGen/pkg/queryexpression"
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/mo"
)

var (
	ErrIDIsNotValid   = errors.New("ID is not valid")
	ErrReportNotFound = errors.New("Report not found")
)

type GetReportsQuery struct {
	Filter FilterQuery

	Sort []SortFields

	Limit mo.Option[int64]

	Offset mo.Option[int64]
}

//go:generate go run -mod=mod github.com/0B1t322/RepoGen

//repogen:filter
type (
	FilterQuery = queryexpression.QueryExpression[FilterFields]

	FilterFields struct {
		ID          mo.Option[filter.FilterField[string]]
		IDs         mo.Option[filter.FilterField[[]string]]
		Name        mo.Option[filter.FilterField[string]]
		Date        mo.Option[filter.FilterField[time.Time]]
		Implementer mo.Option[filter.FilterField[string]]
		Reporter    mo.Option[filter.FilterField[string]]
		State       mo.Option[aggregate.ReportState]
	}
)

//repogen:sort
type SortFields struct {
	Name mo.Option[sortorder.SortOrder]
	Date mo.Option[sortorder.SortOrder]
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
	) (aggregate.Report, error)

	// CreateReport create report and return it
	// 	don't have catchable errors
	CreateReport(
		ctx context.Context,
		report *aggregate.Report,
	) error

	// DeleteReport delete report by id
	// 	catchable errors:
	// 		ErrIDIsNotValid
	// 		ErrReportNotFound
	DeleteReport(
		ctx context.Context,
		report aggregate.Report,
	) error

	// GetReports return reports acording to filters
	// 	don't have catchable errors
	GetReports(
		ctx context.Context,
		query GetReportsQuery,
	) ([]aggregate.Report, error)

	// UpdateReport update reports by id and not nil optionals
	// 	catchable errors:
	// 		ErrIDIsNotValid
	// 		ErrReportNotFound
	UpdateReport(
		ctx context.Context,
		report aggregate.Report,
	) error

	// CountByFilter count reports accroding to filter
	// 	don't have catchable errors
	CountByFilter(
		ctx context.Context,
		query GetReportsQuery,
	) (int64, error)
}
