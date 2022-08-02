package salary

import (
	"context"

	"github.com/samber/mo"
)

type SalaryService interface {
	GetApprovedReportsIds(
		ctx context.Context,
		token string,
		userId mo.Option[string],
	) ([]string, error)
}