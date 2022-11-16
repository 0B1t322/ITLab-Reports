package drafts

import (
	"github.com/0B1t322/RepoGen/pkg/filter"
	draftsrepo "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/repository"
	drafts "github.com/RTUITLab/ITLab-Reports/internal/domain/drafts/service"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type DraftsMarshaller struct {
}

func (DraftsMarshaller) MarshallGetDraftReq(
	req GetDraftReq,
	user aggregate.User,
) drafts.GetDraftReq {
	return drafts.GetDraftReq{
		ID: req.ID,
		By: user,
	}
}

func (DraftsMarshaller) MarshallGetDraftsReq(
	req GetDraftsReq,
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
		},
	}
}

func (DraftsMarshaller) MarshallDeleteDraftReq(
	req DeleteDraftReq,
	user aggregate.User,
) drafts.DeleteDraftReq {
	return drafts.DeleteDraftReq{
		ID: req.ID,
		By: user,
	}
}

func (DraftsMarshaller) MarshallUpdateDraftReq(
	req UpdateDraftReq,
	user aggregate.User,
) drafts.UpdateDraftReq {
	return drafts.UpdateDraftReq{
		ID: req.ID,
		Name: lo.Ternary(
			req.Name != "",
			mo.Some(req.Name),
			mo.None[string](),
		),
		Text: lo.Ternary(
			req.Text != "",
			mo.Some(req.Text),
			mo.None[string](),
		),
		Implementer: lo.Ternary(
			req.Implementer != "",
			mo.Some(req.Implementer),
			mo.None[string](),
		),
		By: user,
	}
}

func (DraftsMarshaller) MarshallCreateDraftReq(
	req CreateDraftReq,
	user aggregate.User,
) drafts.CreateDraftReq {
	return drafts.CreateDraftReq{
		Name: req.Name,
		Text: req.Text,
		Implementer: lo.Ternary(
			req.Implementer != "",
			mo.Some(req.Implementer),
			mo.None[string](),
		),
		By: user,
	}
}
