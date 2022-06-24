package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetReportReq pb.GetReportReq

func (g *GetReportReq) ToEndpointReq() *reqresp.GetReportReq {
	return &reqresp.GetReportReq{
		ID: g.ReportId,
	}
}

func DecodeGetReportReq(
	ctx context.Context,
	grpcReq *pb.GetReportReq,
) (*GetReportReq, error) {
	return (*GetReportReq)(grpcReq), nil
}