package reports_test

import (
	"context"
	"testing"
	"time"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/samber/mo"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/config"
	entAssignees "github.com/RTUITLab/ITLab-Reports/entity/assignees"
	entReport "github.com/RTUITLab/ITLab-Reports/entity/report"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/service/reports/reportservice"

	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/stretchr/testify/require"
)

func Service_Tests(t *testing.T, service reports.Service) {
	t.Run(
		"Create",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					model, err := report.NewReport(
						"some_name",
						"some_text",
						"some_reporter",
						"some_implementer",
					)
					require.NoError(t, err)

					created, err := service.CreateReport(
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
				},
			)

			t.Run(
				"ValidationError",
				func(t *testing.T) {
					t.Run(
						"EmptyName",
						func(t *testing.T) {
							_, err := service.CreateReport(
								context.Background(),
								&report.Report{
									Report: &entReport.Report{
										Name: "",
										Date: time.Now(),
										Text: "some",
									},
									Assignees: &entAssignees.Assignees{
										Reporter:    "some_reporter",
										Implementer: "some_implementer",
									},
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return errors.Is(err, reports.ErrValidationError) && errors.Is(err, report.ErrNameEmpty)
								},
							)

						},
					)

					t.Run(
						"EmptyText",
						func(t *testing.T) {
							_, err := service.CreateReport(
								context.Background(),
								&report.Report{
									Report: &entReport.Report{
										Name: "some_name",
										Date: time.Now(),
										Text: "",
									},
									Assignees: &entAssignees.Assignees{
										Reporter:    "some_reporter",
										Implementer: "some_implementer",
									},
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return errors.Is(err, reports.ErrValidationError) && errors.Is(err, report.ErrTextEmpty)
								},
							)
						},
					)

					t.Run(
						"EmptyImplementer",
						func(t *testing.T) {
							_, err := service.CreateReport(
								context.Background(),
								&report.Report{
									Report: &entReport.Report{
										Name: "some_name",
										Date: time.Now(),
										Text: "some",
									},
									Assignees: &entAssignees.Assignees{
										Reporter:    "some_reporter",
										Implementer: "",
									},
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return errors.Is(err, reports.ErrValidationError) && errors.Is(err, report.ErrImplementorEmpty)
								},
							)
						},
					)

					t.Run(
						"EmptyRepoter",
						func(t *testing.T) {
							_, err := service.CreateReport(
								context.Background(),
								&report.Report{
									Report: &entReport.Report{
										Name: "some_name",
										Date: time.Now(),
										Text: "some",
									},
									Assignees: &entAssignees.Assignees{
										Reporter:    "",
										Implementer: "some_implementer",
									},
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return errors.Is(err, reports.ErrValidationError) && errors.Is(err, report.ErrReporterEmpty)
								},
							)
						},
					)
				},
			)
		},
	)

	t.Run(
		"DeleteReport",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					model, err := report.NewReport(
						"some_name",
						"some_text",
						"some_reporter",
						"some_implementer",
					)
					require.NoError(t, err)

					created, err := service.CreateReport(
						context.Background(),
						model,
					)
					require.NoError(t, err)

					err = service.DeleteReport(
						context.Background(),
						created.GetID(),
					)
					require.NoError(t, err)

					_, err = service.GetReport(
						context.Background(),
						created.GetID(),
					)
					require.ErrorIs(t, err, reports.ErrReportNotFound)
				},
			)

			t.Run(
				"Fail",
				func(t *testing.T) {
					err := service.DeleteReport(
						context.Background(),
						"some_id",
					)
					require.Condition(
						t,
						func() (success bool) {
							return err == reports.ErrReportIDNotValid || err == reports.ErrReportNotFound
						},
					)
				},
			)
		},
	)

	t.Run(
		"UpdateReport",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"UpdateNothing",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_name",
								"some_text",
								"some_reporter",
								"some_implementer",
							)
							require.NoError(t, err)

							created, err := service.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							updated, err := service.UpdateReport(
								context.Background(),
								created.GetID(),
								reportdomain.UpdateReportParams{},
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
						"UpdateOnlyName",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_name",
								"some_text",
								"some_reporter",
								"some_implementer",
							)
							require.NoError(t, err)

							created, err := service.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							updated, err := service.UpdateReport(
								context.Background(),
								created.GetID(),
								reportdomain.UpdateReportParams{
									Name: mo.Some("new_name"),
								},
							)
							require.NoError(t, err)

							created.Report.Date = updated.GetDate()
							created.Report.Name = "new_name"

							require.Equal(
								t,
								created,
								updated,
							)
						},
					)

					t.Run(
						"UpdateOnlyText",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_name",
								"some_text",
								"some_reporter",
								"some_implementer",
							)
							require.NoError(t, err)

							created, err := service.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							updated, err := service.UpdateReport(
								context.Background(),
								created.GetID(),
								reportdomain.UpdateReportParams{
									Text: mo.Some("updated_text"),
								},
							)
							require.NoError(t, err)

							created.Report.Text = "updated_text"
							created.Report.Date = updated.GetDate()

							require.Equal(
								t,
								created,
								updated,
							)
						},
					)

					t.Run(
						"UpdateOnlyImplementer",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_name",
								"some_text",
								"some_reporter",
								"some_implementer",
							)
							require.NoError(t, err)

							created, err := service.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							updated, err := service.UpdateReport(
								context.Background(),
								created.GetID(),
								reportdomain.UpdateReportParams{
									Implementer: mo.Some("new_implementer"),
								},
							)
							require.NoError(t, err)

							created.Assignees.Implementer = "new_implementer"
							created.Report.Date = updated.GetDate()

							require.Equal(
								t,
								created,
								updated,
							)
						},
					)

					t.Run(
						"All",
						func(t *testing.T) {
							model, err := report.NewReport(
								"some_name",
								"some_text",
								"some_reporter",
								"some_implementer",
							)
							require.NoError(t, err)

							created, err := service.CreateReport(
								context.Background(),
								model,
							)
							require.NoError(t, err)

							updated, err := service.UpdateReport(
								context.Background(),
								created.GetID(),
								reportdomain.UpdateReportParams{
									Implementer: mo.Some("new_implementer"),
									Name:        mo.Some("new_name"),
									Text:        mo.Some("new_text"),
								},
							)
							require.NoError(t, err)

							created.Assignees.Implementer = "new_implementer"
							created.Report.Name = "new_name"
							created.Report.Text = "new_text"
							created.Report.Date = updated.GetDate()

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
				"Failure",
				func(t *testing.T) {
					t.Run(
						"NotFoundOrInvalidID",
						func(t *testing.T) {
							_, err := service.UpdateReport(
								context.Background(),
								"some_id",
								reportdomain.UpdateReportParams{},
							)
							require.Condition(
								t,
								func() (success bool) {
									return err == reports.ErrReportIDNotValid || err == reports.ErrReportNotFound
								},
							)
						},
					)

					t.Run(
						"Validations",
						func(t *testing.T) {
							t.Run(
								"EmptyName",
								func(t *testing.T) {
									model, err := report.NewReport(
										"some_name",
										"some_text",
										"some_reporter",
										"some_implementer",
									)
									require.NoError(t, err)

									created, err := service.CreateReport(
										context.Background(),
										model,
									)
									require.NoError(t, err)

									_, err = service.UpdateReport(
										context.Background(),
										created.GetID(),
										reportdomain.UpdateReportParams{
											Name: mo.Some(""),
										},
									)
									require.Condition(
										t,
										func() (success bool) {
											return errors.Is(err, reports.ErrValidationError) && errors.Is(err, report.ErrNameEmpty)
										},
									)

								},
							)

							t.Run(
								"EmptyText",
								func(t *testing.T) {
									model, err := report.NewReport(
										"some_name",
										"some_text",
										"some_reporter",
										"some_implementer",
									)
									require.NoError(t, err)

									created, err := service.CreateReport(
										context.Background(),
										model,
									)
									require.NoError(t, err)

									_, err = service.UpdateReport(
										context.Background(),
										created.GetID(),
										reportdomain.UpdateReportParams{
											Text: mo.Some(""),
										},
									)
									require.Condition(
										t,
										func() (success bool) {
											return errors.Is(err, reports.ErrValidationError) && errors.Is(err, report.ErrTextEmpty)
										},
									)
								},
							)

							t.Run(
								"EmptyImplementor",
								func(t *testing.T) {
									model, err := report.NewReport(
										"some_name",
										"some_text",
										"some_reporter",
										"some_implementer",
									)
									require.NoError(t, err)

									created, err := service.CreateReport(
										context.Background(),
										model,
									)
									require.NoError(t, err)

									_, err = service.UpdateReport(
										context.Background(),
										created.GetID(),
										reportdomain.UpdateReportParams{
											Implementer: mo.Some(""),
										},
									)
									require.Condition(
										t,
										func() (success bool) {
											return errors.Is(err, reports.ErrValidationError) && errors.Is(err, report.ErrImplementorEmpty)
										},
									)
								},
							)
						},
					)
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
						"some_name",
						"some_text",
						"some_reporter",
						"some_implementer",
					)
					require.NoError(t, err)

					created, err := service.CreateReport(
						context.Background(),
						model,
					)
					require.NoError(t, err)

					get, err := service.GetReport(
						context.Background(),
						created.GetID(),
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
				"Failure",
				func(t *testing.T) {
					t.Run(
						"InvalidIDOrNotFind",
						func(t *testing.T) {
							_, err := service.GetReport(
								context.Background(),
								"some_id",
							)
							require.Condition(
								t,
								func() (success bool) {
									return err == reports.ErrReportNotFound || err == reports.ErrReportIDNotValid
								},
							)
						},
					)
				},
			)
		},
	)

	t.Run(
		"GetReports",
		func(t *testing.T) {
			t.Run(
				"Fail",
				func(t *testing.T) {
					t.Run(
						"EmptySortParams",
						func(t *testing.T) {
							require.Condition(
								t,
								func() (success bool) {
									_, err := service.GetReports(
										context.Background(),
										&reportdomain.GetReportsParams{
											Filter: &reportdomain.GetReportsFilter{
												SortParams: []reportdomain.GetReportsSort{
													{
														NameSort: mo.Some[ordertype.OrderType](ordertype.ASC),
													},
													{},
												},
											},
										},
									)
									return errors.Is(err, reports.ErrGetReportsBadParams) && errors.Unwrap(err).Error() == "empty sort argument"
								},
							)
						},
					)

					t.Run(
						"MoreThanOneFieldSetInItem",
						func(t *testing.T) {
							require.Condition(
								t,
								func() (success bool) {
									_, err := service.GetReports(
										context.Background(),
										&reportdomain.GetReportsParams{
											Filter: &reportdomain.GetReportsFilter{
												SortParams: []reportdomain.GetReportsSort{
													{
														NameSort: mo.Some[ordertype.OrderType](ordertype.ASC),
														DateSort: mo.Some[ordertype.OrderType](ordertype.ASC),
													},
												},
											},
										},
									)
									return errors.Is(err, reports.ErrGetReportsBadParams) && errors.Unwrap(err).Error() == "only one sort field can be set in slice item"
								},
							)
						},
					)

					t.Run(
						"IsDuplicatedFields",
						func(t *testing.T) {
							require.Condition(
								t,
								func() (success bool) {
									_, err := service.GetReports(
										context.Background(),
										&reportdomain.GetReportsParams{
											Filter: &reportdomain.GetReportsFilter{
												SortParams: []reportdomain.GetReportsSort{
													{
														NameSort: mo.Some[ordertype.OrderType](ordertype.ASC),
													},
													{
														DateSort: mo.Some[ordertype.OrderType](ordertype.ASC),
													},
													{
														NameSort: mo.Some[ordertype.OrderType](ordertype.ASC),
													},
												},
											},
										},
									)
									return errors.Is(err, reports.ErrGetReportsBadParams) && errors.Unwrap(err).Error() == "you can't sort by field twice"
								},
							)
						},
					)
				},
			)
		},
	)
}

func TestFunc_ReportService(t *testing.T) {
	config := config.GetConfigFrom("./../../.env")
	service, err := reportservice.New(
		reportservice.WithMongoRepository(config.MongoDB.TestURI),
	)
	require.NoError(t, err)

	Service_Tests(t, service)
}
