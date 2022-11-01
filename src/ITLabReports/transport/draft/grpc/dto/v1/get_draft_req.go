package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetDraftReq pb.GetDraftReq

func (g *GetDraftReq) ToEndpointReq() *reqresp.GetReportReq {
	return &reqresp.GetReportReq{
		ID: g.DraftId,
	}
}

func DecodeGetDraftReq(
	ctx context.Context,
	grpcReq *pb.GetDraftReq,
) (*GetDraftReq, error) {
	return (*GetDraftReq)(grpcReq), nil
}
