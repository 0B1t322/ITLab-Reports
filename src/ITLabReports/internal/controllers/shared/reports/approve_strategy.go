package reports

import (
	"context"

	"github.com/0B1t322/RepoGen/pkg/filter"
	reportsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab-Reports/internal/services/salary"
	"github.com/samber/do"
	"github.com/samber/mo"
)

type (
	ApprovedStateReq struct {
		State  aggregate.ReportState
		Token  string
		UserID mo.Option[string]
	}

	ApprovedStateStrategy interface {
		SetApproved(req ApprovedStateReq) (reportsrepo.FilterQuery, error)
	}
)

type ExternalApprovedStrategy struct {
	salaryService salary.SalaryService
}

func NewExternalApprovedStrategy(salaryService salary.SalaryService) *ExternalApprovedStrategy {
	return &ExternalApprovedStrategy{
		salaryService: salaryService,
	}
}

func NewExternalApprovedStrategyFrom(i *do.Injector) (*ExternalApprovedStrategy, error) {
	return NewExternalApprovedStrategy(
		do.MustInvoke[salary.SalaryService](i),
	), nil
}

func (e ExternalApprovedStrategy) SetApproved(
	req ApprovedStateReq,
) (query reportsrepo.FilterQuery, err error) {
	switch req.State {
	case aggregate.ReportStatePaid:
		ids, err := e.salaryService.GetApprovedReportsIds(
			context.Background(),
			req.Token,
			req.UserID,
		)
		if err != nil {
			return query, err
		}

		return reportsrepo.Query().
			Expression(
				reportsrepo.Expression().
					IDs(ids, filter.IN),
			).
			Build(), nil
	case aggregate.ReportStateCreated:
		ids, err := e.salaryService.GetApprovedReportsIds(
			context.Background(),
			req.Token,
			req.UserID,
		)
		if err != nil {
			return query, err
		}

		return reportsrepo.Query().
			Expression(
				reportsrepo.Expression().
					IDs(ids, filter.NIN),
			).
			Build(), nil
	default:
		return
	}
}

type InternalApprovedStrategy struct{}

func NewInternalApprovedStrategy() *InternalApprovedStrategy {
	return &InternalApprovedStrategy{}
}

func (i InternalApprovedStrategy) SetApproved(
	req ApprovedStateReq,
) (query reportsrepo.FilterQuery, err error) {
	switch req.State {
	case aggregate.ReportStatePaid:
		return reportsrepo.Query().
			Expression(
				reportsrepo.Expression().
					State(aggregate.ReportStatePaid),
			).
			Build(), nil
	case aggregate.ReportStateCreated:
		return reportsrepo.Query().
			Expression(
				reportsrepo.Expression().
					State(aggregate.ReportStateCreated),
			).
			Build(), nil
	default:
		return
	}
}
