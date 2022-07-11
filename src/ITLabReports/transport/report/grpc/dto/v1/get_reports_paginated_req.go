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
)

type GetReportsPaginatedReq struct {
	Params        *report.GetReportsParams
	ApprovedState pb.GetReportsPaginatedReq_ApprovedState
}

func (g *GetReportsPaginatedReq) ToEndpointReq() *reqresp.GetReportsReq {
	return &reqresp.GetReportsReq{
		Params:        g.Params,
	}
}

func (g *GetReportsPaginatedReq) IsOnlyApprovedReports() bool {
	return g.ApprovedState == pb.GetReportsPaginatedReq_APPROVED
}

func (g *GetReportsPaginatedReq) SetOnlyApprovedReports(ids ...string) {
	g.Params.Filter.ReportsId = &filter.FilterField[[]string]{
		Value:     ids,
		Operation: filter.IN,
	}
}

func (g *GetReportsPaginatedReq) IsOnlyNotApprovedReports() bool {
	return g.ApprovedState == pb.GetReportsPaginatedReq_NOT_APPROVED
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

	// Limit
	if limit := grpcReq.GetLimit(); limit >= 1 {
		req.Params.Limit.SetValue(int64(limit))
	}

	// Offset
	if offset := grpcReq.GetOffset(); offset >= 0 {
		req.Params.Offset.SetValue(int64(offset))
	}

	// Decode state
	if state := grpcReq.ApprovedState; state != nil {
		req.ApprovedState = *state
	} else {
		req.ApprovedState = pb.GetReportsPaginatedReq_ALL
	}

	// Decode dateBegin
	if grpcReq.DateBegin != nil {
		date := grpcReq.DateBegin.AsTime()
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
	if grpcReq.DateEnd != nil {
		date := grpcReq.DateEnd.AsTime()
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

	// Decode match
	if match := grpcReq.MatchParams; match != nil {
		if name := match.GetName(); name != "" {
			req.Params.Filter.Name = &filter.FilterField[string]{
				Operation: filter.LIKE,
				Value:     name,
			}
		}

		if date := match.GetDate(); date != nil {
			req.Params.Filter.Date = &filter.FilterField[string]{
				Operation: filter.EQ,
				Value:     date.AsTime().UTC().Format(time.RFC3339Nano),
			}
		}

		if implementer := match.GetImplementer(); implementer != "" {
			req.Params.Filter.Implementer = &filter.FilterField[string]{
				Operation: filter.EQ,
				Value:     implementer,
			}
		}

		if reporter := match.GetReporter(); reporter != "" {
			req.Params.Filter.Reporter = &filter.FilterField[string]{
				Operation: filter.EQ,
				Value:     reporter,
			}
		}
	}

	// Decode sort
	if sort := grpcReq.SortParams; sort != nil {
		if sort.Name != nil {
			req.Params.Filter.NameSort.SetValue(utils.OrderTypeFromGRPC(*sort.Name))
		}

		if sort.Date != nil {
			req.Params.Filter.DateSort.SetValue(utils.OrderTypeFromGRPC(*sort.Date))
		}

	} else {
		req.Params.Filter.DateSort.SetValue(ordertype.ASC)
	}

	return req, nil
}
