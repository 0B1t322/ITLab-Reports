package drafts

import (
	"context"

	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/0B1t322/RepoGen/pkg/queryexpression"
	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/adapters/sortorder"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/mo"
)

var (
	ErrDraftIDNotValid = errors.New("Draft id is not valid")
	ErrDraftNotFound   = errors.New("Draft not found")
)

type GetDraftsQuery struct {
	Filter FilterQuery

	Sort []SortFields

	Limit mo.Option[int64]

	Offset mo.Option[int64]
}

//go:generate go run -mod=mod github.com/0B1t322/RepoGen

//repogen:filter
type (
	FilterQuery = queryexpression.QueryExpression[FilterFields]

	FilterFields struct {
		ID          mo.Option[filter.FilterField[string]]
		Implementer mo.Option[filter.FilterField[string]]
		Reporter    mo.Option[filter.FilterField[string]]
	}
)

//repogen:sort
type SortFields struct {
	Name mo.Option[sortorder.SortOrder]
	Date mo.Option[sortorder.SortOrder]
}

type DraftRepository interface {
	// GetDraft return draft by id
	// 	catchable errors:
	//
	// 		1. ErrDraftIDNotValid
	//
	// 		2. ErrDraftNotFound
	//
	GetDraft(
		ctx context.Context,
		id string,
	) (aggregate.Draft, error)

	// GetDrafts return drafts by query
	GetDrafts(
		ctx context.Context,
		query GetDraftsQuery,
	) ([]aggregate.Draft, error)

	CountDrafts(
		ctx context.Context,
		query GetDraftsQuery,
	) (int64, error)

	// CreateDraft create draft
	// 	don't have catchable errors
	CreateDraft(
		ctx context.Context,
		draft *aggregate.Draft,
	) error

	// UpdateDraft update draft
	// 		catchable errors:
	//
	// 		1. ErrDraftIDNotValid
	//
	// 		2. ErrDraftNotFound
	UpdateDraft(
		ctx context.Context,
		draft aggregate.Draft,
	) error

	// DeleteDraft delete draft
	// 		catchable errors:
	//
	// 		1. ErrDraftIDNotValid
	//
	// 		2. ErrDraftNotFound
	DeleteDraft(
		ctx context.Context,
		draft aggregate.Draft,
	) error
}
