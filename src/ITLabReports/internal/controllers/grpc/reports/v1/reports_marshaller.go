package reports

import (
	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	reportsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/repository"
	reports "github.com/RTUITLab/ITLab-Reports/internal/domain/reports/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	reportsgrpc "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/RTUITLab/ITLab/proto/shared"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type ReportsMarshaller struct{}

func (ReportsMarshaller) MarshalGetReportsReq(
	req *reportsgrpc.GetReportsReq,
	user aggregate.User,
) reports.GetReportsReq {
	r := reports.GetReportsReq{
		Query: reportsrepo.GetReportsQuery{},
	}

	sort := reportsrepo.SortBuilder()

	switch req.SortedBy {
	case reportsgrpc.GetReportsReq_DATE.Enum():
		r.Query.Sort = sort.Date(sortorder.DESC).Build()
	case reportsgrpc.GetReportsReq_NAME.Enum():
		r.Query.Sort = sort.Name(sortorder.DESC).Build()
	}

	if user.IsUser() {
		r.Query.Filter.Or = append(
			r.Query.Filter.Or,
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().
						Reporter(user.ID, filter.EQ),
				).
				Build(),
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().
						Implementer(user.ID, filter.EQ),
				).
				Build(),
		)
	}

	return r
}

func (rm ReportsMarshaller) MarshallGetReportsPaginatedReq(
	req *reportsgrpc.GetReportsPaginatedReq,
	user aggregate.User,
) reports.GetReportsReq {
	var query reportsrepo.GetReportsQuery
	{
		if params := req.FilterParams; params != nil {
			query.Filter = rm.marshallFilterParams(params, user)
		}

		query.Sort = rm.marshallSortParams(req.OrderParams)

		if pagination := req.Pagination; pagination != nil {
			query.Limit = lo.Ternary(
				pagination.Limit != 0,
				mo.Some(pagination.Limit),
				mo.None[int64](),
			)

			query.Offset = lo.Ternary(
				pagination.Offset != 0,
				mo.Some(pagination.Offset),
				mo.None[int64](),
			)
		}
	}

	return reports.GetReportsReq{
		Query: query,
	}
}

func (rm ReportsMarshaller) marshallFilterParams(
	p *reportsgrpc.GetReportsPaginatedReq_FilterParams,
	user aggregate.User,
) (q reportsrepo.FilterQuery) {
	b := reportsrepo.Expression()
	if p.DateBegin != nil {
		q.And = append(
			q.And,
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().Date(
						p.DateBegin.AsTime(),
						filter.GTE,
					),
				).
				Build(),
		)
	}

	if p.DateEnd != nil {
		q.And = append(
			q.And,
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().Date(
						p.DateEnd.AsTime(),
						filter.LTE,
					),
				).
				Build(),
		)
	}

	if name := p.GetNameMatch(); name != "" {
		b = b.Name(
			name,
			filter.LIKE,
		)
	}

	if user.IsUser() {
		q.Expression = b.Build()
		q.Or = append(
			q.Or,
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().Reporter(
						user.ID,
						filter.EQ,
					),
				).
				Build(),
			reportsrepo.Query().
				Expression(
					reportsrepo.Expression().Implementer(
						user.ID,
						filter.EQ,
					),
				).
				Build(),
		)
		return
	}

	if id := p.GetImplementerId(); id != "" {
		b = b.Implementer(
			id,
			filter.EQ,
		)
	}

	if id := p.GetReporterId(); id != "" {
		b = b.Reporter(
			id,
			filter.EQ,
		)
	}

	q.Expression = b.Build()

	return
}

func (rm ReportsMarshaller) marshallSortParams(
	params []*reportsgrpc.GetReportsPaginatedReq_OrderParams,
) []reportsrepo.SortFields {
	b := reportsrepo.SortBuilder()

	for _, param := range params {
		switch param.Field {
		case reportsgrpc.GetReportsPaginatedReq_OrderParams_DATE:
			b = b.Date(
				ProtoSortOrderTo(param.Ordering),
			)
		case reportsgrpc.GetReportsPaginatedReq_OrderParams_NAME:
			b = b.Name(
				ProtoSortOrderTo(param.Ordering),
			)
		}
	}

	return b.Build()
}

func (rm ReportsMarshaller) MarshalGetReportsPaginatedResp(
	req *reportsgrpc.GetReportsPaginatedReq,
	reports []aggregate.Report,
	count int64,
) *reportsgrpc.GetReportsPaginatedResp {
	var (
		limit  int64
		offset int64
	)
	{
		if req.Pagination != nil {
			limit = req.Pagination.Limit
			offset = req.Pagination.Offset
		}
	}
	return &reportsgrpc.GetReportsPaginatedResp{
		Reports: lo.Map(
			reports,
			func(r aggregate.Report, _ int) *types.Report {
				return ReportFrom(r)
			},
		),
		PaginationInfo: &shared.PaginationInfo{
			Count:       int64(len(reports)),
			HasMore:     HasMore(count, limit, offset),
			Limit:       limit,
			Offset:      offset,
			TotalResult: count,
		},
	}
}
