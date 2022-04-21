package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
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
				GetReportsFilterFieldsWithOrAnd: report.GetReportsFilterFieldsWithOrAnd{
				},
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

func DecodeGetReportsReq(
	ctx context.Context,
	r *http.Request,
) (*GetReportsReq, error) {
	values := r.URL.Query()

	var (
		sortedBy = values.Get("sorted_by")
	)

	req := &GetReportsReq{
		SortedBy: sortedBy,
	}

	return req, nil
}