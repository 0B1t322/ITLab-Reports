package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/utils"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetDraftsResp pb.GetDraftsResp

func GetDraftsRespFrom(from *reqresp.GetReportsResp) *GetDraftsResp {
	var drafts []*types.Draft
	{
		for _, d := range from.Reports {
			drafts = append(drafts, utils.DraftToPBType(d))
		}
	}

	return &GetDraftsResp{
		Drafts: drafts,
	}
}

func EncodeGetDraftsResp(
	ctx context.Context,
	resp *GetDraftsResp,
) (*pb.GetDraftsResp, error) {
	return (*pb.GetDraftsResp)(resp), nil
}
