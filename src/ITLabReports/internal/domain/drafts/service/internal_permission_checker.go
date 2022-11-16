package drafts

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

func (i *internalPermissionChecker) CanGetDraft(
	ctx context.Context,
	draft aggregate.Draft,
	user aggregate.User,
) error {
	if !user.IsAdminOrSuperAdmin() && !draft.UserIsDraftOwner(user) {
		return errors.Wrap(fmt.Errorf("You are not admin"), ErrCantGetDraft)
	}

	return nil
}

func (i *internalPermissionChecker) CanUpdateDraft(
	ctx context.Context,
	draft aggregate.Draft,
	user aggregate.User,
) error {
	if !user.IsAdminOrSuperAdmin() && !draft.UserIsDraftOwner(user) {
		return errors.Wrap(fmt.Errorf("You are not admin"), ErrCantUpdateDraft)
	}

	return nil
}

func (i *internalPermissionChecker) CanDeleteDraft(
	ctx context.Context,
	draft aggregate.Draft,
	user aggregate.User,
) error {
	if !user.IsAdminOrSuperAdmin() && !draft.UserIsDraftOwner(user) {
		return errors.Wrap(fmt.Errorf("You are not admin"), ErrCantDeleteDraft)
	}

	return nil
}
