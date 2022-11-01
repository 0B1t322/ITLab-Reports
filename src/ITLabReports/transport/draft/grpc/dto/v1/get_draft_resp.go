package dto

import (
	"context"

	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetDraftResp pb.GetDraftResp

func (g *GetDraftResp) IsError() bool {
	switch g.Result.(type) {
	case *pb.GetDraftResp_Error:
		return true
	}
	return false
}

func (g *GetDraftResp) GetReporter() string {
	switch v := g.Result.(type) {
	case *pb.GetDraftResp_Draft:
		return v.Draft.GetAssignees().GetReporter()
	}
	return ""
}

func EncodeGetDraftResp(
	ctx context.Context,
	resp *GetDraftResp,
) (*pb.GetDraftResp, error) {
	return (*pb.GetDraftResp)(resp), nil
}
