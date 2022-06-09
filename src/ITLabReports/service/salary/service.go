package salary

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
)

type SalaryService interface {
	GetApprovedReportsIds(
		ctx context.Context,
		token string,
		userId optional.Optional[string],
	) ([]string, error)
}