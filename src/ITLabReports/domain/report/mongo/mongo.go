package mongo

import (
	"context"
	"time"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	domain "github.com/RTUITLab/ITLab-Reports/domain/report"
	modelAssignes "github.com/RTUITLab/ITLab-Reports/entity/assignees"
	modelReport "github.com/RTUITLab/ITLab-Reports/entity/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/mongobuildertofilter"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/ordertypetosortorder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/0B1t322/MongoBuilder/operators/query"
	"github.com/0B1t322/MongoBuilder/operators/sort"
	"github.com/0B1t322/MongoBuilder/operators/update"
	"github.com/0B1t322/MongoBuilder/utils"
)

type Report = modelReport.Report
type Assignees = modelAssignes.Assignees

type MongoRepository struct {
	client *mongo.Client

	db *mongo.Database

	reports *mongo.Collection

	collectionName string
}

type MongoRepositoryOptions func(m *MongoRepository)

func WithCollectionName(collection string) MongoRepositoryOptions {
	return func(m *MongoRepository) {
		m.collectionName = collection
	}
}

func WithClient(client *mongo.Client) MongoRepositoryOptions {
	return func(m *MongoRepository) {
		m.client = client
	}
}

func (m *MongoRepository) createClient(
	ctx context.Context,
	connString string,
) error {
	conn, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(connString),
	)
	if err != nil {
		return err
	}

	m.client = conn
	return nil
}

func New(
	ctx context.Context,
	connString string,
	opts ...MongoRepositoryOptions,
) (*MongoRepository, error) {
	r := &MongoRepository{}

	for _, opt := range opts {
		opt(r)
	}

	if r.client == nil {
		if err := r.createClient(ctx, connString); err != nil {
			return nil, err
		}
	}

	connStr, err := connstring.Parse(connString)
	if err != nil {
		return nil, err
	}

	r.db = r.client.Database(connStr.Database)

	if r.collectionName == "" {
		r.collectionName = "reports"
	}
	r.reports = r.db.Collection(r.collectionName)
	return r, nil
}

type MongoReportModel struct {
	ID        primitive.ObjectID      `bson:"_id"`
	Name      string                  `bson:"name"`
	Date      time.Time               `bson:"date"`
	Text      string                  `bson:"text"`
	Assignees MongoAssignesModel      `bson:"assignees"`
	State     modelReport.ReportState `bson:"state,omitempty"`
}

type MongoAssignesModel struct {
	Reporter    string `bson:"reporter"`
	Implementer string `bson:"implementer"`
}

func (m MongoReportModel) ToAgragate() *report.Report {
	return &report.Report{
		Report: &Report{
			ID:    m.ID.Hex(),
			Name:  m.Name,
			Date:  m.Date,
			Text:  m.Text,
			State: m.State,
		},
		Assignees: &Assignees{
			Reporter:    m.Assignees.Reporter,
			Implementer: m.Assignees.Implementer,
		},
	}
}

func fromAgragate(report *report.Report) MongoReportModel {
	return MongoReportModel{
		Name: report.Report.Name,
		Date: report.Report.Date,
		Text: report.Report.Text,
		Assignees: MongoAssignesModel{
			Reporter:    report.Assignees.Reporter,
			Implementer: report.Assignees.Implementer,
		},
		State: report.Report.State,
	}
}

func (m *MongoRepository) getObjectID(id string) (primitive.ObjectID, error) {
	if objId, err := primitive.ObjectIDFromHex(id); err != nil {
		return primitive.NilObjectID, domain.ErrIDIsNotValid
	} else {
		return objId, nil
	}
}

func (m *MongoRepository) GetReport(
	ctx context.Context,
	id string,
) (*report.Report, error) {
	objId, err := m.getObjectID(id)
	if err != nil {
		return nil, err
	}

	var get MongoReportModel
	{
		sr := m.reports.FindOne(
			ctx,
			query.EQField("_id", objId),
		)

		if err := sr.Err(); err == mongo.ErrNoDocuments {
			return nil, domain.ErrReportNotFound
		} else if err != nil {
			return nil, err
		}

		if err := sr.Decode(&get); err != nil {
			return nil, err
		}
	}

	return get.ToAgragate(), nil

}

func (m *MongoRepository) CreateReport(
	ctx context.Context,
	report *report.Report,
) (*report.Report, error) {
	model := fromAgragate(report)
	model.ID = primitive.NewObjectID()
	ir, err := m.reports.InsertOne(
		ctx,
		model,
	)
	if err != nil {
		return nil, err
	}

	model.ID = ir.InsertedID.(primitive.ObjectID)

	return model.ToAgragate(), nil
}

func (m *MongoRepository) DeleteReport(
	ctx context.Context,
	id string,
) error {
	objId, err := m.getObjectID(id)
	if err != nil {
		return err
	}

	result, err := m.reports.DeleteOne(
		ctx,
		query.EQField("_id", objId),
	)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrReportNotFound
	}

	return nil
}

