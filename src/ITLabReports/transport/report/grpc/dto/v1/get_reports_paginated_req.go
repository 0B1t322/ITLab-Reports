package dto

import (
	"context"
	"time"

	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/utils"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	"github.com/samber/mo"
)

type GetReportsPaginatedReq struct {
	Params        *report.GetReportsParams
	PaidState pb.GetReportsPaginatedReq_FilterParams_PaidState
}

func (g *GetReportsPaginatedReq) ToEndpointReq() *reqresp.GetReportsReq {
	return &reqresp.GetReportsReq{
		Params:        g.Params,
	}
}

func (g *GetReportsPaginatedReq) IsOnlyApprovedReports() bool {
	return g.PaidState == pb.GetReportsPaginatedReq_FilterParams_PAID
}

func (g *GetReportsPaginatedReq) SetOnlyApprovedReports(ids ...string) {
	g.Params.Filter.ReportsId = &filter.FilterField[[]string]{
		Value:     ids,
		Operation: filter.IN,
	}
}

func (g *GetReportsPaginatedReq) IsOnlyNotApprovedReports() bool {
	return g.PaidState == pb.GetReportsPaginatedReq_FilterParams_NOT_PAID
}

func (g *GetReportsPaginatedReq) SetOnlyNotApprovedReports(ids ...string) {
	g.Params.Filter.ReportsId = &filter.FilterField[[]string]{
		Value:     ids,
		Operation: filter.NIN,
	}
}

func (g *GetReportsPaginatedReq) SetImplementerAndReporter(implementer, reporter string) {
	// If this method call it's mean user don't have access to another reports
	// So nill filters
	g.Params.Filter.Implementer = nil
	g.Params.Filter.Reporter = nil

	g.Params.Filter.Or = append(
		g.Params.Filter.Or,
		&report.GetReportsFilterFieldsWithOrAnd{
			GetReportsFilterFields: report.GetReportsFilterFields{
				Reporter: &filter.FilterField[string]{
					Value:     reporter,
					Operation: filter.EQ,
				},
			},
		},
		&report.GetReportsFilterFieldsWithOrAnd{
			GetReportsFilterFields: report.GetReportsFilterFields{
				Reporter: &filter.FilterField[string]{
					Value:     implementer,
					Operation: filter.EQ,
				},
			},
		},
	)
}

func DecodeGetReportsPaginatedReq(
	ctx context.Context,
	grpcReq *pb.GetReportsPaginatedReq,
) (*GetReportsPaginatedReq, error) {
	req := &GetReportsPaginatedReq{
		Params: &report.GetReportsParams{
			Filter: &report.GetReportsFilter{},
		},
	}

	if pagination := grpcReq.GetPagination(); pagination != nil {
		// Limit
		if limit := pagination.GetLimit(); limit >= 1 {
			req.Params.Limit = mo.Some(limit)
		}

		// Offset
		if offset := pagination.GetOffset(); offset >= 0 {
			req.Params.Offset = mo.Some(offset)
		}
	}

	// Decode FilterParams
	if filterParams := grpcReq.FilterParams; filterParams != nil {
		if name := filterParams.GetNameMatch(); name != "" {
			req.Params.Filter.Name = &filter.FilterField[string]{
				Operation: filter.LIKE,
				Value:     name,
			}
		}

		if implementer := filterParams.GetImplementerId(); implementer != "" {
			req.Params.Filter.Implementer = &filter.FilterField[string]{
				Operation: filter.EQ,
				Value:     implementer,
			}
		}

		if reporter := filterParams.GetReporterId(); reporter != "" {
			req.Params.Filter.Reporter = &filter.FilterField[string]{
				Operation: filter.EQ,
				Value:     reporter,
			}
		}

		// Decode state
		if state := filterParams.PaidState; state != nil {
			req.PaidState = *state
		} else {
			req.PaidState = pb.GetReportsPaginatedReq_FilterParams_ALL
		}

		// Decode dateBegin
		if dateBegin := filterParams.DateBegin; dateBegin != nil {
			date := dateBegin.AsTime()
			req.Params.Filter.And = append(
				req.Params.Filter.And,
				&report.GetReportsFilterFieldsWithOrAnd{
					GetReportsFilterFields: report.GetReportsFilterFields{
						Date: &filter.FilterField[string]{
							Operation: filter.GTE,
							Value:     date.UTC().Format(time.RFC3339Nano),
						},
					},
				},
			)
		}

		// Decode dateEnd
		if dateEnd := filterParams.DateEnd; dateEnd != nil {
			date := dateEnd.AsTime()
			req.Params.Filter.And = append(
				req.Params.Filter.And,
				&report.GetReportsFilterFieldsWithOrAnd{
					GetReportsFilterFields: report.GetReportsFilterFields{
						Date: &filter.FilterField[string]{
							Operation: filter.LTE,
							Value:     date.UTC().Format(time.RFC3339Nano),
						},
					},
				},
			)
		}
	}

	// Decode sort
	if sort := grpcReq.OrderParams; sort != nil {
		var sortArgs []report.GetReportsSort
		{
			for _, sortField := range sort {
				if sortField.Field == pb.GetReportsPaginatedReq_OrderParams_DATE {
					sortArgs = append(
						sortArgs, 
						report.GetReportsSort{
							DateSort: mo.Some(utils.OrderTypeFromGRPC(sortField.Ordering)),
						},
					)
				}
				if sortField.Field == pb.GetReportsPaginatedReq_OrderParams_NAME {
					sortArgs = append(
						sortArgs, 
						report.GetReportsSort{
							NameSort: mo.Some(utils.OrderTypeFromGRPC(sortField.Ordering)),
						},
					)
				}
			}
		}
		req.Params.Filter.SortParams = sortArgs
	} else {
		req.Params.Filter.SortParams = []report.GetReportsSort{
			{
				DateSort: mo.Some[ordertype.OrderType](ordertype.ASC),
			},
		}
	}

	return req, nil
}
