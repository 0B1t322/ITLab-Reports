package dto

import (
	"context"

	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetReportsPaginatedResp pb.GetReportsPaginatedResp

func EncodeGetReportsPaginatedResp(
	ctx context.Context,
	resp *GetReportsPaginatedResp,
) (*pb.GetReportsPaginatedResp, error) {
	return (*pb.GetReportsPaginatedResp)(resp), nil
}
