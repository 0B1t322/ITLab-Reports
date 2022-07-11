package dto

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
)

type GetReportsReq struct {
	SortedBy string

	Implementer string

	Reporter string
}

func (g *GetReportsReq) SetImplementerAndReporter(implementer, reporter string) {
	g.Implementer = implementer
	g.Reporter = reporter
}

func (g *GetReportsReq) ToEndpointReq() *reqresp.GetReportsReq {
	req := &reqresp.GetReportsReq{
		Params: &report.GetReportsParams{
			Filter: &report.GetReportsFilter{
				GetReportsFilterFieldsWithOrAnd: report.GetReportsFilterFieldsWithOrAnd{},
			},
		},
	}

	switch g.SortedBy {
	case "name":
		req.Params.Filter.NameSort = *optional.NewOptional[ordertype.OrderType](ordertype.ASC)
	case "date":
		req.Params.Filter.DateSort = *optional.NewOptional[ordertype.OrderType](ordertype.ASC)
	}

	if g.Implementer != "" && g.Reporter != "" {
		req.SetImplementerAndReporter(g.Implementer, g.Reporter)
	}

	return req
}

func DecodeGetReportsListReq(
	ctx context.Context,
	grpcReq *pb.GetReportsReq,
) (*GetReportsReq, error) {
	var sortedBy string
	{
		if grpcReq.SortedBy == nil {
			sortedBy = "name"
		} else {
			switch *grpcReq.SortedBy {
			case pb.GetReportsReq_NAME:
				sortedBy = "name"
			case pb.GetReportsReq_DATE:
				sortedBy = "date"
			}
		}
	}
	return &GetReportsReq{
		SortedBy: sortedBy,
	}, nil
}
