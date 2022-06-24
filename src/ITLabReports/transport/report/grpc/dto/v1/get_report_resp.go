package dto

import (
	"context"

	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetReportResp pb.GetReportResp

func (g *GetReportResp) GetImplementer() string {
	var implementer string
	{
		switch g := g.Result.(type) {
		case *pb.GetReportResp_Report:
			implementer = g.Report.Assignees.Implementer
		}
	}

	return implementer
}

func (g *GetReportResp) GetReporter() string {
	var reporter string
	{
		switch g := g.Result.(type) {
		case *pb.GetReportResp_Report:
			reporter = g.Report.Assignees.Reporter
		}
	}

	return reporter
}

func (g *GetReportResp) IsError() bool {
	switch g.Result.(type) {
	case *pb.GetReportResp_Error:
		return true
	}
	return false
}

func EncodeGetReportResp(
	ctx context.Context,
	resp *GetReportResp,
) (*pb.GetReportResp, error) {
	return (*pb.GetReportResp)(resp), nil
}