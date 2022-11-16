package salary

import (
	"context"

	"github.com/samber/mo"
)

type testModeSalaryService struct{}

func NewTestModeSalaryService() SalaryService {
	return &testModeSalaryService{}
}

func (t *testModeSalaryService) GetApprovedReportsIds(
	ctx context.Context,
	token string,
	userId mo.Option[string],
) ([]string, error) {
	return []string{}, nil
}
