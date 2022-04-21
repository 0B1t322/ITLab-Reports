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
	DateBegin string
	DateEnd   string
	Employee  string
}

func (g *GetReportsForEmployeeReq) GetEmployee() string {
	return g.Employee
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
									Value:     g.Employee,
									Operation: filter.EQ,
								},
							},
						},
						{
							GetReportsFilterFields: report.GetReportsFilterFields{
								Reporter: &filter.FilterField[string]{
									Value:     g.Employee,
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
		if dateBegin := g.DateBegin; dateBegin != "" {
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

		if dateEnd := g.DateEnd; dateEnd != "" {
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
	values := r.URL.Query()

	req := &GetReportsForEmployeeReq{
		DateBegin: values.Get("dateBegin"),
		DateEnd:   values.Get("dateEnd"),
		Employee:  vars["employee"],
	}

	return req, nil
}
