package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/gorilla/mux"
)

type GetReportsReq struct {
	sortedBy string

	implementer string

	reporter string
}

func (g *GetReportsReq) SetImplementerAndReporter(implementer, reporter string) {
	g.implementer = implementer
	g.reporter = reporter
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

	switch g.sortedBy {
	case "name":
		req.Params.Filter.NameSort = *optional.NewOptional[ordertype.OrderType](ordertype.ASC)
	case "date":
		req.Params.Filter.DateSort = *optional.NewOptional[ordertype.OrderType](ordertype.ASC)
	}

	if g.implementer != "" && g.reporter != "" {
		req.SetImplementerAndReporter(g.implementer, g.reporter)
	}


	return req
}

func DecodeGetReportsReq(
	ctx context.Context,
	r *http.Request,
) (*GetReportsReq, error) {
	vars := mux.Vars(r)

	var (
		sortedBy = vars["sorted_by"]
	)

	req := &GetReportsReq{
		sortedBy: sortedBy,
	}

	return req, nil
}