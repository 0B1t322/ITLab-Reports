package migrations

import (
	"context"

	stage "github.com/0B1t322/MongoBuilder/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/field"
	"github.com/RTUITLab/ITLab-Reports/internal/services/salary"
	"github.com/RTUITLab/ITLab-Reports/internal/services/token"
	"github.com/samber/mo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddStateMigration struct {
	SalaryService salary.SalaryService
	tokenService  token.TokenService
}

func NewAddStateMigration(
	salaryService salary.SalaryService,
	tokenService token.TokenService,
) *AddStateMigration {
	return &AddStateMigration{
		SalaryService: salaryService,
		tokenService:  tokenService,
	}
}

func (a *AddStateMigration) Version() uint64 {
	return 5
}

func (a *AddStateMigration) Description() string {
	return "add reports state field"
}

func (a *AddStateMigration) Up(db *mongo.Database) error {
	token, err := a.tokenService.RequestToken(context.Background())
	if err != nil {
		return err
	}

	ids, err := a.SalaryService.GetApprovedReportsIds(
		context.Background(),
		token,
		mo.None[string](),
	)
	if err != nil {
		return err
	}

	var objIds bson.A
	{
		for _, id := range ids {
			objId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				continue
			}
			objIds = append(objIds, objId)
		}
	}

	pipeline := []bson.M{
		stage.AddFields(
			stage.AddFieldArg().
				AddField(
					"state",
					aggregation.Switch(
						aggregation.SwitchArg().
							AddCase(
								aggregation.In(
									field.Field("_id"),
									objIds,
								),
								"paid",
							).
							Default(
								"created",
							),
					),
				),
		),
		stage.Merge(
			stage.MergeArg().
				IntoCollection("reports").
				WhenMatched(stage.WhenMatchedArg().Merge()),
		),
	}

	_, err = db.Collection("reports").Aggregate(
		context.Background(),
		pipeline,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *AddStateMigration) Down(db *mongo.Database) error {
	pipeline := []bson.M{
		stage.Unset(
			stage.UnsetArg().Unset("state"),
		),
		stage.Merge(
			stage.MergeArg().
				IntoCollection("reports").
				WhenMatched(stage.WhenMatchedArg().Replace()),
		),
	}

	_, err := db.Collection("reports").Aggregate(
		context.Background(),
		pipeline,
	)
	if err != nil {
		return err
	}

	return nil
}
