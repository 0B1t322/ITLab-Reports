package dto

import (
	"context"

	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetDraftsPaginatedResp pb.GetDraftsPaginatedResp

func EncodeGetDraftsPaginatedResp(
	ctx context.Context,
	resp *GetDraftsPaginatedResp,
) (*pb.GetDraftsPaginatedResp, error) {
	return (*pb.GetDraftsPaginatedResp)(resp), nil
}
