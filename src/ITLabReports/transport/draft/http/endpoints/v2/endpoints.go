package endpoints

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	v1 "github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type GetDrafts = endpoint.Endpoint[*dto.GetDraftsReq, *dto.GetDraftsResp]

type Endpoints struct {
	GetDrafts GetDrafts
}

func NewEndpoints(
	e report.Endpoints,
) Endpoints {
	return Endpoints{
		GetDrafts: makeGetDraftsEndpoint(e),
	}
}

func makeGetDraftsEndpoint(
	e report.Endpoints,
) GetDrafts {
	return func(ctx context.Context, req *dto.GetDraftsReq) (*dto.GetDraftsResp, error) {
		drafts, err := e.GetReports(ctx, req.ToEndpointReq())
		if err != nil {
			return nil, err
		}

		countDraft, err := e.CountReports(
			ctx,
			&reqresp.CountReportsReq{
				Params: &req.Query.Params.Filter.GetReportsFilterFieldsWithOrAnd,
			},
		)
		if err != nil {
			return nil, err
		}

		offset := 0
		limit := 0
		count := len(drafts.Reports)
		totalResult := countDraft.Count

		if req.Query.Params.Limit.IsPresent() {
			limit = int(req.Query.Params.Limit.MustGet())
		}

		if req.Query.Params.Offset.IsPresent() {
			offset = int(req.Query.Params.Offset.MustGet())
		}

		HasMore := false
		{
			if limit != 0 && totalResult-int64(offset)-int64(limit) > 0 {
				HasMore = true
			}
		}

		return &dto.GetDraftsResp{
			Count:       count,
			HasMore:     HasMore,
			Items:       v1.GetDraftsRespFrom(drafts).Drafts,
			Limit:       limit,
			Offset:      offset,
			TotalResult: int(totalResult),
		}, nil

	}
}
