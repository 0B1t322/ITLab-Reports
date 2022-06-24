package dto

import (
	"context"

	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetReportImplementerResp struct {
	*pb.GetReportImplementerResp
	Reporter string
}

func (g *GetReportImplementerResp) GetImplementer() string {
	var implementer string
	{
		switch g := g.Result.(type) {
		case *pb.GetReportImplementerResp_Implementer:
			implementer = g.Implementer
		}
	}

	return implementer
}

func (g *GetReportImplementerResp) GetReporter() string {
	return g.Reporter
}

func (g *GetReportImplementerResp) IsError() bool {
	switch g.Result.(type) {
	case *pb.GetReportImplementerResp_Error:
		return true
	}
	return false
}

func EncodeGetReportImplementerResp(
	ctx context.Context,
	resp *GetReportImplementerResp,
) (*pb.GetReportImplementerResp, error) {
	return resp.GetReportImplementerResp, nil
}
