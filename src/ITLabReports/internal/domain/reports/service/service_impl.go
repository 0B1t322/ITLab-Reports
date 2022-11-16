package reports

import (
	"context"
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/do"
)

type (
	ReportsServicePermissionChecker interface {
		// CanSetPaid check if user can set report paid
		//
		// throws errors:
		//
		// 1. wrapped ErrCantSetReportPaid
		CanSetPaid(ctx context.Context, report aggregate.Report, user aggregate.User) error

		// CanGetReport check if user can get this report
		//
		// throws errors:
		//
		//  1. wrapped CantGetReports
		CanGetReport(ctx context.Context, report aggregate.Report, user aggregate.User) error
	}

	DraftService interface {
		// GetDraft return draft by id
		//
		// catchable errors:
		//
		// 	1. ErrDraftNotFound
		//
		//  2. wrapped ErrCantCreateReportFromDraft
		GetDraft(ctx context.Context, id string, by aggregate.User) (aggregate.Draft, error)

		// DeleteDraft delete draft by id
		//
		// catchable errors:
		//
		// 1. ErrDraftNotFound
		//
		// 2. wrapped ErrCantDeleteDraft
		DeleteDraft(ctx context.Context, id string, by aggregate.User) error
	}
)

type ReportsServiceImpl struct {
	repo              reports.ReportRepository
	draftService      DraftService
	permissionChecker ReportsServicePermissionChecker
}

func NewReportsServiceImpl(
	repo reports.ReportRepository,
	draftService DraftService,
) *ReportsServiceImpl {
	s := &ReportsServiceImpl{
		repo:              repo,
		draftService:      draftService,
		permissionChecker: NewInternalPermissionChecker(),
	}

	return s
}

func NewReportsServiceImplFrom(i *do.Injector) (*ReportsServiceImpl, error) {
	repo := do.MustInvoke[reports.ReportRepository](i)
	draftService := do.MustInvoke[DraftService](i)

	return NewReportsServiceImpl(repo, draftService), nil
}

// GetReport return report by id
//
// catchable errors:
//  1. ErrReportNotFound
func (r *ReportsServiceImpl) GetReport(
	ctx context.Context,
	req GetReportReq,
) (aggregate.Report, error) {
	report, err := r.repo.GetReport(
		ctx,
		req.ID,
	)
	if err == reports.ErrIDIsNotValid || err == reports.ErrReportNotFound {
		return aggregate.Report{}, ErrReportNotFound
	} else if err != nil {
		return aggregate.Report{}, errors.Wrap(err, ErrFailedToGetReport)
	}

	if err := r.permissionChecker.CanGetReport(ctx, report, req.By); err != nil {
		return aggregate.Report{}, err
	}

	return report, nil
}

// GetReports return reports acording to filters
//
// don't have catchable errors
func (r *ReportsServiceImpl) GetReports(
	ctx context.Context,
	req GetReportsReq,
) ([]aggregate.Report, error) {
	reports, err := r.repo.GetReports(
		ctx,
		req.Query,
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedToGetReports)
	}

	return reports, nil
}

// CreateReport create report and return it
//
//	catchable errors:
//		ErrValidationError as target
//
// Target errors catch by:
//
//	errors.Is(err, ErrValidationError)
func (r *ReportsServiceImpl) CreateReport(
	ctx context.Context,
	req CreateReportReq,
) (aggregate.Report, error) {
	report, err := aggregate.NewReport(
		req.Name,
		req.Text,
		req.By.ID,
		req.Implementer.OrElse(req.By.ID),
	)
	if err != nil {
		return aggregate.Report{}, err
	}

	if err := r.repo.CreateReport(ctx, &report); err != nil {
		return aggregate.Report{}, errors.Wrap(err, ErrFailedToCreateReport)
	}

	return report, nil
}

// CountReport count report according to filter and return count
//
// don't have catchable errors
func (r *ReportsServiceImpl) CountReports(
	ctx context.Context,
	req GetReportsReq,
) (int64, error) {
	count, err := r.repo.CountByFilter(
		ctx,
		req.Query,
	)
	if err != nil {
		return 0, errors.Wrap(err, ErrFailedToCountReports)
	}

	return count, nil
}

// CreateReportFrom draft get draft and create report from it
//
// throws errors:
//
//  1. ErrFailedCreateReportFromDraft
//
//  2. wrapped ErrReportValidation from aggregate package
//
//  3. wrapped ErrCantCreateReportFromDraft
//
// If some uncatchable error occured, it will be wrapped by ErrFailedCreateReportFromDraft
func (r *ReportsServiceImpl) CreateReportFromDraft(
	ctx context.Context,
	req CreateReportFromDraftReq,
) (aggregate.Report, error) {
	draft, err := r.draftService.GetDraft(ctx, req.DraftID, req.By)
	if err == ErrDraftNotFound || err == ErrCantCreateReportFromDraft {
		return aggregate.Report{}, err
	} else if err != nil {
		return aggregate.Report{}, errors.Wrap(err, ErrFailedToCreateReportFromDraft)
	}

	report := draft.ToReport()
	// Update date of report
	report.Date = time.Now().UTC()

	if err := r.repo.CreateReport(
		ctx,
		&report,
	); err != nil {
		return aggregate.Report{}, errors.Wrap(err, ErrFailedToCreateReportFromDraft)
	}

	// Don't check error because if we can get it we can delete it
	if err := r.draftService.DeleteDraft(ctx, req.DraftID, req.By); err != nil {
		return aggregate.Report{}, errors.Wrap(err, ErrFailedToCreateReportFromDraft)
	}

	return report, nil
}

// SetReportPaid
//
// throws errors:
//
//  1. ErrReportsIsAlreadyPaid
//
//  2. ErrReportNotFound
//
//  3. wrapped ErrCantSetReportPaid
//
// If some uncatchable error occured, it will be wrapped by ErrFailedToSetReportPaid
func (r *ReportsServiceImpl) SetReportPaid(
	ctx context.Context,
	req SetReportPaidReq,
) error {
	report, err := r.GetReport(
		ctx,
		GetReportReq{
			ID: req.ID,
		},
	)
	if err == ErrReportNotFound {
		return err
	} else if err != nil {
		return errors.Wrap(errors.Unwrap(err), ErrFailedToSetReportPaid)
	}

	if err := r.permissionChecker.CanSetPaid(ctx, report, req.By); err != nil {
		return err
	}

	report.SetPaid()

	if err := r.repo.UpdateReport(ctx, report); err != nil {
		return errors.Wrap(err, ErrFailedToSetReportPaid)
	}

	return nil
}
