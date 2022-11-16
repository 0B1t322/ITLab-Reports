package reports

import (
	"fmt"
	"strings"
	"time"

	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v1/dto"
	reportsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type ReportMarshaller struct {
}

func (ReportMarshaller) MarshallReportView(report aggregate.Report) ReportView {
	return ReportViewFrom(report)
}

func (ReportMarshaller) MarshallReportsView(reports []aggregate.Report) []ReportView {
	return ReportsViewFrom(reports)
}

func (ReportMarshaller) MarshallGetReportReq(
	req GetReportReq,
	user aggregate.User,
) reports.GetReportReq {
	return reports.GetReportReq{
		ID: req.ID,
		By: user,
	}
}

func (ReportMarshaller) MarshallGetReportsForEmployeeReq(
	req GetEmployeeReportsReq,
	user aggregate.User,
) reports.GetReportsReq {
	query := reportsrepo.Query().
		Or(
			reportsrepo.Expression().
				Reporter(
					req.EmployeeID,
					filter.EQ,
				),
		).
		Or(
			reportsrepo.Expression().
				Implementer(
					req.EmployeeID,
					filter.EQ,
				),
		).
		Build()

	if !req.DateBegin.UTC().Equal(time.Time{}) {
		query.And = append(
			query.And,
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().
						Date(
							req.DateBegin.UTC(),
							filter.GTE,
						),
				).
				Build(),
		)
	}

	if !req.DateEnd.UTC().Equal(time.Time{}) {
		query.And = append(
			query.And,
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().
						Date(
							req.DateEnd.UTC(),
							filter.LTE,
						),
				).
				Build(),
		)
	}

	return reports.GetReportsReq{
		Query: reportsrepo.GetReportsQuery{
			Filter: query,
			Sort: reportsrepo.SortBuilder().
				Date(sortorder.ASC).
				Build(),
		},
	}
}

func (ReportMarshaller) MarshallGetReportsReq(
	req GetReportsReq,
	user aggregate.User,
) reports.GetReportsReq {
	var query reportsrepo.FilterQuery
	{
		if user.IsUser() {
			query = reportsrepo.Query().
				Or(
					reportsrepo.Expression().
						Reporter(
							user.ID,
							filter.EQ,
						),
				).
				Or(
					reportsrepo.Expression().
						Implementer(
							user.ID,
							filter.EQ,
						),
				).
				Build()
		} else if user.IsAdminOrSuperAdmin() {
			query = reportsrepo.Query().Build()
		}
	}

	var sort []reportsrepo.SortFields
	{
		switch req.SortedBy {
		case dto.ReportsSortedByFields_Name:
			sort = reportsrepo.SortBuilder().
				Name(sortorder.ASC).
				Build()
		default:
			sort = reportsrepo.SortBuilder().
				Date(sortorder.ASC).
				Build()
		}
	}

	return reports.GetReportsReq{
		Query: reportsrepo.GetReportsQuery{
			Filter: query,
			Sort:   sort,
		},
	}
}

func (ReportMarshaller) MarshallCreateReportReq(
	req CreateReportReq,
	user aggregate.User,
) (reports.CreateReportReq, error) {
	var (
		name string
		text string
	)

	if req.Name != "" {
		name = req.Name
		text = req.Text
	} else {
		splitedText := strings.Split(req.Text, "@\n\t\n@")
		if len(splitedText) < 2 {
			return reports.CreateReportReq{}, fmt.Errorf("invalid text and name format")
		}

		name = splitedText[0]
		text = splitedText[1]
	}

	return reports.CreateReportReq{
		Name: name,
		Text: text,
		Implementer: lo.Ternary(
			req.Implementer != nil,
			mo.Some(lo.FromPtr(req.Implementer)),
			mo.None[string](),
		),
		By: user,
	}, nil
}

func (ReportMarshaller) MarshallCreateReportFromDraftReq(
	req CreateReportFromDraftReq,
	user aggregate.User,
) reports.CreateReportFromDraftReq {
	return reports.CreateReportFromDraftReq{
		DraftID: req.DraftID,
		By:      user,
	}
}
