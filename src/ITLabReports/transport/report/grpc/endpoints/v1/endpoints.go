package endpoints

import (
	"github.com/sirupsen/logrus"
	"context"

	"github.com/RTUITLab/ITLab-Reports/pkg/endpoint"
	"github.com/RTUITLab/ITLab-Reports/service/reports"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/dto/v1"
	types "github.com/RTUITLab/ITLab/proto/reports/types"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetReport = endpoint.Endpoint[*dto.GetReportReq, *dto.GetReportResp]
type GetReportImplementer = endpoint.Endpoint[*dto.GetReportImplementerReq, *dto.GetReportImplementerResp]

type Endpoints struct {
	GetReport            GetReport
	GetReportImplementer GetReportImplementer
}

func NewEndpoint(
	e report.Endpoints,
) Endpoints {
	return Endpoints{
		GetReport: makeGetReportEndpoint(e),
		GetReportImplementer: makeGetReportImplementerEndpoint(e),
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
		logrus.Info("Get report endpoint")
		return &dto.GetReportResp{
			Result: &pb.GetReportResp_Report{
				Report: &types.Report{
					Id:   resp.Report.GetID(),
					Name: resp.Report.GetName(),
					Text: resp.Report.GetText(),
					Assignees: &types.Assignees{
						Reporter:    resp.Report.GetReporter(),
						Implementer: resp.Report.GetImplementer(),
					},
					Date: timestamppb.New(resp.Report.GetDate()),
				},
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
				Reporter:                 "",
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
			Reporter:                 resp.GetReporter(),
		}, nil
	}
}