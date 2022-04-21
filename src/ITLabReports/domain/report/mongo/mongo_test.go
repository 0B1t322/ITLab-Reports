package mongo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/config"
	domain "github.com/RTUITLab/ITLab-Reports/domain/report"
	assigneesEntity "github.com/RTUITLab/ITLab-Reports/entity/assignees"
	reportEntity "github.com/RTUITLab/ITLab-Reports/entity/report"

	"github.com/RTUITLab/ITLab-Reports/domain/report/mongo"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFunc_Filters(t *testing.T) {
	m := mongo.MongoRepository{}

	t.Run(
		"Filters",
		func(t *testing.T) {
			t.Run(
				"Name",
				func(t *testing.T) {
					t.Run(
						"EQ",
						func(t *testing.T) {
							actual := m.BuildFilters(
								&domain.GetReportsFilterFieldsWithOrAnd{
									GetReportsFilterFields: domain.GetReportsFilterFields{
										Name: &filter.FilterField[string]{
											Value:     "some",
											Operation: filter.EQ,
										},
									},
								},
							).Object()

							expected := bson.M{
								"name": "some",
							}

							require.Equal(
								t,
								expected,
								actual,
							)
						},
					)

					t.Run(
						"GT",
						func(t *testing.T) {
							actual := m.BuildFilters(
								&domain.GetReportsFilterFieldsWithOrAnd{
									GetReportsFilterFields: domain.GetReportsFilterFields{
										Name: &filter.FilterField[string]{
											Value:     "some",
											Operation: filter.GT,
										},
									},
								},
							).Object()

							expected := bson.M{
								"name": bson.M{
									"$gt": "some",
								},
							}

							require.Equal(
								t,
								expected,
								actual,
							)
						},
					)

				},
			)

			t.Run(
				"Date",
				func(t *testing.T) {
					t.Run(
						"EQ",
						func(t *testing.T) {
							date := time.Now().UTC()
							actual := m.BuildFilters(
								&domain.GetReportsFilterFieldsWithOrAnd{
									GetReportsFilterFields: domain.GetReportsFilterFields{
										Date: &filter.FilterField[string]{
											Value:     date.Format(time.RFC3339Nano),
											Operation: filter.EQ,
										},
									},
								},
							).Object()

							expected := bson.M{
								"date": date,
							}

							require.Equal(
								t,
								expected,
								actual,
							)
						},
					)
				},
			)
		},
	)

	t.Run(
		"With_Or_And",
		func(t *testing.T) {
			date := time.Now().UTC()

			actual := m.BuildFilters(
				&domain.GetReportsFilterFieldsWithOrAnd{
					And: []*domain.GetReportsFilterFieldsWithOrAnd{
						{
							Or: []*domain.GetReportsFilterFieldsWithOrAnd{
								{
									GetReportsFilterFields: domain.GetReportsFilterFields{
										Name: &filter.FilterField[string]{
											Value:     "some",
											Operation: filter.LIKE,
										},
									},
								},
								{
									GetReportsFilterFields: domain.GetReportsFilterFields{
										Implementer: &filter.FilterField[string]{
											Value:     "id_1",
											Operation: filter.EQ,
										},
									},
								},
								{
									GetReportsFilterFields: domain.GetReportsFilterFields{
										Implementer: &filter.FilterField[string]{
											Value:     "id_2",
											Operation: filter.EQ,
										},
									},
								},
							},
						},
						{
							GetReportsFilterFields: domain.GetReportsFilterFields{
								Date: &filter.FilterField[string]{
									Value:     date.Format(time.RFC3339Nano),
									Operation: filter.GTE,
								},
							},
						},
					},
				},
			).Object()

			expected := bson.M{
				"$and": bson.A{
					bson.M{
						"$or": bson.A{
							bson.M{
								"name": bson.M{
									"$regex":   "some",
									"$options": "i",
								},
							},
							bson.M{
								"assignees.implementer": "id_1",
							},
							bson.M{
								"assignees.implementer": "id_2",
							},
						},
					},
					bson.M{
						"date": bson.M{
							"$gte": date,
						},
					},
				},
			}

			require.Equal(
				t,
				expected,
				actual,
			)
		},
	)

}

