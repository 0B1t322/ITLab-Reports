// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: rtuitlab/itlab/reports/v1/drafts.proto

package services

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DraftsClient is the client API for Drafts service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DraftsClient interface {
	// Return draft by id
	// If draft not found return DRAFT_NOT_FOUND error
	GetDraft(ctx context.Context, in *GetDraftReq, opts ...grpc.CallOption) (*GetDraftResp, error)
	// Return list of drafts
	GetDrafts(ctx context.Context, in *GetDraftsReq, opts ...grpc.CallOption) (*GetDraftsResp, error)
	// Return paginated list of drafts
	GetDraftsPaginated(ctx context.Context, in *GetDraftsPaginatedReq, opts ...grpc.CallOption) (*GetDraftsPaginatedResp, error)
}

type draftsClient struct {
	cc grpc.ClientConnInterface
}

func NewDraftsClient(cc grpc.ClientConnInterface) DraftsClient {
	return &draftsClient{cc}
}

func (c *draftsClient) GetDraft(ctx context.Context, in *GetDraftReq, opts ...grpc.CallOption) (*GetDraftResp, error) {
	out := new(GetDraftResp)
	err := c.cc.Invoke(ctx, "/rtuitlab.itlab.reports.v1.Drafts/GetDraft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *draftsClient) GetDrafts(ctx context.Context, in *GetDraftsReq, opts ...grpc.CallOption) (*GetDraftsResp, error) {
	out := new(GetDraftsResp)
	err := c.cc.Invoke(ctx, "/rtuitlab.itlab.reports.v1.Drafts/GetDrafts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *draftsClient) GetDraftsPaginated(ctx context.Context, in *GetDraftsPaginatedReq, opts ...grpc.CallOption) (*GetDraftsPaginatedResp, error) {
	out := new(GetDraftsPaginatedResp)
	err := c.cc.Invoke(ctx, "/rtuitlab.itlab.reports.v1.Drafts/GetDraftsPaginated", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DraftsServer is the server API for Drafts service.
// All implementations must embed UnimplementedDraftsServer
// for forward compatibility
type DraftsServer interface {
	// Return draft by id
	// If draft not found return DRAFT_NOT_FOUND error
	GetDraft(context.Context, *GetDraftReq) (*GetDraftResp, error)
	// Return list of drafts
	GetDrafts(context.Context, *GetDraftsReq) (*GetDraftsResp, error)
	// Return paginated list of drafts
	GetDraftsPaginated(context.Context, *GetDraftsPaginatedReq) (*GetDraftsPaginatedResp, error)
	mustEmbedUnimplementedDraftsServer()
}

// UnimplementedDraftsServer must be embedded to have forward compatible implementations.
type UnimplementedDraftsServer struct {
}

func (UnimplementedDraftsServer) GetDraft(context.Context, *GetDraftReq) (*GetDraftResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDraft not implemented")
}
func (UnimplementedDraftsServer) GetDrafts(context.Context, *GetDraftsReq) (*GetDraftsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDrafts not implemented")
}
func (UnimplementedDraftsServer) GetDraftsPaginated(context.Context, *GetDraftsPaginatedReq) (*GetDraftsPaginatedResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDraftsPaginated not implemented")
}
func (UnimplementedDraftsServer) mustEmbedUnimplementedDraftsServer() {}

// UnsafeDraftsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DraftsServer will
// result in compilation errors.
type UnsafeDraftsServer interface {
	mustEmbedUnimplementedDraftsServer()
}

func RegisterDraftsServer(s grpc.ServiceRegistrar, srv DraftsServer) {
	s.RegisterService(&Drafts_ServiceDesc, srv)
}

func _Drafts_GetDraft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDraftReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DraftsServer).GetDraft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rtuitlab.itlab.reports.v1.Drafts/GetDraft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DraftsServer).GetDraft(ctx, req.(*GetDraftReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Drafts_GetDrafts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDraftsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DraftsServer).GetDrafts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rtuitlab.itlab.reports.v1.Drafts/GetDrafts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DraftsServer).GetDrafts(ctx, req.(*GetDraftsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Drafts_GetDraftsPaginated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDraftsPaginatedReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DraftsServer).GetDraftsPaginated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rtuitlab.itlab.reports.v1.Drafts/GetDraftsPaginated",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DraftsServer).GetDraftsPaginated(ctx, req.(*GetDraftsPaginatedReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Drafts_ServiceDesc is the grpc.ServiceDesc for Drafts service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Drafts_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rtuitlab.itlab.reports.v1.Drafts",
	HandlerType: (*DraftsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDraft",
			Handler:    _Drafts_GetDraft_Handler,
		},
		{
			MethodName: "GetDrafts",
			Handler:    _Drafts_GetDrafts_Handler,
		},
		{
			MethodName: "GetDraftsPaginated",
			Handler:    _Drafts_GetDraftsPaginated_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rtuitlab/itlab/reports/v1/drafts.proto",
}
