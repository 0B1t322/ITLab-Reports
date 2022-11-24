package app

import (
	shared "github.com/RTUITLab/ITLab-Reports/internal/controllers/shared/reports"
	"github.com/samber/do"
)

func (a *App) ConfigureSharedControllersOptions() {
	// Configure reports approved strategy
	do.Provide(
		a.injector,
		func(i *do.Injector) (shared.ApprovedStateStrategy, error) {
			return shared.NewExternalApprovedStrategyFrom(i)
		},
	)
}
