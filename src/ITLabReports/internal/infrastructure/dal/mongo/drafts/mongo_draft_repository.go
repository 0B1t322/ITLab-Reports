package drafts

import (
	"context"

	"github.com/0B1t322/MongoBuilder/operators/query"
	"github.com/0B1t322/RepoGen/pkg/filter"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/models"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/utils"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDraftRepository struct {
	collection *mongo.Collection
	utils.IDChecker
}

func NewMongoDraftRepository(db *mongo.Database) *MongoDraftRepository {
	return &MongoDraftRepository{
		collection: db.Collection("drafts"),
		IDChecker:  utils.NewIDChecker(drafts.ErrDraftIDNotValid),
	}
}

func NewMongoDraftRepositoryFrom(i *do.Injector) (*MongoDraftRepository, error) {
	db := do.MustInvoke[*mongo.Database](i)

	return NewMongoDraftRepository(db), nil
}

// GetDraft return draft by id
// catchable errors:
//
//  1. ErrDraftIDNotValid
//
//  2. ErrDraftNotFound
func (m *MongoDraftRepository) GetDraft(ctx context.Context, id string) (aggregate.Draft, error) {
	_, err := m.ParseID(id)
	if err != nil {
		return aggregate.Draft{}, err
	}

	get, err := m.GetDrafts(
		ctx,
		drafts.GetDraftsQuery{
			Filter: drafts.Query().
				Expression(
					drafts.Expression().ID(id, filter.EQ),
				).
				Build(),
		},
	)
	if err != nil {
		return aggregate.Draft{}, err
	}

	if len(get) == 0 {
		return aggregate.Draft{}, drafts.ErrDraftNotFound
	}

	return get[0], nil
}

// GetDrafts return drafts by query
func (m *MongoDraftRepository) GetDrafts(
	ctx context.Context,
	query drafts.GetDraftsQuery,
) ([]aggregate.Draft, error) {
	cur, err := m.collection.Aggregate(ctx, MarshallQuery(query))
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var drafts []models.Draft
	{
		if err := cur.All(ctx, &drafts); err != nil {
			return nil, err
		}
	}

	return lo.Map(
		drafts,
		func(draft models.Draft, _ int) aggregate.Draft {
			return models.DraftFromModel(draft)
		},
	), nil
}

func (m *MongoDraftRepository) CountDrafts(
	ctx context.Context,
	query drafts.GetDraftsQuery,
) (int64, error) {
	cur, err := m.collection.Aggregate(ctx, MarshallQueryForCount(query))
	if err != nil {
		return 0, err
	}

	defer cur.Close(ctx)

	type Count struct {
		Count int64 `bson:"count"`
	}

	var count Count
	{
		cur.Next(ctx)

		if err := cur.Decode(&count); err != nil {
			return 0, err
		}
	}

	return count.Count, nil
}

// CreateDraft create draft
//
//	don't have catchable errors
func (m *MongoDraftRepository) CreateDraft(ctx context.Context, draft *aggregate.Draft) error {
	draft.ID = primitive.NewObjectID().Hex()

	draftModel, err := models.NewDraftModel(*draft)
	if err != nil {
		return err
	}

	_, err = m.collection.InsertOne(ctx, draftModel)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDraft update draft
//
//	catchable errors:
//
//	1. ErrDraftIDNotValid
//
//	2. ErrDraftNotFound
func (m *MongoDraftRepository) UpdateDraft(ctx context.Context, draft aggregate.Draft) error {
	id, err := m.ParseID(draft.ID)
	if err != nil {
		return err
	}

	draftModel, err := models.NewDraftModel(draft)
	if err != nil {
		return err
	}

	res, err := m.collection.UpdateOne(
		ctx,
		query.EQField(models.DraftFieldsID.String(), id),
		utils.UpdateQuery(draftModel),
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return drafts.ErrDraftNotFound
	}

	return nil
}

// DeleteDraft delete draft
//
//	catchable errors:
//
//	1. ErrDraftIDNotValid
//
//	2. ErrDraftNotFound
func (m *MongoDraftRepository) DeleteDraft(ctx context.Context, draft aggregate.Draft) error {
	objId, err := m.ParseID(draft.ID)
	if err != nil {
		return err
	}

	res, err := m.collection.DeleteOne(ctx, query.EQField(models.DraftFieldsID.String(), objId))
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return drafts.ErrDraftNotFound
	}

	return nil
}