func TestFunc_Sort(t *testing.T) {
	m := mongo.MongoRepository{}

	t.Run(
		"Name",
		func(t *testing.T) {
			actual := m.BuildSort(
				&domain.GetReportsSort{
					NameSort: *optional.NewOptional[ordertype.OrderType](ordertype.ASC),
				},
			)

			expected := bson.D{
				{"name", 1},
			}

			require.Equal(
				t,
				expected,
				actual,
			)
		},
	)

	t.Run(
		"Date",
		func(t *testing.T) {
			actual := m.BuildSort(
				&domain.GetReportsSort{
					DateSort: *optional.NewOptional[ordertype.OrderType](ordertype.DESC),
				},
			)

			expected := bson.D{
				{"date", -1},
			}

			require.Equal(
				t,
				expected,
				actual,
			)
		},
	)
}

func TestFunc_MongoRepository(t *testing.T) {
	cfg := config.GetConfigFrom("../../../.env")
	mongoReportRepo, err := mongo.New(
		context.Background(),
		cfg.MongoDB.TestURI,
	)
	require.NoError(t, err)

	var reportRepo domain.ReportRepository = mongoReportRepo

	t.Run(
		"CreateReport",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					model, err := report.NewReport(
						"some_report_name",
						"some_markdown_text",
						"reporter_id_1",
						"implementer_id_2",
					)
					require.NoError(t, err)

					created, err := reportRepo.CreateReport(
						context.Background(),
						model,
					)
					require.NoError(t, err)

					model.Report.ID = created.Report.ID

					require.Equal(
						t,
						model,
						created,
					)

					defer reportRepo.DeleteReport(
						context.Background(),
						created.Report.ID,
					)
				},
			)
		},
	)

	t .Run(
		"DeleteReport",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					model, err := report.NewReport(
						"some_report_nam_that_will_be_deleted",
						"some_markdown_text",
						"reporter_id_1",
						"implementer_id_2",
					)
					require.NoError(t, err)

					created, err := reportRepo.CreateReport(
						context.Background(),
						model,
					)
					require.NoError(t, err)

					err = reportRepo.DeleteReport(
						context.Background(),
						created.Report.ID,
					)
					require.NoError(t, err)

					_, err = reportRepo.GetReport(
						context.Background(),
						created.Report.ID,
					)
					require.ErrorIs(t, err, domain.ErrReportNotFound)
				},
			)

			t.Run(
				"NotFound",
				func(t *testing.T) {
					err = reportRepo.DeleteReport(
						context.Background(),
						primitive.NewObjectID().Hex(),
					)
					require.ErrorIs(t, err, domain.ErrReportNotFound)
				},
			)

			t.Run(
				"IdIsNotValid",
				func(t *testing.T) {
					err = reportRepo.DeleteReport(
						context.Background(),
						"some_id",
					)
					require.ErrorIs(t, err, domain.ErrIDIsNotValid)
				},
			)
		},
	)


	t.Run(
		"GetReport",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					model, err := report.NewReport(
						"some_report_nam_that_will_be_deleted",
						"some_markdown_text",
						"reporter_id_1",
						"implementer_id_2",
					)
					require.NoError(t, err)

					created, err := reportRepo.CreateReport(
						context.Background(),
						model,
					)
					require.NoError(t, err)

					defer reportRepo.DeleteReport(
						context.Background(),
						created.Report.ID,
					)

					get, err := reportRepo.GetReport(
						context.Background(),
						created.Report.ID,
					)
					require.NoError(t, err)

					require.Equal(
						t,
						created,
						get,
					)

				},
			)

			t.Run(
				"NotFound",
				func(t *testing.T) {
					_, err := reportRepo.GetReport(
						context.Background(),
						primitive.NewObjectID().Hex(),
					)
					require.ErrorIs(
						t,
						err,
						domain.ErrReportNotFound,
					)
				},
			)

			t.Run(
				"InvalidID",
				func(t *testing.T) {
					_, err := reportRepo.GetReport(
						context.Background(),
						"some_id",
					)
					require.ErrorIs(
						t,
						err,
						domain.ErrIDIsNotValid,
					)
				},
			)
		},
	)

	t.Run(
		"Update",
		func(t *testing.T) {

			t.Run(
				"Success",
				func(t *testing.T) {

					t.Run(
						"Empty",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_report_that_will_be_not_updated",
								"some_text",
								"reporter_id_1",
								"implementor_id_1",
							)
							require.NoError(t, err)

							created, err := reportRepo.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							defer reportRepo.DeleteReport(
								context.Background(),
								created.Report.ID,
							)

							updated, err := reportRepo.UpdateReport(
								context.Background(),
								created.Report.ID,
								domain.UpdateReportParams{
								},
							)
							require.NoError(t, err)

							require.Equal(
								t,
								created,
								updated,
							)
						},
					)

					t.Run(
						"OnlyName",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_report_that_will_be_updated",
								"some_text",
								"reporter_id_1",
								"implementor_id_1",
							)
							require.NoError(t, err)

							created, err := reportRepo.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							defer reportRepo.DeleteReport(
								context.Background(),
								created.Report.ID,
							)

							const newName = "updatedName"

							updated, err := reportRepo.UpdateReport(
								context.Background(),
								created.Report.ID,
								domain.UpdateReportParams{
									Name: *optional.NewOptional(newName),
								},
							)
							require.NoError(t, err)

							created.Report.Name = newName
							created.Report.Date = updated.Report.Date

							require.Equal(
								t,
								created,
								updated,
							)
						},
					)

					t.Run(
						"AllFields",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_report_that_will_be_updated",
								"some_text",
								"reporter_id_1",
								"implementor_id_1",
							)
							require.NoError(t, err)

							created, err := reportRepo.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							defer reportRepo.DeleteReport(
								context.Background(),
								created.Report.ID,
							)

							const (
								newName = "updated_name"
								newImplementer = "updated_implementer"
								newText = "updated_text"
							)

							updated, err := reportRepo.UpdateReport(
								context.Background(),
								created.Report.ID,
								domain.UpdateReportParams{
									Name: *optional.NewOptional(newName),
									Text: *optional.NewOptional(newText),
									Implementer: *optional.NewOptional(newImplementer),
								},
							)
							require.NoError(t, err)

							created.Report.Name = newName
							created.Report.Text = newText
							created.Assignees.Implementer = newImplementer
							created.Report.Date = updated.Report.Date

							require.Equal(
								t,
								created,
								updated,
							)
						},
					)
				},
			)

			t.Run(
				"Fail",
				func(t *testing.T) {
					t.Run(
						"NotFound",
						func(t *testing.T) {
							_, err := reportRepo.UpdateReport(
								context.Background(),
								primitive.NewObjectID().Hex(),
								domain.UpdateReportParams{},
							)
							require.ErrorIs(t, err, domain.ErrReportNotFound)
						},
					)

					t.Run(
						"InvalidID",
						func(t *testing.T) {
							_, err := reportRepo.UpdateReport(
								context.Background(),
								"some_id",
								domain.UpdateReportParams{},
							)
							require.ErrorIs(t, err, domain.ErrIDIsNotValid)
						},
					)
				},
			)
		},
	)

	t.Run(
		"CountDocuments",
		func(t *testing.T) {
			oldCount, err := reportRepo.CountByFilter(
				context.Background(),
				&domain.GetReportsFilterFieldsWithOrAnd{},
			)
			require.NoError(t, err)

			model, err := report.NewReport(
				"some_report_that_will_be_not_updated",
				"some_text",
				"reporter_id_1",
				"implementor_id_1",
			)
			require.NoError(t, err)

			created, err := reportRepo.CreateReport(
				context.Background(),
				model,
			)
			require.NoError(t, err)

			defer reportRepo.DeleteReport(
				context.Background(),
				created.Report.ID,
			)

			newCount, err := reportRepo.CountByFilter(
				context.Background(),
				&domain.GetReportsFilterFieldsWithOrAnd{},
			)
			require.NoError(t, err)

			require.Greater(t, newCount, oldCount)
		},
	)

	t.Run(
		"GetReports",
		func(t *testing.T) {
			const Day = time.Hour * 24

			day := func(days int) time.Time {

				result := time.Now()
				for i := 0; i < days; i++ {
					result = result.Add(Day * -1)
				}

				return result
			}

			var createdReports []*report.Report

			for i := 0; i < 10; i++ {
				created, err := reportRepo.CreateReport(
					context.Background(),
					&report.Report{
						Report: &reportEntity.Report{
							Name: fmt.Sprintf("report_%v", i),
							Text: fmt.Sprintf("text_%v", i),
							Date: day(i).UTC().Round(time.Millisecond),
						},
						Assignees: &assigneesEntity.Assignees{
							Implementer: fmt.Sprintf("implementer_%v", i),
							Reporter: fmt.Sprintf("reporter_%v", i),
						},
					},
				)
				require.NoError(t, err)

				createdReports = append(createdReports, created)

				defer reportRepo.DeleteReport(
					context.Background(),
					created.Report.ID,
				)
			}

			t.Run(
				"All",
				func(t *testing.T) {
					getReports, err := reportRepo.GetReports(
						context.Background(),
						&domain.GetReportsParams{
							Filter: &domain.GetReportsFilter{},
						},
					)
					require.NoError(t, err)

					for _, createdReport := range createdReports {
						require.Contains(t, getReports, createdReport)
					}
				},
			)

			t.Run(
				"Names_Contains_1_or_9",
				func(t *testing.T) {
					getReports, err := reportRepo.GetReports(
						context.Background(),
						&domain.GetReportsParams{
							Filter: &domain.GetReportsFilter{
								GetReportsFilterFieldsWithOrAnd: domain.GetReportsFilterFieldsWithOrAnd{
									Or: []*domain.GetReportsFilterFieldsWithOrAnd{
										{
											GetReportsFilterFields: domain.GetReportsFilterFields{
												Name: &filter.FilterField[string]{
													Value: "1",
													Operation: filter.LIKE,
												},
											},
										},
										{
											GetReportsFilterFields: domain.GetReportsFilterFields{
												Name: &filter.FilterField[string]{
													Value: "9",
													Operation: filter.LIKE,
												},
											},
										},
									},
								},
							},
						},
					)
					require.NoError(t, err)

					for _, getReport := range getReports {
						require.Condition(
							t,
							func() (success bool) {
								return getReport.Report.Name == "report_1" || getReport.Report.Name == "report_9"
							},
						)
					}
				},
			)

			t.Run(
				"Date_betwwen_1_day_ago_and_3_days_ago",
				func(t *testing.T) {
					date := time.Now().UTC()
					getReports, err := reportRepo.GetReports(
						context.Background(),
						&domain.GetReportsParams{
							Filter: &domain.GetReportsFilter{
								GetReportsFilterFieldsWithOrAnd: domain.GetReportsFilterFieldsWithOrAnd{
									And: []*domain.GetReportsFilterFieldsWithOrAnd{
										{
											GetReportsFilterFields: domain.GetReportsFilterFields{
												Date: &filter.FilterField[string]{
													Value: date.Add(-1*Day).Format(time.RFC3339Nano),
													Operation: filter.LTE,
												},
											},
										},
										{
											GetReportsFilterFields: domain.GetReportsFilterFields{
												Date: &filter.FilterField[string]{
													Value: date.Add(-3*Day).Format(time.RFC3339Nano),
													Operation: filter.GTE,
												},
											},
										},
									},
								},
							},
						},
					)
					require.NoError(t, err)

					isInBorder := func(report *report.Report) (success bool) {
						lte := report.Report.Date.Equal(date.Add(-1*Day)) || report.Report.Date.Before(date.Add(-1*Day))

						gte := report.Report.Date.Equal(date.Add(-3*Day)) || report.Report.Date.After(date.Add(-3*Day))
						
						return lte && gte
					}

					for _, report := range getReports {
						require.Condition(
							t,
							func() (success bool) {
								return isInBorder(report)
							},
						)
					}

					// Check that in crated only those reports with this date borders
					require.Condition(
						t,
						func() (success bool) {
							reportThatNotInBorderCount := func() int {
								count := 0
								for _, report := range createdReports {
									if !isInBorder(report) {
										count++
									}
								}
								return count
							}()
							
							return reportThatNotInBorderCount + len(getReports) >= len(createdReports)

						},
					)
				},
			)
		},
	)
}