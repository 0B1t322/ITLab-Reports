package reports

import (
	"context"

	"github.com/0B1t322/MongoBuilder/operators/query"
	"github.com/0B1t322/RepoGen/pkg/filter"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/models"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/utils"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReportsRepository struct {
	collection *mongo.Collection
	utils.IDChecker
}

func NewMongoReportsRepository(db *mongo.Database) *MongoReportsRepository {
	return &MongoReportsRepository{
		collection: db.Collection("reports"),
		IDChecker:  utils.NewIDChecker(reports.ErrIDIsNotValid),
	}
}

func NewMongoReportsRepositoryFrom(i *do.Injector) (*MongoReportsRepository, error) {
	db := do.MustInvoke[*mongo.Database](i)

	return NewMongoReportsRepository(db), nil
}

// GetReport return report by id
// catchable errors:
//
//	ErrIDIsNotValid
//	ErrReportNotFound
func (m *MongoReportsRepository) GetReport(
	ctx context.Context,
	id string,
) (aggregate.Report, error) {
	get, err := m.GetReports(
		ctx,
		reports.GetReportsQuery{
			Filter: reports.Query().
				Expression(
					reports.Expression().ID(id, filter.EQ),
				).
				Build(),
		},
	)
	if err != nil {
		return aggregate.Report{}, err
	}

	if len(get) == 0 {
		return aggregate.Report{}, reports.ErrReportNotFound
	}

	return get[0], nil
}

// CreateReport create report and return it
//
//	don't have catchable errors
func (m *MongoReportsRepository) CreateReport(ctx context.Context, report *aggregate.Report) error {
	report.ID = primitive.NewObjectID().Hex()

	reportModel, err := models.NewReportModel(*report)
	if err != nil {
		return err
	}

	_, err = m.collection.InsertOne(
		ctx,
		reportModel,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteReport delete report by id
//
//	catchable errors:
//		ErrIDIsNotValid
//		ErrReportNotFound
func (m *MongoReportsRepository) DeleteReport(ctx context.Context, report aggregate.Report) error {
	id, err := m.ParseID(report.ID)
	if err != nil {
		return err
	}

	dr, err := m.collection.DeleteOne(
		ctx,
		query.EQ(models.ReportFieldsID.String(), id),
	)
	if err != nil {
		return err
	}

	if dr.DeletedCount == 0 {
		return reports.ErrReportNotFound
	}

	return nil
}

// GetReports return reports acording to filters
//
//	don't have catchable errors
func (m *MongoReportsRepository) GetReports(
	ctx context.Context,
	query reports.GetReportsQuery,
) ([]aggregate.Report, error) {
	cur, err := m.collection.Aggregate(
		ctx,
		MarshallQuery(query),
	)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var reports []models.Report
	{
		if err := cur.All(ctx, &reports); err != nil {
			return nil, err
		}
	}

	return lo.Map(
		reports,
		func(r models.Report, _ int) aggregate.Report {
			return models.ReportFromModel(r)
		},
	), nil
}

// UpdateReport update reports by id and not nil optionals
//
//	catchable errors:
//		ErrIDIsNotValid
//		ErrReportNotFound
func (m *MongoReportsRepository) UpdateReport(
	ctx context.Context,
	report aggregate.Report,
) error {
	reportModel, err := models.NewReportModel(report)
	if err != nil {
		return reports.ErrIDIsNotValid
	}

	ur, err := m.collection.UpdateOne(
		ctx,
		query.EQField(
			models.ReportFieldsID.String(),
			reportModel.ID,
		),
		utils.UpdateQuery(
			reportModel,
		),
	)
	if err != nil {
		return err
	}

	if ur.MatchedCount == 0 {
		return reports.ErrReportNotFound
	}

	return nil
}

// CountByFilter count reports accroding to filter
//
//	don't have catchable errors
func (m *MongoReportsRepository) CountByFilter(
	ctx context.Context,
	query reports.GetReportsQuery,
) (int64, error) {
	type Count struct {
		Count int64 `bson:"count"`
	}

	var c Count
	{
		cur, err := m.collection.Aggregate(
			ctx,
			MarshallQueryForCount(query),
		)
		if err != nil {
			return 0, err
		}

		defer cur.Close(ctx)

		cur.Next(ctx)

		if err := cur.Decode(&c); err != nil {
			return 0, err
		}
	}

	return c.Count, nil
}
