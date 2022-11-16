package drafts

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
)

var (
	ErrDraftNotFound   = errors.New("Draft not found")
	ErrCantGetDraft    = errors.New("Can't get draft")
	ErrCantUpdateDraft = errors.New("Can't update draft")
	ErrCantDeleteDraft = errors.New("Can't delete draft")

	ErrFailedToGetDraft    = errors.New("Failed to get draft")
	ErrFailedToCreateDraft = errors.New("Failed to create draft")
	ErrFailedToUpdateDraft = errors.New("Failed to update draft")
	ErrFailedToDeleteDraft = errors.New("Failed to delete draft")
	ErrFailedToGetDrafts   = errors.New("Failed to get drafts")
	ErrFailedToCountDrafts = errors.New("Failed to count drafts")
)

type DraftsService interface {
	/*
		GetDraft return draft

		Only author of draft can get it

		throws errors:
			1. ErrDraftNotFound
			2. wrapped ErrCantGetDraft
		If some unexpected error occurred, throws ErrFailedToGetDraft
	*/
	GetDraft(
		ctx context.Context,
		req GetDraftReq,
	) (aggregate.Draft, error)

	/*
		CreateDraft create draft

		throws errors:
			1. wrapped ErrDraftValidation from aggregate package

		If some unexpected error occurred, throws ErrFailedToCreateDraft
	*/
	CreateDraft(
		ctx context.Context,
		req CreateDraftReq,
	) (aggregate.Draft, error)

	/*
		UpdateDraft update draft

		Only author of draft can update it

		throws errors:
			1. ErrDraftNotFound
			2. wrapped ErrCantUpdateDraft
			3. wrapped ErrDraftValidation from aggregate package

		If some unexpected error occurred, throws ErrFailedToUpdateDraft
	*/
	UpdateDraft(
		ctx context.Context,
		req UpdateDraftReq,
	) (aggregate.Draft, error)

	/*
		DeleteDraft delete draft

		Only author of draft can delete it

		throws errors:
			1. ErrDraftNotFound
			2. wrapped ErrCantDeleteDraft

		If some unexpected error occurred, throws ErrFailedToDeleteDraft
	*/
	DeleteDraft(
		ctx context.Context,
		req DeleteDraftReq,
	) error

	/*
		GetDrafts return drafts

		Return only drafts of requested user

		If some unexpected error occurred, throws ErrFailedToGetDrafts
	*/
	GetDrafts(
		ctx context.Context,
		req GetDraftsReq,
	) ([]aggregate.Draft, error)

	/*
		CountDrafts return count of drafts

		Count only drafts of requested user

		If some unexpected error occurred, throws ErrFailedToCountDrafts
	*/
	CountDrafts(
		ctx context.Context,
		req GetDraftsReq,
	) (int64, error)
}
