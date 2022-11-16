package reports

import (
	"time"

	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/reports/v2/dto"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/match"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/sort"
	reportsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type ReportsMarshaller struct {
	matchMarshaller match.MatchMarshaler
	sortMarshaller  sort.SortMarshaler
}

func (ReportsMarshaller) MarshallGetReportsResp(
	req GetReportsReq,
	reports []aggregate.Report,
	count int64,
) GetReportsResp {
	return GetReportsResp{
		Count:       int64(len(reports)),
		Reports:     ReportsViewFrom(reports),
		HasMore:     HasMore(count, req.Limit, req.Offset),
		Limit:       req.Limit,
		Offset:      req.Offset,
		TotalResult: count,
	}
}

func (r ReportsMarshaller) MarshallGetReportsReq(
	req GetReportsReq,
	user aggregate.User,
) (reports.GetReportsReq, error) {
	filter, err := r.marshallGetReportsFilterQuery(
		req,
		user,
	)
	if err != nil {
		return reports.GetReportsReq{}, err
	}

	sort, err := r.marshallGetReportsSort(req)
	if err != nil {
		return reports.GetReportsReq{}, err
	}

	request := reports.GetReportsReq{
		Query: reportsrepo.GetReportsQuery{
			Filter: filter,
			Sort:   sort,
			Limit: lo.Ternary(
				req.Limit == 0,
				mo.None[int64](),
				mo.Some(req.Limit),
			),
			Offset: lo.Ternary(
				req.Offset == 0,
				mo.None[int64](),
				mo.Some(req.Offset),
			),
		},
	}

	return request, nil
}

func (r ReportsMarshaller) marshallGetReportsFilterQuery(
	req GetReportsReq,
	user aggregate.User,
) (reportsrepo.FilterQuery, error) {
	filterQuery := reportsrepo.FilterQuery{}
	{
		// If user sets reporters for him
		if user.IsUser() {
			filterQuery.Or = append(
				filterQuery.Or,
				reportsrepo.Query().Expression(
					reportsrepo.Expression().
						Reporter(user.ID, filter.EQ),
				).Build(),
				reportsrepo.Query().Expression(
					reportsrepo.Expression().
						Implementer(user.ID, filter.EQ),
				).Build(),
			)
		}

		if !req.DateBegin.Equal(time.Time{}) {
			filterQuery.And = append(filterQuery.And, reportsrepo.Query().Expression(
				reportsrepo.Expression().
					Date(req.DateBegin, filter.GTE),
			).Build())
		}

		if !req.DateEnd.Equal(time.Time{}) {
			filterQuery.And = append(filterQuery.And, reportsrepo.Query().Expression(
				reportsrepo.Expression().
					Date(req.DateEnd, filter.LTE),
			).Build())
		}

		params := r.matchMarshaller.Marshal(req.Match)
		for _, param := range params {
			switch dto.ReportMatchFields(param.Field()) {
			case dto.ReportMatchFields_Name:
				value, err := dto.Name_MatchParam.ValueFromString(param.Value())
				if err != nil {
					continue
				}

				filterQuery.And = append(
					filterQuery.And,
					reportsrepo.Query().Expression(
						reportsrepo.Expression().Name(
							value,
							filter.EQ,
						),
					).Build(),
				)
			case dto.ReportMatchFields_Date:
				value, err := dto.Date_MatchParam.ValueFromString(param.Value())
				if err != nil {
					continue
				}

				filterQuery.And = append(
					filterQuery.And,
					reportsrepo.Query().Expression(
						reportsrepo.Expression().Date(
							value,
							filter.EQ,
						),
					).Build(),
				)

			case dto.ReportMatchFields_Reporter:
				value, err := dto.Reporter_MatchParam.ValueFromString(param.Value())
				if err != nil {
					continue
				}

				if !user.IsAdminOrSuperAdmin() && user.ID != value {
					continue
				}

				filterQuery.And = append(
					filterQuery.And,
					reportsrepo.Query().Expression(
						reportsrepo.Expression().Reporter(
							value,
							filter.EQ,
						),
					).Build(),
				)

			case dto.ReportMatchFields_Implementer:
				value, err := dto.Implementer_MatchParam.ValueFromString(param.Value())
				if err != nil {
					continue
				}

				if !user.IsAdminOrSuperAdmin() && user.ID != value {
					continue
				}

				filterQuery.And = append(
					filterQuery.And,
					reportsrepo.Query().Expression(
						reportsrepo.Expression().Implementer(
							value,
							filter.EQ,
						),
					).Build(),
				)
			}
		}

	}

	return filterQuery, nil
}

func (r ReportsMarshaller) marshallGetReportsSort(
	req GetReportsReq,
) ([]reportsrepo.SortFields, error) {
	params := r.sortMarshaller.Marshal(req.SortBy)
	if err := params.IsFieldsRepeat(); err != nil {
		return nil, err
	}

	builder := reportsrepo.SortBuilder()
	if len(params) == 0 {
		return builder.
			Date(sortorder.ASC).
			Build(), nil
	}

	for _, param := range params {
		switch param.Field() {
		case string(dto.ReportSortFields_Name):
			builder.Name(param.Order().ToSortOrder())
		case string(dto.ReportSortFields_Date):
			builder.Date(param.Order().ToSortOrder())
		}
	}

	return builder.Build(), nil
}
