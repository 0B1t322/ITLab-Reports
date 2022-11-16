package reports

import (
	"strings"
	"time"

	stage "github.com/0B1t322/MongoBuilder/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/query"
	"github.com/0B1t322/MongoBuilder/operators/sort"
	"github.com/0B1t322/MongoBuilder/utils"

	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/adapters/filterbuilder"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/adapters/sortorder"
	"github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	Filter reports.FilterQuery
	Sort   []reports.SortFields
)

func (f Filter) BuildID(to *bson.M) {
	if id, ok := f.Expression.ID.Get(); ok {
		id.BuildTo(
			filterbuilder.New(
				to,
				models.ReportFieldsID.String(),
				filterbuilder.WithFieldFormatter(
					filterbuilder.StringIdMarshaller(),
				),
			),
		)
	}
}

func (f Filter) BuildIDs(to *bson.M) {
	if ids, ok := f.Expression.IDs.Get(); ok {
		ids.BuildTo(
			filterbuilder.New(
				to,
				models.ReportFieldsID.String(),
				filterbuilder.WithFieldFormatter(
					filterbuilder.SliceIdMarshaller(),
				),
			),
		)
	}
}

func (f Filter) BuildName(to *bson.M) {
	if name, ok := f.Expression.Name.Get(); ok {
		name.BuildTo(
			filterbuilder.New[string](
				to,
				models.ReportFieldsName.String(),
			),
		)
	}
}

func (f Filter) BuildDate(to *bson.M) {
	if date, ok := f.Expression.Date.Get(); ok {
		date.BuildTo(
			filterbuilder.New[time.Time](
				to,
				models.ReportFieldsDate.String(),
			),
		)
	}
}

func (f Filter) BuildImplementer(to *bson.M) {
	if implementer, ok := f.Expression.Implementer.Get(); ok {
		implementer.BuildTo(
			filterbuilder.New[string](
				to,
				strings.Join(
					[]string{
						models.ReportFieldsAssignees.String(),
						models.AssigneesFieldsImplementer.String(),
					},
					".",
				),
			),
		)
	}
}

func (f Filter) BuildReporter(to *bson.M) {
	if reporter, ok := f.Expression.Reporter.Get(); ok {
		reporter.BuildTo(
			filterbuilder.New[string](
				to,
				strings.Join(
					[]string{
						models.ReportFieldsAssignees.String(),
						models.AssigneesFieldsReporter.String(),
					},
					".",
				),
			),
		)
	}
}

func (f Filter) BuildState(to *bson.M) {
	if state, ok := f.Expression.State.Get(); ok {
		*to = utils.MergeBsonM(
			*to,
			query.EQ(models.ReportFieldsState.String(), state),
		)
	}
}

func (f Filter) BuildExpression(to *bson.M) {
	f.BuildID(to)
	f.BuildIDs(to)
	f.BuildName(to)
	f.BuildDate(to)
	f.BuildReporter(to)
	f.BuildImplementer(to)
	f.BuildState(to)
}

func (m Filter) BuildOr(to *bson.M) {
	if len(m.Or) > 0 {
		*to = utils.MergeBsonM(
			*to,
			aggregation.Or(
				func() (preds []interface{}) {
					for _, filter := range m.Or {
						preds = append(preds, Filter(filter).BuildAll())
					}
					return preds
				}()...,
			),
		)
	}
}

func (m Filter) BuildAnd(to *bson.M) {
	if len(m.And) > 0 {
		*to = utils.MergeBsonM(
			*to,
			aggregation.And(
				func() (preds []interface{}) {
					for _, filter := range m.And {
						preds = append(preds, Filter(filter).BuildAll())
					}
					return preds
				}()...,
			),
		)
	}
}

func (m Filter) BuildAll() bson.M {
	b := bson.M{}
	{
		m.BuildExpression(&b)
		m.BuildOr(&b)
		m.BuildAnd(&b)
	}
	return b
}

func (m Filter) Build() bson.M {
	b := m.BuildAll()

	if len(b) > 0 {
		return stage.Match(
			b,
		)
	}
	return b
}

func (s Sort) Build() bson.M {
	var args []sort.SortArger
	{
		for _, s := range s {
			s.Name.ForEach(
				func(order sortorder.SortOrder) {
					args = append(args, sort.SortArg(
						models.ReportFieldsName.String(),
						sortorder.ToSortOrder(order),
					))
				},
			)

			s.Date.ForEach(
				func(order sortorder.SortOrder) {
					args = append(args, sort.SortArg(
						models.ReportFieldsName.String(),
						sortorder.ToSortOrder(order),
					))
				},
			)
		}
	}

	if len(args) > 0 {
		return stage.Sort(
			args...,
		)
	}

	return bson.M{}
}

func MarshallQuery(q reports.GetReportsQuery) (pipeline []bson.M) {
	if filterStage := Filter(q.Filter).Build(); len(filterStage) > 0 {
		pipeline = append(pipeline, filterStage)
	}

	q.Offset.ForEach(
		func(value int64) {
			pipeline = append(pipeline, stage.Skip(value))
		},
	)

	q.Limit.ForEach(
		func(value int64) {
			pipeline = append(pipeline, stage.Limit(int(value)))
		},
	)

	if sortStage := Sort(q.Sort).Build(); len(sortStage) > 0 {
		pipeline = append(pipeline, sortStage)
	}

	return
}

func MarshallQueryForCount(q reports.GetReportsQuery) []bson.M {
	var pipeline []bson.M
	{
		if filterStage := Filter(q.Filter).Build(); len(filterStage) > 0 {
			pipeline = append(pipeline, filterStage)
		}

		pipeline = append(pipeline, stage.Count("count"))
	}
	return pipeline
}
