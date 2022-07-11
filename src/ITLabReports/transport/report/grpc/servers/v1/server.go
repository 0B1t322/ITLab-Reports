package servers

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/grpc/handlers/v1"
	im "github.com/RTUITLab/ITLab-Reports/transport/report/grpc/middlewares"
	internalMiddlewares "github.com/RTUITLab/ITLab-Reports/transport/report/middlewares"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	pb.UnimplementedReportsServer

	// Handlers
	getReport            gt.Handler
	getReportImplementer gt.Handler
	getReports           gt.Handler
	getPaginatedReports  gt.Handler
}

type serverOptions struct {
	auther                  middlewares.Auther
	approvedReportsIdGetter internalMiddlewares.ApprovedReportsIdsGetter
}

type ServerOptions func(s *serverOptions)

func WithAuther(a middlewares.Auther) ServerOptions {
	return func(s *serverOptions) {
		s.auther = a
	}
}

func WithApprovedreportsIdGetter(a internalMiddlewares.ApprovedReportsIdsGetter) ServerOptions {
	return func(s *serverOptions) {
		s.approvedReportsIdGetter = a
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
		getReport:            handlers.GetReportHandler(e),
		getReportImplementer: handlers.GetReportImplementerHandler(e),
		getReports:           handlers.GetReportsHandler(e),
		getPaginatedReports:  handlers.GetReportsPaginatedHandler(e),
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

	ends.GetReports.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
		middlewares.RunMiddlewareIfAllFail(
			middlewares.SetReporterAndImplementer[*dto.GetReportsReq, *dto.GetReportsResp](),
			middlewares.IsAdmin[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
			middlewares.IsSuperAdmin[*dto.GetReportsReq, *dto.GetReportsResp](opt.auther),
		),
	)

	ends.GetReportsPaginated.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetReportsPaginatedReq, *dto.GetReportsPaginatedResp](opt.auther),
		internalMiddlewares.SetApprovedStateReportsIds[*dto.GetReportsPaginatedReq, *dto.GetReportsPaginatedResp](opt.approvedReportsIdGetter, opt.auther),
		middlewares.RunMiddlewareIfAllFail(
			middlewares.SetReporterAndImplementer[*dto.GetReportsPaginatedReq, *dto.GetReportsPaginatedResp](),
			middlewares.IsAdmin[*dto.GetReportsPaginatedReq, *dto.GetReportsPaginatedResp](opt.auther),
			middlewares.IsSuperAdmin[*dto.GetReportsPaginatedReq, *dto.GetReportsPaginatedResp](opt.auther),
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

// Return reports list without pagaination
func (s *gRPCServer) GetReports(
	ctx context.Context,
	req *pb.GetReportsReq,
) (*pb.GetReportsResp, error) {
	_, resp, err := s.getReports.ServeGRPC(ctx, req)
	if err != nil {
		return nil, HandleErrors(err)
	}

	return resp.(*pb.GetReportsResp), nil
}

// Return reports list with pagaination
func (s *gRPCServer) GetReportsPaginated(
	ctx context.Context,
	req *pb.GetReportsPaginatedReq,
) (*pb.GetReportsPaginatedResp, error) {
	_, resp, err := s.getPaginatedReports.ServeGRPC(ctx, req)
	if err != nil {
		return nil, HandleErrors(err)
	}

	return resp.(*pb.GetReportsPaginatedResp), nil
}
