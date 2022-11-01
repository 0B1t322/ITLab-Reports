package servers

import (
	"context"

	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/grpc/handlers/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	pb "github.com/RTUITLab/ITLab/proto/reports/v1"
	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	pb.UnimplementedDraftsServer

	// Handlers
	getDraft           gt.Handler
	getDrafts          gt.Handler
	getDraftsPaginated gt.Handler
}

type serverOptions struct {
	auther middlewares.Auther
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
) pb.DraftsServer {
	e := BuildMiddlewares(
		endpoints.NewEndpoint(ends),
		MergeServerOptions(opts...),
	)

	return &gRPCServer{
		getDraft:           handlers.GetDraftHandler(e),
		getDrafts:          handlers.GetDraftsHandler(e),
		getDraftsPaginated: handlers.GetDraftsPaginatedHandler(e),
	}
}

func BuildMiddlewares(
	ends endpoints.Endpoints,
	opt *serverOptions,
) endpoints.Endpoints {
	// Add to check error in oneof
	ends.GetDraft.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetDraftReq, *dto.GetDraftResp](opt.auther),
		middlewares.CheckIsError[*dto.GetDraftReq, *dto.GetDraftResp](),
		middlewares.CheckUserIsReporter[*dto.GetDraftReq, *dto.GetDraftResp](),
	)

	ends.GetDrafts.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetDraftsReq, *dto.GetDraftsResp](opt.auther),
		middlewares.SetUserID[*dto.GetDraftsReq, *dto.GetDraftsResp](),
	)

	ends.GetDraftsPaginated.AddCustomMiddlewares(
		middlewares.Auth[*dto.GetDraftsPaginatedReq, *dto.GetDraftsPaginatedResp](opt.auther),
		middlewares.SetUserID[*dto.GetDraftsPaginatedReq, *dto.GetDraftsPaginatedResp](),
	)

	return ends
}

func (s *gRPCServer) GetDraft(ctx context.Context, req *pb.GetDraftReq) (*pb.GetDraftResp, error) {
	_, resp, err := s.getDraft.ServeGRPC(ctx, req)
	if err != nil {
		return nil, HandleErrors(err)
	}
	return resp.(*pb.GetDraftResp), nil
}

func (s *gRPCServer) GetDrafts(ctx context.Context, req *pb.GetDraftsReq) (*pb.GetDraftsResp, error) {
	_, resp, err := s.getDrafts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, HandleErrors(err)
	}
	return resp.(*pb.GetDraftsResp), nil
}

func (s *gRPCServer) GetDraftsPaginated(ctx context.Context, req *pb.GetDraftsPaginatedReq) (*pb.GetDraftsPaginatedResp, error) {
	_, resp, err := s.getDraftsPaginated.ServeGRPC(ctx, req)
	if err != nil {
		return nil, HandleErrors(err)
	}
	return resp.(*pb.GetDraftsPaginatedResp), nil
}
