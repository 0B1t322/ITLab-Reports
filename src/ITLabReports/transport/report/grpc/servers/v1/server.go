package servers

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/handlers/v1"
	im "github.com/RTUITLab/ITLab-Reports/transport/report/grpc/middlewares"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	pb.UnimplementedReportsServer

	// Handlers
	getReport            gt.Handler
	getReportImplementer gt.Handler
}

type serverOptions struct {
	auther       middlewares.Auther
}

type ServerOptions func(s *serverOptions)

func WithAuther(a middlewares.Auther) ServerOptions {
	return func(s *serverOptions) {
		s.auther = a
	}
}

func MergeServerOptions(opts ...ServerOptions) *serverOptions {
	s := &serverOptions{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func NewServer(
	ends report.Endpoints,
	opts ...ServerOptions,
) pb.ReportsServer {
	e := BuildMiddlewares(
		endpoints.NewEndpoint(ends),
		MergeServerOptions(opts...),
	)

	return &gRPCServer{
		getReport: handlers.GetReportHandler(e),
		getReportImplementer: handlers.GetReportImplementerHandler(e),
	}
}

func BuildMiddlewares(
	ends endpoints.Endpoints,
	opt *serverOptions,
) endpoints.Endpoints {
	// Add to check error in oneof
	ends.GetReport.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
		middlewares.RunMiddlewareIfAllFail(
			im.CheckUserIsReporterOrImplementerIfNotError[*dto.GetReportReq, *dto.GetReportResp](),
			middlewares.IsSuperAdmin[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
			middlewares.IsAdmin[*dto.GetReportReq, *dto.GetReportResp](opt.auther),
		),
	)

	ends.GetReportImplementer.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportImplementerReq, *dto.GetReportImplementerResp](opt.auther),
		middlewares.RunMiddlewareIfAllFail(
			im.CheckUserIsReporterOrImplementerIfNotError[*dto.GetReportImplementerReq, *dto.GetReportImplementerResp](),
			middlewares.IsSuperAdmin[*dto.GetReportImplementerReq, *dto.GetReportImplementerResp](opt.auther),
			middlewares.IsAdmin[*dto.GetReportImplementerReq, *dto.GetReportImplementerResp](opt.auther),
		),
	)

	return ends
}

// Return report implementer by id
// If report not found return REPORT_NOT_FOUND error
func (s *gRPCServer) GetReportImplementer(
	ctx context.Context,
	req *pb.GetReportImplementerReq,
) (*pb.GetReportImplementerResp, error) {
	_, resp, err := s.getReportImplementer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, HandleErrors(err)
	}

	return resp.(*pb.GetReportImplementerResp), nil
}

// Return report by id
// If report not found return REPORT_NOT_FOUND error
func (s *gRPCServer) GetReport(
	ctx context.Context,
	req *pb.GetReportReq,
) (*pb.GetReportResp, error) {
	_, resp, err := s.getReport.ServeGRPC(ctx, req)
	if err != nil {
		return nil, HandleErrors(err)
	}

	return resp.(*pb.GetReportResp), nil
}