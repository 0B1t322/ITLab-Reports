package salary

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
)

type testModeSalaryService struct{}

func NewTestModeSalaryService() SalaryService {
	return &testModeSalaryService{}
}

func (t *testModeSalaryService) GetApprovedReportsIds(
	ctx context.Context,
	token string,
	userId optional.Optional[string],
) ([]string, error) {
	return []string{}, nil
}