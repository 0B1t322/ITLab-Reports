package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetReportImplementerReq pb.GetReportImplementerReq

func (g *GetReportImplementerReq) ToEndpointReq() *reqresp.GetReportReq {
	return &reqresp.GetReportReq{
		ID: g.ReportId,
	}
}

func DecodeGetReportImplementerReq(
	ctx context.Context,
	grpcReq *pb.GetReportImplementerReq,
) (*GetReportImplementerReq, error) {
	return (*GetReportImplementerReq)(grpcReq), nil
}