func (m *MongoRepository) BuildFilters(
	filter *domain.GetReportsFilterFieldsWithOrAnd,
) bson.M {
	b := bson.M{}

	if filter == nil {
		return b
	}

	if filter.Name != nil {
		filter.Name.BuildTo(
			mongobuildertofilter.NewBuilderQueryAdapter[string](
				&b,
				"name",
			),
		)
	}

	if filter.ReportID != nil {
		filter.ReportID.BuildTo(
			mongobuildertofilter.NewBuilderQueryAdapter(
				&b,
				"_id",
				mongobuildertofilter.
					NewBuilderQueryAdapterOptions[string]().
					WithFieldFormatter(
						mongobuildertofilter.StringIdMarshaller(),
					),
			),
		)
	}

	if filter.ReportsId != nil {
		filter.ReportsId.BuildTo(
			mongobuildertofilter.NewBuilderQueryAdapter(
				&b,
				"_id",
				mongobuildertofilter.
					NewBuilderQueryAdapterOptions[[]string]().
					WithFieldFormatter(
						mongobuildertofilter.SliceIdMarshaller(),
					),
			),
		)
	}

	if filter.Date != nil {
		filter.Date.BuildTo(
			mongobuildertofilter.NewBuilderQueryAdapter(
				&b,
				"date",
				mongobuildertofilter.
					NewBuilderQueryAdapterOptions[string]().
					WithFieldFormatter(
						mongobuildertofilter.FieldToTimeMarshaller[string](),
					),
			),
		)
	}

	if filter.Implementer != nil {
		filter.Implementer.BuildTo(
			mongobuildertofilter.NewBuilderQueryAdapter[string](
				&b,
				"assignees.implementer",
			),
		)
	}

	if filter.Reporter != nil {
		filter.Reporter.BuildTo(
			mongobuildertofilter.NewBuilderQueryAdapter[string](
				&b,
				"assignees.reporter",
			),
		)
	}

	if filter.Or != nil {
		b = utils.MergeBsonM(
			b,
			query.Or(
				func() (preds []bson.M) {
					for _, f := range filter.Or {
						preds = append(preds, m.BuildFilters(f))
					}

					return preds
				}()...,
			),
		)
	}

	if filter.And != nil {
		b = utils.MergeBsonM(
			b,
			query.And(
				func() (preds []bson.M) {
					for _, f := range filter.And {
						preds = append(preds, m.BuildFilters(f))
					}

					return preds
				}()...,
			),
		)
	}

	return b
}

func (m *MongoRepository) BuildFindOptions(
	params *domain.GetReportsParams,
) *options.FindOptions {
	opt := options.Find().
		SetSort(m.BuildSort(params.Filter.SortParams))

	if params.Limit.IsPresent() {
		opt.SetLimit(params.Limit.MustGet())
	}

	if params.Offset.IsPresent() {
		opt.SetSkip(params.Offset.MustGet())
	}

	return opt
}

func (m *MongoRepository) BuildSort(
	sortParams []domain.GetReportsSort,
) any {
	if len(sortParams) == 0 {
		return nil
	}

	sortArgs := []sort.SortArger{}
	{
		for _, s := range sortParams {
			if s.NameSort.IsPresent() {
				sortArgs = append(
					sortArgs,
					sort.SortArg(
						"name",
						ordertypetosortorder.ToSortOrder(s.NameSort.MustGet()),
					),
				)
			}

			if s.DateSort.IsPresent() {
				sortArgs = append(
					sortArgs,
					sort.SortArg(
						"date",
						ordertypetosortorder.ToSortOrder(s.DateSort.MustGet()),
					),
				)
			}
		}
	}
	return sort.Sort(sortArgs...)["$sort"]
}

func (m *MongoRepository) GetReports(
	ctx context.Context,
	params *domain.GetReportsParams,
) ([]*report.Report, error) {
	var mgReports []MongoReportModel
	{
		cur, err := m.reports.Find(
			ctx,
			m.BuildFilters(&params.Filter.GetReportsFilterFieldsWithOrAnd),
			m.BuildFindOptions(params),
		)
		if err != nil {
			return nil, err
		}
		defer cur.Close(ctx)

		if err := cur.All(ctx, &mgReports); err != nil {
			return nil, err
		}
	}

	var reports []*report.Report
	{
		for _, mgReport := range mgReports {
			reports = append(reports, mgReport.ToAgragate())
		}
	}

	return reports, nil
}

func (m *MongoRepository) BuildUpdateFields(params domain.UpdateReportParams) bson.M {
	var updateArgs []update.SetArger
	{
		if params.Name.IsPresent() {
			updateArgs = append(
				updateArgs,
				update.SetArg("name", params.Name.MustGet()),
			)
		}

		if params.Text.IsPresent() {
			updateArgs = append(
				updateArgs,
				update.SetArg("text", params.Text.MustGet()),
			)
		}

		if params.Implementer.IsPresent() {
			updateArgs = append(
				updateArgs,
				update.SetArg("assignees.implementer", params.Implementer.MustGet()),
			)
		}
	}
	if len(updateArgs) == 0 {
		return bson.M{"$set": bson.M{}}
	}
	return update.Set(
		updateArgs...,
	)
}

func (m *MongoRepository) UpdateReport(
	ctx context.Context,
	id string,
	params domain.UpdateReportParams,
) (*report.Report, error) {
	mgId, err := m.getObjectID(id)
	if err != nil {
		return nil, err
	}

	ur, err := m.reports.UpdateOne(
		ctx,
		query.EQField("_id", mgId),
		m.BuildUpdateFields(params),
	)
	if err != nil {
		return nil, err
	}

	if ur.MatchedCount == 0 {
		return nil, domain.ErrReportNotFound
	}

	return m.GetReport(ctx, id)
}

func (m *MongoRepository) CountByFilter(
	ctx context.Context,
	params *domain.GetReportsFilterFieldsWithOrAnd,
) (int64, error) {
	count, err := m.reports.CountDocuments(
		ctx,
		m.BuildFilters(params),
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}
