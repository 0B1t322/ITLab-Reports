package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetDraftsReq struct {
	UserID string
}

func (g *GetDraftsReq) SetUserID(id string) {
	g.UserID = id
}

func (g *GetDraftsReq) ToEndpointReq() *reqresp.GetReportsReq {
	return &reqresp.GetReportsReq{
		Params: &report.GetReportsParams{
			Filter: &report.GetReportsFilter{
				GetReportsFilterFieldsWithOrAnd: report.GetReportsFilterFieldsWithOrAnd{
					GetReportsFilterFields: report.GetReportsFilterFields{
						Reporter: &filter.FilterField[string]{
							Value:     g.UserID,
							Operation: filter.EQ,
						},
					},
				},
			},
		},
	}
}

func DecodeGetDraftsReq(
	ctx context.Context,
	grpcReq *pb.GetDraftsReq,
) (*GetDraftsReq, error) {
	return &GetDraftsReq{}, nil
}
