package drafts

import (
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	draftsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/repository"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type DraftsMarshaller struct {
}

func (DraftsMarshaller) MarshallGetDraftsReq(
	req GetDraftsReq,
	user aggregate.User,
) drafts.GetDraftsReq {
	return drafts.GetDraftsReq{
		Query: draftsrepo.GetDraftsQuery{
			Limit: lo.Ternary(
				req.Limit != 0,
				mo.Some(req.Limit),
				mo.None[int64](),
			),
			Offset: lo.Ternary(
				req.Offset != 0,
				mo.Some(req.Offset),
				mo.None[int64](),
			),
			Sort: draftsrepo.SortBuilder().
				Date(sortorder.DESC).
				Build(),
		},
	}
}

func (DraftsMarshaller) MarshallGetDraftsResp(
	req GetDraftsReq,
	drafts []aggregate.Draft,
	count int64,
) GetDraftsResp {
	return GetDraftsResp{
		Drafts: lo.Map(
			drafts,
			func(draft aggregate.Draft, _ int) DraftView {
				return DraftViewFrom(draft)
			},
		),
		Count:      int64(len(drafts)),
		Offset:     req.Offset,
		Limit:      req.Limit,
		TotalCount: count,
		HasMore:    HasMore(count, req.Limit, req.Offset),
	}
}
