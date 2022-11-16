package drafts

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"
	"github.com/samber/do"
)

type DraftsServicePermissionsChecker interface {
	/*
		Checks if user has permission to get draft

		throws errors:
			1. wrapped ErrCantGetDraft
	*/
	CanGetDraft(ctx context.Context, draft aggregate.Draft, user aggregate.User) error

	/*
		Checks if user has permission to update draft

		throws errors:
			1. wrapped ErrCantUpdateDraft
	*/
	CanUpdateDraft(ctx context.Context, draft aggregate.Draft, user aggregate.User) error

	/*
		Checks if user has permission to delete draft

		throws errors:
			1. wrapped ErrCantDeleteDraft
	*/
	CanDeleteDraft(ctx context.Context, draft aggregate.Draft, user aggregate.User) error
}

type DraftsServiceImpl struct {
	repo              drafts.DraftRepository
	permissionChecker DraftsServicePermissionsChecker
}

func NewDraftsServiceImpl(
	repo drafts.DraftRepository,
) *DraftsServiceImpl {
	s := &DraftsServiceImpl{
		repo:              repo,
		permissionChecker: NewInternalPermissionChecker(),
	}

	return s
}

func NewDraftsServiceImplFrom(i *do.Injector) (*DraftsServiceImpl, error) {
	repo := do.MustInvoke[drafts.DraftRepository](i)

	return NewDraftsServiceImpl(repo), nil
}

/*
GetDraft return draft

# Only author of draft can get it

throws errors:
 1. ErrDraftNotFound
 2. wrapped ErrCantGetDraft

If some unexpected error occurred, throws ErrFailedToGetDraft
*/
func (d *DraftsServiceImpl) GetDraft(
	ctx context.Context,
	req GetDraftReq,
) (aggregate.Draft, error) {
	draft, err := d.repo.GetDraft(
		ctx,
		req.ID,
	)
	if err == drafts.ErrDraftNotFound || err == drafts.ErrDraftIDNotValid {
		return aggregate.Draft{}, ErrDraftNotFound
	} else if err != nil {
		return aggregate.Draft{}, errors.Wrap(err, ErrFailedToGetDraft)
	}

	err = d.permissionChecker.CanGetDraft(ctx, draft, req.By)
	if err != nil {
		return aggregate.Draft{}, errors.Wrap(err, ErrCantGetDraft)
	}

	return draft, nil
}

/*
CreateDraft create draft

throws errors:
 1. wrapped ErrDraftValidation from aggregate package

If some unexpected error occurred, throws ErrFailedToCreateDraft
*/
func (d *DraftsServiceImpl) CreateDraft(
	ctx context.Context,
	req CreateDraftReq,
) (aggregate.Draft, error) {
	draft, err := aggregate.NewDraft(
		req.Name,
		req.Text,
		req.By.ID,
		req.Implementer.OrElse(req.By.ID),
	)
	if err != nil {
		return aggregate.Draft{}, err
	}

	err = d.repo.CreateDraft(ctx, &draft)
	if err != nil {
		return aggregate.Draft{}, errors.Wrap(err, ErrFailedToCreateDraft)
	}

	return draft, nil
}

/*
UpdateDraft update draft

# Only author of draft can update it

throws errors:
 1. ErrDraftNotFound
 2. wrapped ErrCantUpdateDraft
 3. wrapped ErrDraftValidation from aggregate package

If some unexpected error occurred, throws ErrFailedToUpdateDraft
*/
func (d *DraftsServiceImpl) UpdateDraft(
	ctx context.Context,
	req UpdateDraftReq,
) (aggregate.Draft, error) {
	draft, err := d.repo.GetDraft(
		ctx,
		req.ID,
	)
	if err == drafts.ErrDraftNotFound || err == drafts.ErrDraftIDNotValid {
		return aggregate.Draft{}, ErrDraftNotFound
	} else if err != nil {
		return aggregate.Draft{}, errors.Wrap(err, ErrFailedToUpdateDraft)
	}

	err = d.permissionChecker.CanUpdateDraft(ctx, draft, req.By)
	if err != nil {
		return aggregate.Draft{}, errors.Wrap(err, ErrCantUpdateDraft)
	}

	var (
		validator   = aggregate.NewDraftValidator()
		validErrors []error
	)
	{
		req.Name.ForEach(
			func(name string) {
				if err := validator.ValidateName(name); err != nil {
					validErrors = append(validErrors, err)
				} else {
					draft.SetName(name)
				}
			},
		)

		req.Text.ForEach(
			func(text string) {
				if err := validator.ValidateText(text); err != nil {
					validErrors = append(validErrors, err)
				} else {
					draft.SetText(text)
				}
			},
		)

		req.Implementer.ForEach(
			func(implementerID string) {
				if err := validator.ValidateAssignees(
					valueobject.Assignees{Reporter: draft.Assignees.Reporter, Implementer: implementerID},
				); err != nil {
					validErrors = append(validErrors, err)
				} else {
					draft.SetAssignees(
						draft.Assignees.Reporter,
						implementerID,
					)
				}
			},
		)
	}
	if err := validator.Merge(validErrors...); err != nil {
		return aggregate.Draft{}, err
	}

	err = d.repo.UpdateDraft(
		ctx,
		draft,
	)
	if err != nil {
		return aggregate.Draft{}, errors.Wrap(err, ErrFailedToUpdateDraft)
	}

	return draft, nil
}

/*
DeleteDraft delete draft

# Only author of draft can delete it

throws errors:
 1. ErrDraftNotFound
 2. wrapped ErrCantDeleteDraft

If some unexpected error occurred, throws ErrFailedToDeleteDraft
*/
func (d *DraftsServiceImpl) DeleteDraft(ctx context.Context, req DeleteDraftReq) error {
	draft, err := d.repo.GetDraft(
		ctx,
		req.ID,
	)
	if err == drafts.ErrDraftNotFound || err == drafts.ErrDraftIDNotValid {
		return ErrDraftNotFound
	} else if err != nil {
		return errors.Wrap(err, ErrFailedToDeleteDraft)
	}

	err = d.permissionChecker.CanDeleteDraft(ctx, draft, req.By)
	if err != nil {
		return err
	}

	err = d.repo.DeleteDraft(
		ctx,
		draft,
	)
	if err != nil {
		return errors.Wrap(err, ErrFailedToDeleteDraft)
	}

	return nil
}

/*
GetDrafts return drafts

# Return only drafts of requested user

If some unexpected error occurred, throws ErrFailedToGetDrafts
*/
func (d *DraftsServiceImpl) GetDrafts(
	ctx context.Context,
	req GetDraftsReq,
) ([]aggregate.Draft, error) {
	drafts, err := d.repo.GetDrafts(
		ctx,
		req.Query,
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrFailedToGetDrafts)
	}

	return drafts, nil
}

/*
CountDrafts return count of drafts

# Count only drafts of requested user

If some unexpected error occurred, throws ErrFailedToCountDrafts
*/
func (d *DraftsServiceImpl) CountDrafts(
	ctx context.Context,
	req GetDraftsReq,
) (int64, error) {
	count, err := d.repo.CountDrafts(
		ctx,
		req.Query,
	)

	if err != nil {
		return 0, errors.Wrap(err, ErrFailedToCountDrafts)
	}

	return count, nil
}
