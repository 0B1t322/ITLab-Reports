package report_test

import (
	"testing"

	model "github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/entity/assignees"
	"github.com/RTUITLab/ITLab-Reports/entity/report"
	"github.com/stretchr/testify/require"
)

func TestFunc_ReportCreate(t *testing.T) {
	t.Run(
		"NewReport",
		func(t *testing.T) {
			r, err := model.NewReport(
				"report_name",
				"some_text",
				"reporter_id",
				"implementor_id",
			)
			require.NoError(t, err)

			require.Equal(
				t,
				&model.Report{
					Report: &report.Report{
						ID: "",
						Name: "report_name",
						Text: "some_text",
						Date: r.Report.Date,
						State: report.ReportStateCreated,
					},
					Assignees: &assignees.Assignees{
						Reporter: "reporter_id",
						Implementer: "implementor_id",
					},
				},
				r,
			)
		},
	)

	t.Run(
		"NewReport_WithError",
		func(t *testing.T) {
			t.Run(
				"EmptyName",
				func(t *testing.T) {
					_, err := model.NewReport(
						"",
						"some_text",
						"reporter_id",
						"implementor_id",
					)
					require.ErrorIs(t, err, model.ErrNameEmpty)
				},
			)

			t.Run(
				"EmptyText",
				func(t *testing.T) {
					_, err := model.NewReport(
						"name",
						"",
						"reporter_id",
						"implementor_id",
					)
					require.ErrorIs(t, err, model.ErrTextEmpty)
				},
			)

			t.Run(
				"EmptyReportId",
				func(t *testing.T) {
					_, err := model.NewReport(
						"name",
						"some_text",
						"",
						"implementor_id",
					)
					require.ErrorIs(t, err, model.ErrReporterEmpty)
				},
			)

			t.Run(
				"EmptyImplementId",
				func(t *testing.T) {
					_, err := model.NewReport(
						"name",
						"some_text",
						"report_id",
						"",
					)
					require.ErrorIs(t, err, model.ErrImplementorEmpty)
				},
			)
		},
	)
}

func TestFunc_ReportBuilder(t *testing.T) {
	t.Run(
		"Create",
		func(t *testing.T) {
			r, err := model.NewReportBuilder().
						SetName("some_name").
						SetText("some_text").
						SetReporter("reporter_id").
						SetImplementor("implementor_id").
						Create()
			require.NoError(t, err)


			require.Equal(
				t,
				&model.Report{
					Report: &report.Report{
						ID: "",
						Name: "some_name",
						Text: "some_text",
						Date: r.Report.Date,
					},
					Assignees: &assignees.Assignees{
						Reporter: "reporter_id",
						Implementer: "implementor_id",
					},
				},
				r,
			)
		},
	)

	t.Run(
		"WithError",
		func(t *testing.T) {
			t.Run(
				"EmptyName",
				func(t *testing.T) {
					_, err := model.NewReportBuilder().
						SetText("some_text").
						SetReporter("reporter_id").
						SetImplementor("implementor_id").
						Create()
					require.ErrorIs(t, err, model.ErrNameEmpty)
				},
			)

			t.Run(
				"EmptyText",
				func(t *testing.T) {
					_, err := model.NewReportBuilder().
						SetName("some_name").
						SetReporter("reporter_id").
						SetImplementor("implementor_id").
						Create()
					require.ErrorIs(t, err, model.ErrTextEmpty)
				},
			)

			t.Run(
				"EmptyReportId",
				func(t *testing.T) {
					_, err := model.NewReportBuilder().
						SetName("some_name").
						SetText("some_text").
						SetImplementor("implementor_id").
						Create()
					require.ErrorIs(t, err, model.ErrReporterEmpty)
				},
			)

			t.Run(
				"EmptyImplementId",
				func(t *testing.T) {
					_, err := model.NewReportBuilder().
						SetName("some_name").
						SetText("some_text").
						SetReporter("reporter_id").
						Create()
					require.ErrorIs(t, err, model.ErrImplementorEmpty)
				},
			)
		},
	)
}