// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: proto/cropper.proto

package cropper

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

// UrlCropperClient is the client API for UrlCropper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UrlCropperClient interface {
	CropUrl(ctx context.Context, in *CropRequest, opts ...grpc.CallOption) (*CroppedUrl, error)
}

type urlCropperClient struct {
	cc grpc.ClientConnInterface
}

func NewUrlCropperClient(cc grpc.ClientConnInterface) UrlCropperClient {
	return &urlCropperClient{cc}
}

func (c *urlCropperClient) CropUrl(ctx context.Context, in *CropRequest, opts ...grpc.CallOption) (*CroppedUrl, error) {
	out := new(CroppedUrl)
	err := c.cc.Invoke(ctx, "/rpc.UrlCropper/CropUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UrlCropperServer is the srv API for UrlCropper service.
// All implementations must embed UnimplementedUrlCropperServer
// for forward compatibility
type UrlCropperServer interface {
	CropUrl(context.Context, *CropRequest) (*CroppedUrl, error)
	mustEmbedUnimplementedUrlCropperServer()
}

// UnimplementedUrlCropperServer must be embedded to have forward compatible implementations.
type UnimplementedUrlCropperServer struct {
}

func (UnimplementedUrlCropperServer) CropUrl(context.Context, *CropRequest) (*CroppedUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CropUrl not implemented")
}
func (UnimplementedUrlCropperServer) mustEmbedUnimplementedUrlCropperServer() {}

// UnsafeUrlCropperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UrlCropperServer will
// result in compilation errors.
type UnsafeUrlCropperServer interface {
	mustEmbedUnimplementedUrlCropperServer()
}

func RegisterUrlCropperServer(s grpc.ServiceRegistrar, srv UrlCropperServer) {
	s.RegisterService(&UrlCropper_ServiceDesc, srv)
}

func _UrlCropper_CropUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlCropperServer).CropUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.UrlCropper/CropUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlCropperServer).CropUrl(ctx, req.(*CropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UrlCropper_ServiceDesc is the grpc.ServiceDesc for UrlCropper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UrlCropper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.UrlCropper",
	HandlerType: (*UrlCropperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CropUrl",
			Handler:    _UrlCropper_CropUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/cropper.proto",
}
