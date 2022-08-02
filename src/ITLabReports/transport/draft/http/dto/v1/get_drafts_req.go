package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/samber/mo"
)

type GetDraftsReq struct {
	UserID	string
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
							Value: g.UserID,
							Operation: filter.EQ,
						},
					},
				},
				SortParams: []report.GetReportsSort{
					{
						DateSort: mo.Some[ordertype.OrderType](ordertype.ASC),
					},
				},
			},
		},
	}
}

func DecodeGetDraftsReq(
	ctx context.Context,
	r *http.Request,
) (*GetDraftsReq, error) {
	return &GetDraftsReq{}, nil
}