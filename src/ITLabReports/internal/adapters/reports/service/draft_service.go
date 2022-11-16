package reports

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/do"
)

type ReportsDraftService struct {
	draftsService drafts.DraftsService
}

func NewReportsDraftService(draftsService drafts.DraftsService) *ReportsDraftService {
	return &ReportsDraftService{
		draftsService: draftsService,
	}
}

func NewReportsDraftServiceFrom(i *do.Injector) (*ReportsDraftService, error) {
	draftService := do.MustInvoke[drafts.DraftsService](i)

	return NewReportsDraftService(draftService), nil
}

func (rds *ReportsDraftService) GetDraft(
	ctx context.Context,
	id string,
	by aggregate.User,
) (aggregate.Draft, error) {
	draft, err := rds.draftsService.GetDraft(
		ctx,
		drafts.GetDraftReq{
			ID: id,
			By: by,
		},
	)
	if err == drafts.ErrDraftNotFound {
		return aggregate.Draft{}, reports.ErrDraftNotFound
	} else if errors.Is(err, drafts.ErrCantGetDraft) {
		return aggregate.Draft{}, errors.Wrap(err, reports.ErrCantCreateReportFromDraft)
	} else if err != nil {
		return aggregate.Draft{}, err
	}

	return draft, nil
}

func (rds *ReportsDraftService) DeleteDraft(
	ctx context.Context,
	id string,
	by aggregate.User,
) error {
	err := rds.draftsService.DeleteDraft(
		ctx,
		drafts.DeleteDraftReq{
			ID: id,
			By: by,
		},
	)
	if err == drafts.ErrDraftNotFound {
		return reports.ErrDraftNotFound
	} else if errors.Is(err, drafts.ErrCantDeleteDraft) {
		return err
	} else if err != nil {
		return err
	}

	return nil
}
