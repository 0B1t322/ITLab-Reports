package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/samber/mo"
)

type GetDraftsPaginatedReq struct {
	Params *report.GetReportsParams
	UserID string
}

func (g *GetDraftsPaginatedReq) SetUserID(id string) {
	g.UserID = id
}

func (g *GetDraftsPaginatedReq) ToEndpointReq() *reqresp.GetReportsReq {
	g.Params.Filter.GetReportsFilterFieldsWithOrAnd.GetReportsFilterFields.Reporter = &filter.FilterField[string]{
		Value:     g.UserID,
		Operation: filter.EQ,
	}

	return &reqresp.GetReportsReq{
		Params: g.Params,
	}
}

func DecodeGetDraftsPaginatedReq(
	ctx context.Context,
	grpcReq *pb.GetDraftsPaginatedReq,
) (*GetDraftsPaginatedReq, error) {
	req := &GetDraftsPaginatedReq{
		Params: &report.GetReportsParams{
			Filter: &report.GetReportsFilter{},
		},
	}

	if pagination := grpcReq.GetPagination(); pagination != nil {
		if limit := pagination.GetLimit(); limit >= 1 {
			req.Params.Limit = mo.Some(limit)
		}

		if offset := pagination.GetOffset(); offset >= 0 {
			req.Params.Offset = mo.Some(offset)
		}
	}

	return req, nil
}
