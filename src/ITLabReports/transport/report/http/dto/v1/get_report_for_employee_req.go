package dto

import (
	"context"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/gorilla/mux"
)

type GetReportsForEmployeeReq struct {
	dateBegin string
	dateEnd   string
	employee  string
}

func (g *GetReportsForEmployeeReq) GetEmployee() string {
	return g.employee
}

func (g *GetReportsForEmployeeReq) ToEndpointReq() *reqresp.GetReportsReq {
	req := &reqresp.GetReportsReq{
		Params: &report.GetReportsParams{
			Filter: &report.GetReportsFilter{
				GetReportsFilterFieldsWithOrAnd: report.GetReportsFilterFieldsWithOrAnd{
					Or: []*report.GetReportsFilterFieldsWithOrAnd{
						{
							GetReportsFilterFields: report.GetReportsFilterFields{
								Implementer: &filter.FilterField[string]{
									Value:     g.employee,
									Operation: filter.EQ,
								},
							},
						},
						{
							GetReportsFilterFields: report.GetReportsFilterFields{
								Reporter: &filter.FilterField[string]{
									Value:     g.employee,
									Operation: filter.EQ,
								},
							},
						},
					},
				},
			},
		},
	}

	var dateAnd []*report.GetReportsFilterFieldsWithOrAnd
	{
		if dateBegin := g.dateBegin; dateBegin != "" {
			dateAnd = append(
				dateAnd,
				&report.GetReportsFilterFieldsWithOrAnd{
					GetReportsFilterFields: report.GetReportsFilterFields{
						Date: &filter.FilterField[string]{
							Value:     dateBegin,
							Operation: filter.GTE,
						},
					},
				},
			)
		}

		if dateEnd := g.dateEnd; dateEnd != "" {
			dateAnd = append(
				dateAnd,
				&report.GetReportsFilterFieldsWithOrAnd{
					GetReportsFilterFields: report.GetReportsFilterFields{
						Date: &filter.FilterField[string]{
							Value:     dateEnd,
							Operation: filter.LTE,
						},
					},
				},
			)
		}
	}

	if len(dateAnd) != 0 {
		req.Params.Filter.And = append(req.Params.Filter.And, dateAnd...)
	}
	return req
}

func DecodeGetReportsForEmployeeReq(
	ctx context.Context,
	r *http.Request,
) (*GetReportsForEmployeeReq, error) {
	vars := mux.Vars(r)

	req := &GetReportsForEmployeeReq{
		dateBegin: vars["dateBegin"],
		dateEnd:   vars["dateEnd"],
		employee:  vars["employee"],
	}

	return req, nil
}
