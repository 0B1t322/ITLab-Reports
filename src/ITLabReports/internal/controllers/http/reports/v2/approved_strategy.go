package reports

import (
	"context"

	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v2/dto"
	reportsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab-Reports/internal/services/salary"
	"github.com/samber/do"
	"github.com/samber/mo"
)

type (
	ApprovedStateReq struct {
		State  dto.ApprovedState
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
	case dto.ApprovedState_Approved:
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
	case dto.ApprovedState_Not:
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
	case dto.ApprovedState_Approved:
		return reportsrepo.Query().
			Expression(
				reportsrepo.Expression().
					State(aggregate.ReportStatePaid),
			).
			Build(), nil
	case dto.ApprovedState_Not:
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
