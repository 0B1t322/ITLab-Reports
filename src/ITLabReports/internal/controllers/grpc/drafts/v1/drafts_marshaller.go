package drafts

import (
	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	draftsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/repository"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	reportsgrpc "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/RTUITLab/ITLab/proto/shared"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type DraftMarshaller struct {
}

func (DraftMarshaller) MarshallGetDraftsReq(
	req *reportsgrpc.GetDraftsReq,
	user aggregate.User,
) drafts.GetDraftsReq {
	return drafts.GetDraftsReq{
		Query: draftsrepo.GetDraftsQuery{
			Filter: draftsrepo.Query().
				Expression(
					draftsrepo.Expression().Reporter(
						user.ID,
						filter.EQ,
					),
				).
				Build(),
			Sort: draftsrepo.SortBuilder().
				Date(sortorder.DESC).
				Build(),
		},
	}
}

func (DraftMarshaller) MarshallGetDraftsPaginatedReq(
	req *reportsgrpc.GetDraftsPaginatedReq,
	user aggregate.User,
) drafts.GetDraftsReq {
	var (
		limit  int64
		offset int64
	)
	{
		if pagination := req.GetPagination(); pagination != nil {
			limit = pagination.GetLimit()
			offset = pagination.GetOffset()
		}
	}

	return drafts.GetDraftsReq{
		Query: draftsrepo.GetDraftsQuery{
			Filter: draftsrepo.Query().
				Expression(
					draftsrepo.Expression().Reporter(
						user.ID,
						filter.EQ,
					),
				).
				Build(),
			Sort: draftsrepo.SortBuilder().Build(),
			Limit: lo.Ternary(
				limit != 0,
				mo.Some(limit),
				mo.None[int64](),
			),
			Offset: lo.Ternary(
				offset != 0,
				mo.Some(offset),
				mo.None[int64](),
			),
		},
	}
}

func (DraftMarshaller) MarshallGetDraftsResp(
	req *reportsgrpc.GetDraftsPaginatedReq,
	drafts []aggregate.Draft,
	count int64,
) *reportsgrpc.GetDraftsPaginatedResp {
	var (
		offset int64
		limit  int64
	)
	{
		if pagination := req.GetPagination(); pagination != nil {
			limit = pagination.GetLimit()
			offset = pagination.GetOffset()
		}
	}
	return &reportsgrpc.GetDraftsPaginatedResp{
		Drafts: lo.Map(
			drafts,
			func(draft aggregate.Draft, _ int) *types.Draft {
				return DraftFrom(draft)
			},
		),
		PaginationInfo: &shared.PaginationInfo{
			Count:       int64(len(drafts)),
			HasMore:     HasMore(count, limit, offset),
			Limit:       limit,
			Offset:      offset,
			TotalResult: count,
		},
	}
}
