package report

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type Endpoints struct {
	
	// CreateReport create report and return it
	// 	catchable errors:
	// 		ErrValidationError as target
	// Target errors catch by:
	// 		errors.Is(err, ErrValidationError)
	CreateReport endpoint.Endpoint[*reqresp.CreateReportReq, *reqresp.CreateReportResp]

	// DeleteReport delete report by id
	// 	catchable errors:
	// 		ErrReportIDNotValid
	// 		ErrReportNotFound
	DeleteReport endpoint.Endpoint[*reqresp.DeleteReportReq, *reqresp.DeleteReportResp]

	// UpdateReport update reports by id and not nil optionals
	// Name, Text, Implemtner fields can't be empty
	// 	catchable errors:
	// 		ErrReportIDNotValid
	// 		ErrReportNotFound
	// 		ErrValidationError as target
	// Target errors catch by:
	// 		errors.Is(err, ErrValidationError)
	UpdateReport endpoint.Endpoint[*reqresp.UpdateReportReq, *reqresp.UpdateReportResp]

	// GetReport return report by id
	// 	catchable errors:
	// 		ErrReportIDNotValid
	// 		ErrReportNotFound
	GetReport endpoint.Endpoint[*reqresp.GetReportReq, *reqresp.GetReportResp]

	// GetReports return reports acording to filters
	// 
	// don't have catchable errors
	GetReports endpoint.Endpoint[*reqresp.GetReportsReq, *reqresp.GetReportsResp]

	// CountReport count report according to filter and return count
	// 
	// don't have catchable errors
	CountReports endpoint.Endpoint[*reqresp.CountReportsReq, *reqresp.CountReportsResp]
}

type EndpointsOptions func(e *Endpoints)

func MakeEndpoints(
	s reports.Service,
	opts ...EndpointsOptions,
) Endpoints {
	e := Endpoints{
		CreateReport: makeCreateEndpoint(s),
		DeleteReport: makeDeleteEndpoint(s),
		UpdateReport: makeUpdateEndpoint(s),
		GetReport: makeGetEndpoint(s),
		GetReports: makeGetReportsEndpoint(s),
		CountReports: makeCountReportsEndpoint(s),
	}

	for _, opt := range opts {
		opt(&e)
	}

	return e
}

func makeCreateEndpoint(
	s	reports.Service,
) endpoint.Endpoint[*reqresp.CreateReportReq, *reqresp.CreateReportResp] {
	return func(
		ctx			context.Context,
		request		*reqresp.CreateReportReq,
	) (responce *reqresp.CreateReportResp, err error) {
		created, err := s.CreateReport(
			ctx,
			request.Report,
		)
		if err != nil {
			return nil, err
		}

		return &reqresp.CreateReportResp{
			Report: created,
		}, nil
	}
}

func makeDeleteEndpoint(
	s	reports.Service,
) endpoint.Endpoint[*reqresp.DeleteReportReq, *reqresp.DeleteReportResp] {
	return func(
		ctx			context.Context,
		request		*reqresp.DeleteReportReq,
	) (responce *reqresp.DeleteReportResp, err error) {
		err = s.DeleteReport(
			ctx,
			request.ID,
		)
		if err != nil {
			return nil, err
		}

		return &reqresp.DeleteReportResp{}, nil
	}
}

func makeUpdateEndpoint(
	s	reports.Service,
) endpoint.Endpoint[*reqresp.UpdateReportReq, *reqresp.UpdateReportResp] {
	return func(
		ctx			context.Context,
		request		*reqresp.UpdateReportReq,
	) (responce *reqresp.UpdateReportResp, err error) {
		updated, err := s.UpdateReport(
			ctx,
			request.ID,
			request.Params,
		)
		if err != nil {
			return nil, err
		}

		return &reqresp.UpdateReportResp{
			Report: updated,
		}, nil
	}
}

func makeGetEndpoint(
	s	reports.Service,
) endpoint.Endpoint[*reqresp.GetReportReq, *reqresp.GetReportResp] {
	return func(
		ctx			context.Context,
		request		*reqresp.GetReportReq,
	) (responce *reqresp.GetReportResp, err error) {
		get, err := s.GetReport(ctx, request.ID)
		if err != nil {
			return nil, err
		}
		
		return &reqresp.GetReportResp{
			Report: get,
		}, nil
	}
}

func makeGetReportsEndpoint(
	s	reports.Service,
) endpoint.Endpoint[*reqresp.GetReportsReq, *reqresp.GetReportsResp] {
	return func(
		ctx			context.Context,
		request		*reqresp.GetReportsReq,
	) (responce *reqresp.GetReportsResp, err error) {
		reports, err := s.GetReports(
			ctx,
			request.Params,
		)
		if err != nil {
			return nil, err
		}

		return &reqresp.GetReportsResp{
			Reports: reports,
		}, nil
	}
}

func makeCountReportsEndpoint(
	s	reports.Service,
) endpoint.Endpoint[*reqresp.CountReportsReq, *reqresp.CountReportsResp] {
	return func(
		ctx			context.Context,
		request		*reqresp.CountReportsReq,
	) (responce *reqresp.CountReportsResp, err error) {
		count, err := s.CountReports(ctx, request.Params)
		if err != nil {
			return nil, err
		}

		return &reqresp.CountReportsResp{
			Count: count,
		}, nil
	}
}