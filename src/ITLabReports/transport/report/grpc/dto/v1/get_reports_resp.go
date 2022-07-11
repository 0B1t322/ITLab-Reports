package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/utils"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetReportsResp pb.GetReportsResp

func GetReportsRespFrom(resp *reqresp.GetReportsResp) *GetReportsResp {
	var reports []*types.Report
	{
		for _, r := range resp.Reports {
			reports = append(reports, utils.ReportToPBType(r))
		}
	}

	return &GetReportsResp{
		Reports: reports,
	}
}

func EncodeGetReportsResp(
	ctx context.Context,
	resp *GetReportsResp,
) (*pb.GetReportsResp, error) {
	return (*pb.GetReportsResp)(resp), nil
}