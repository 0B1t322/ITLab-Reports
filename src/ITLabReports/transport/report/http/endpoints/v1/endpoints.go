package endpoints

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
)

type Endpoints struct {
	GetReport             endpoint.Endpoint[*dto.GetReportReq, *dto.GetReportResp]
	GetReportsForEmployee endpoint.Endpoint[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp]
	CreateReport          endpoint.Endpoint[*dto.CreateReportReq, *dto.CreateReportResp]
	GetReports            endpoint.Endpoint[*dto.GetReportsReq, *dto.GetReportsResp]
}

func NewEndpoints(
	e report.Endpoints,
) Endpoints {
	return Endpoints{
		GetReport:             makeGetReport(e),
		GetReportsForEmployee: makeGetReportsForEmployeeEndpoint(e),
		CreateReport:          makeCreateReportEndpoint(e),
		GetReports:            makeGetReportsEndpoint(e),
	}
}

func makeGetReport(e report.Endpoints) endpoint.Endpoint[*dto.GetReportReq, *dto.GetReportResp] {
	return func(ctx context.Context, req *dto.GetReportReq) (*dto.GetReportResp, error) {
		resp, err := e.GetReport(
			ctx,
			req.ToEndopointReq(),
		)
		if err != nil {
			return nil, err
		}

		dtoResp := dto.GetReportRespFrom(resp)

		return dtoResp, nil
	}
}

func makeGetReportsForEmployeeEndpoint(
	e report.Endpoints,
) endpoint.Endpoint[*dto.GetReportsForEmployeeReq, *dto.GetReportsResp] {
	return func(
		ctx context.Context,
		request *dto.GetReportsForEmployeeReq,
	) (responce *dto.GetReportsResp, err error) {
		resp, err := e.GetReports(
			ctx,
			request.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		dtoResp := dto.GetReportsRespFrom(resp)

		return &dtoResp, nil
	}
}

func makeCreateReportEndpoint(
	e report.Endpoints,
) endpoint.Endpoint[*dto.CreateReportReq, *dto.CreateReportResp] {
	return func(
		ctx context.Context,
		request *dto.CreateReportReq,
	) (responce *dto.CreateReportResp, err error) {
		resp, err := e.CreateReport(
			ctx,
			request.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		dtoResp := dto.CreateReportRespFrom(resp)
		return dtoResp, nil
	}
}

func makeGetReportsEndpoint(
	e report.Endpoints,
) endpoint.Endpoint[*dto.GetReportsReq, *dto.GetReportsResp] {
	return func(
		ctx context.Context,
		request *dto.GetReportsReq,
	) (responce *dto.GetReportsResp, err error) {
		resp, err := e.GetReports(
			ctx,
			request.ToEndpointReq(),
		)
		if err != nil {
			return nil, err
		}

		dtoResp := dto.GetReportsRespFrom(resp)

		return &dtoResp, nil
	}
}
