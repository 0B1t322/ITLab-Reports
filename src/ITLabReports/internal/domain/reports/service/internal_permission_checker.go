package reports

import (
	"context"
	"fmt"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
)

type internalPermissionChecker struct {
}

func NewInternalPermissionChecker() *internalPermissionChecker {
	return &internalPermissionChecker{}
}

// CanSetPaid check if user can set report paid
//
// throws errors:
//
// 1. wrapped ErrCantSetReportPaid
//
// 2. ErrReportsIsAlreadyPaid
func (i *internalPermissionChecker) CanSetPaid(
	ctx context.Context,
	report aggregate.Report,
	user aggregate.User,
) error {
	if report.State == aggregate.ReportStatePaid {
		return ErrReportsIsAlreadyPaid
	}

	if !user.IsAdminOrSuperAdmin() {
		return errors.Wrap(fmt.Errorf("You are not admin"), ErrCantSetReportPaid)
	}

	return nil
}

// CanGetReport check if user can get this report
//
// throws errors:
//
//  1. wrapped CantGetReports
func (i *internalPermissionChecker) CanGetReport(
	ctx context.Context,
	report aggregate.Report,
	user aggregate.User,
) error {
	if !user.IsAdminOrSuperAdmin() && !report.UserIsReportOwner(user) &&
		!report.ReportAboutUser(user) {
		return errors.Wrap(fmt.Errorf("You are not admin"), ErrCantGetReport)
	}

	return nil
}
