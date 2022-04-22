package mongo

import (
	"context"
	"time"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	domain "github.com/RTUITLab/ITLab-Reports/domain/report"
	modelAssignes "github.com/RTUITLab/ITLab-Reports/entity/assignees"
	modelReport "github.com/RTUITLab/ITLab-Reports/entity/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/mongobuildertofilter"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/ordertypetomongo"
	builder "github.com/RTUITLab/ITLab-Reports/pkg/dialect/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Report = modelReport.Report
type Assignees = modelAssignes.Assignees

type MongoRepository struct {
	db *mongo.Database

	reports *mongo.Collection

	collectionName string
}

type MongoRepositoryOptions func (m *MongoRepository)

func WithCollectionName(collection string) MongoRepositoryOptions {
	return func(m *MongoRepository) {
		m.collectionName = collection
	}
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

	conn, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(connString),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := connstring.Parse(connString)
	if err != nil {
		return nil, err
	}

	r.db = conn.Database(connStr.Database)

	if r.collectionName == "" {
		r.collectionName = "reports"
	}
	r.reports = r.db.Collection(r.collectionName)

	return r, nil
}



type mongoReportModel struct {
	ID primitive.ObjectID `bson:"_id"`

	Name string `bson:"name"`

	Date time.Time `bson:"date"`

	Text string `bson:"text"`

	Assignees assignesModel `bson:"assignees"`
}

type assignesModel struct {
	Reporter    string `bson:"reporter"`
	Implementer string `bson:"implementer"`
}

func (m mongoReportModel) ToAgragate() *report.Report {
	return &report.Report{
		Report: &Report{
			ID: m.ID.Hex(),
			Name: m.Name,
			Date: m.Date,
			Text: m.Text,
		},
		Assignees: &Assignees{
			Reporter: m.Assignees.Reporter,
			Implementer: m.Assignees.Implementer,
		},
	}
}

func fromAgragate(report *report.Report) mongoReportModel {
	return mongoReportModel{
		Name: report.Report.Name,
		Date: report.Report.Date,
		Text: report.Report.Text,
		Assignees: assignesModel{
			Reporter: report.Assignees.Reporter,
			Implementer: report.Assignees.Implementer,
		},
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

	var get mongoReportModel
	{
		sr := m.reports.FindOne(
			ctx,
			builder.P().
				EQ("_id", objId),
		)

		if err := sr.Err(); err == mongo.ErrNoDocuments {
			return nil, domain.ErrReportNotFound
		} else if err != nil{
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
		builder.P().EQ("_id", objId),
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
) *builder.Predicate {
	b := builder.P()

	if filter == nil {
		return b
	}

	if filter.Name != nil {
		filter.Name.BuildTo(
			mongobuildertofilter.New[string](
				b,
				"name",
			),
		)
	}

	if filter.ReportID != nil {
		filter.ReportID.BuildTo(
			mongobuildertofilter.New[string](
				b,
				"_id",
			),
		)
	}

	if filter.Date != nil {
		filter.Date.BuildTo(
			mongobuildertofilter.New(
				b,
				"date",
				mongobuildertofilter.WithFieldFormatter(
					mongobuildertofilter.FieldToTime[string](),
				),
			),
		)
	}

	if filter.Implementer != nil {
		filter.Implementer.BuildTo(
			mongobuildertofilter.New[string](
				b,
				"assignees.implementer",
			),
		)
	}

	if filter.Reporter != nil {
		filter.Reporter.BuildTo(
			mongobuildertofilter.New[string](
				b,
				"assignees.reporter",
			),
		)
	}

	if filter.Or != nil {
		b.Or(
			func() (preds []*builder.Predicate) {
				for _, f := range filter.Or {
					preds = append(preds, m.BuildFilters(f))
				}

				return preds
			}()...,
		)
	}

	if filter.And != nil {
		b.And(
			func() (preds []*builder.Predicate) {
				for _, f := range filter.And {
					preds = append(preds, m.BuildFilters(f))
				}

				return preds
			}()...,
		)
	}

	return b
}


func (m *MongoRepository) BuildFindOptions(
	params	*domain.GetReportsParams,
) *options.FindOptions {
	opt := options.Find().
				SetSort(m.BuildSort(&params.Filter.GetReportsSort))

	if params.Limit.HasValue() {
		opt.SetLimit(params.Limit.MustGetValue())
	}

	if params.Offset.HasValue() {
		opt.SetLimit(params.Offset.MustGetValue())
	}

	return opt
}

func (m *MongoRepository) BuildSort(
	sort *domain.GetReportsSort,
) any {
	mgSort := bson.D{}

	if sort == nil {
		return mgSort
	}
	
	if sort.DateSort.HasValue() {
		if order := ordertypetomongo.ToMongoOrderType(sort.DateSort.MustGetValue()); order != 0 {
			mgSort = append(mgSort, bson.E{Key: "date", Value: order})
		}
	}

	if sort.NameSort.HasValue() {
		if order := ordertypetomongo.ToMongoOrderType(sort.NameSort.MustGetValue()); order != 0 {
			mgSort = append(mgSort, bson.E{Key: "name", Value: order})
		}
	}

	return mgSort
}

func (m *MongoRepository) GetReports(
	ctx context.Context,
	params *domain.GetReportsParams,
) ([]*report.Report, error) {
	var mgReports []mongoReportModel
	{
		cur, err := m.reports.Find(
			ctx,
			m.BuildFilters(&params.Filter.GetReportsFilterFieldsWithOrAnd),
			m.BuildFindOptions(params),
		)
		if err != nil {
			return nil, err
		}

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

func (m *MongoRepository) BuildUpdateFields(params domain.UpdateReportParams) *builder.SetBuidler {
	return builder.Set().
			SetFieldIf("name", params.Name.MustGetValue(), func() bool {return params.Name.HasValue()}).
			SetFieldIf("text", params.Text.MustGetValue(), func() bool {return params.Text.HasValue()}).
			SetFieldIf("assignees.implementer", params.Implementer.MustGetValue(), func() bool {return params.Implementer.HasValue()}).
			SetFieldIf("date", time.Now().UTC(), func() bool {return params.Name.HasValue() || params.Implementer.HasValue() || params.Text.HasValue()})
}

func (m *MongoRepository) UpdateReport(
	ctx	context.Context,
	id	string,
	params domain.UpdateReportParams,
) (*report.Report, error) {
	mgId, err := m.getObjectID(id)
	if err != nil {
		return nil, err
	}

	ur, err := m.reports.UpdateOne(
		ctx,
		builder.EQ("_id", mgId),
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
