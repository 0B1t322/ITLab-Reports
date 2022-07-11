package endpoints

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/utils"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
)

type GetReport = endpoint.Endpoint[*dto.GetReportReq, *dto.GetReportResp]
type GetReportImplementer = endpoint.Endpoint[*dto.GetReportImplementerReq, *dto.GetReportImplementerResp]
type GetReports = endpoint.Endpoint[*dto.GetReportsReq, *dto.GetReportsResp]
type GetReportsPaginated = endpoint.Endpoint[*dto.GetReportsPaginatedReq, *dto.GetReportsPaginatedResp]

type Endpoints struct {
	GetReport            GetReport
	GetReportImplementer GetReportImplementer
	GetReports           GetReports
	GetReportsPaginated  GetReportsPaginated
}

func NewEndpoint(
	e report.Endpoints,
) Endpoints {
	return Endpoints{
		GetReport:            makeGetReportEndpoint(e),
		GetReportImplementer: makeGetReportImplementerEndpoint(e),
		GetReports:           makeGetReports(e),
		GetReportsPaginated:  makeGetReportsPaginated(e),
	}
}

func makeGetReportEndpoint(
	e report.Endpoints,
) GetReport {
	return func(
		ctx context.Context,
		req *dto.GetReportReq,
	) (responce *dto.GetReportResp, err error) {
		resp, err := e.GetReport(
			ctx,
			req.ToEndpointReq(),
		)
		if err == reports.ErrReportNotFound || err == reports.ErrReportIDNotValid {
			return &dto.GetReportResp{
				Result: &pb.GetReportResp_Error{
					Error: pb.ReportsServiceErrors_REPORT_NOT_FOUND,
				},
			}, nil
		} else if err != nil {
			return nil, err
		}
		return &dto.GetReportResp{
			Result: &pb.GetReportResp_Report{
				Report: utils.ReportToPBType(resp.Report),
			},
		}, nil
	}
}

func makeGetReportImplementerEndpoint(
	e report.Endpoints,
) GetReportImplementer {
	return func(
		ctx context.Context,
		req *dto.GetReportImplementerReq,
	) (responce *dto.GetReportImplementerResp, err error) {
		resp, err := e.GetReport(
			ctx,
			req.ToEndpointReq(),
		)

		if err == reports.ErrReportNotFound || err == reports.ErrReportIDNotValid {
			return &dto.GetReportImplementerResp{
				GetReportImplementerResp: &pb.GetReportImplementerResp{
					Result: &pb.GetReportImplementerResp_Error{
						Error: pb.ReportsServiceErrors_REPORT_NOT_FOUND,
					},
				},
				Reporter: "",
			}, nil
		} else if err != nil {
			return nil, err
		}

		return &dto.GetReportImplementerResp{
			GetReportImplementerResp: &pb.GetReportImplementerResp{
				Result: &pb.GetReportImplementerResp_Implementer{
					Implementer: resp.GetImplementer(),
				},
			},
			Reporter: resp.GetReporter(),
		}, nil
	}
}

func makeGetReports(
	e report.Endpoints,
) GetReports {
	return func(
		ctx context.Context,
		req *dto.GetReportsReq,
	) (responce *dto.GetReportsResp, err error) {
		resp, err := e.GetReports(
			ctx,
			req.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}
		return dto.GetReportsRespFrom(resp), nil
	}
}

func makeGetReportsPaginated(
	e report.Endpoints,
) GetReportsPaginated {
	return func(
		ctx context.Context,
		req *dto.GetReportsPaginatedReq,
	) (responce *dto.GetReportsPaginatedResp, err error) {
		resp, err := e.GetReports(
			ctx,
			req.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		countReport, err := e.CountReports(
			ctx,
			&reqresp.CountReportsReq{
				Params: &req.Params.Filter.GetReportsFilterFieldsWithOrAnd,
			},
		)
		responce = &dto.GetReportsPaginatedResp{}
		{
			responce.Offset = 0
			responce.Limit = 0
			responce.Count = int32(len(resp.Reports))
			responce.TotalResult = int32(countReport.Count)
			if req.Params.Limit.HasValue() {
				responce.Limit = int32(req.Params.Limit.MustGetValue())
			}

			if req.Params.Offset.HasValue() {
				responce.Offset = int32(req.Params.Offset.MustGetValue())
			}

			responce.HasMore = false
			if responce.Limit != 0 && responce.TotalResult - responce.Offset -  responce.Limit > 0 {
				responce.HasMore = true
			}

			for _, r := range resp.Reports {
				responce.Reports = append(responce.Reports, utils.ReportToPBType(r))
			}
		}
		return responce, nil
	}
}
