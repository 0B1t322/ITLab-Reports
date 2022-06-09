package toapprovereportsidgetter

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/service/salary"
	"github.com/RTUITLab/ITLab-Reports/transport/report/middlewares"
)

type toApproveReportsIdGetter struct {
	salaryService salary.SalaryService
}

func ToApproveReportsIdGetter(s salary.SalaryService) middlewares.ApprovedReportsIdsGetter {
	return &toApproveReportsIdGetter{salaryService: s}
}

func (t *toApproveReportsIdGetter) GetApprovedReportsIdsForUser(userId string, token string) ([]string, error) {
	return t.salaryService.GetApprovedReportsIds(
		context.Background(),
		token,
		*optional.NewOptional(userId),
	)
}

func (t *toApproveReportsIdGetter) GetApprovedReportsIds(token string) ([]string, error) {
	return t.salaryService.GetApprovedReportsIds(
		context.Background(),
		token,
		*optional.NewEmptyOptional[string](),
	)
}

