// Code generated by protoc-gen-go.
// source: auth/proto/auth.proto
// DO NOT EDIT!

/*
Package auth is a generated protocol buffer package.

It is generated from these files:
	auth/proto/auth.proto

It has these top-level messages:
	Credentials
	Token
	Verification
*/
package auth

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the credentials data.
type Credentials struct {
	Email    string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *Credentials) Reset()                    { *m = Credentials{} }
func (m *Credentials) String() string            { return proto.CompactTextString(m) }
func (*Credentials) ProtoMessage()               {}
func (*Credentials) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The request/response message containing the jwt data
type Token struct {
	Jwt string `protobuf:"bytes,1,opt,name=jwt" json:"jwt,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// The response message containing the jwt data
type Verification struct {
	IsValid bool `protobuf:"varint,1,opt,name=isValid" json:"isValid,omitempty"`
}

func (m *Verification) Reset()                    { *m = Verification{} }
func (m *Verification) String() string            { return proto.CompactTextString(m) }
func (*Verification) ProtoMessage()               {}
func (*Verification) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*Credentials)(nil), "auth.Credentials")
	proto.RegisterType((*Token)(nil), "auth.Token")
	proto.RegisterType((*Verification)(nil), "auth.Verification")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for AuthService service

type AuthServiceClient interface {
	Authenticate(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Token, error)
	Verify(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Verification, error)
}

type authServiceClient struct {
	cc *grpc.ClientConn
}

func NewAuthServiceClient(cc *grpc.ClientConn) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Authenticate(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := grpc.Invoke(ctx, "/auth.AuthService/Authenticate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Verify(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Verification, error) {
	out := new(Verification)
	err := grpc.Invoke(ctx, "/auth.AuthService/Verify", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AuthService service

type AuthServiceServer interface {
	Authenticate(context.Context, *Credentials) (*Token, error)
	Verify(context.Context, *Token) (*Verification, error)
}

func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	s.RegisterService(&_AuthService_serviceDesc, srv)
}

func _AuthService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Authenticate(ctx, req.(*Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Verify(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _AuthService_Authenticate_Handler,
		},
		{
			MethodName: "Verify",
			Handler:    _AuthService_Verify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("auth/proto/auth.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x4c, 0x90, 0x41, 0x4e, 0x85, 0x30,
	0x10, 0x40, 0xff, 0x57, 0x41, 0x1c, 0x58, 0xe8, 0x44, 0x13, 0x64, 0x65, 0xba, 0x22, 0x31, 0x01,
	0xa3, 0x07, 0x30, 0xc6, 0x1b, 0xa0, 0x61, 0x5f, 0x61, 0x0c, 0xa3, 0x95, 0x92, 0xb6, 0x48, 0xbc,
	0xbd, 0x69, 0xd1, 0x1f, 0x76, 0xef, 0x75, 0xda, 0xcc, 0x4b, 0xe1, 0x4a, 0xce, 0x6e, 0xa8, 0x27,
	0xa3, 0x9d, 0xae, 0x3d, 0x56, 0x01, 0xf1, 0xc4, 0xb3, 0x78, 0x84, 0xf4, 0xd9, 0x50, 0x4f, 0xa3,
	0x63, 0xa9, 0x2c, 0x5e, 0x42, 0x44, 0x5f, 0x92, 0x55, 0xbe, 0xbf, 0xd9, 0x97, 0x67, 0xcd, 0x2a,
	0x58, 0x40, 0x32, 0x49, 0x6b, 0x17, 0x6d, 0xfa, 0xfc, 0x28, 0x0c, 0x0e, 0x2e, 0xae, 0x21, 0x7a,
	0xd5, 0x9f, 0x34, 0xe2, 0x39, 0x1c, 0x7f, 0x2c, 0xee, 0xef, 0xa1, 0x47, 0x51, 0x42, 0xd6, 0x92,
	0xe1, 0x77, 0xee, 0xa4, 0x63, 0x3d, 0x62, 0x0e, 0xa7, 0x6c, 0x5b, 0xa9, 0xb8, 0x0f, 0xb7, 0x92,
	0xe6, 0x5f, 0xef, 0x15, 0xa4, 0x4f, 0xb3, 0x1b, 0x5e, 0xc8, 0x7c, 0x73, 0x47, 0x78, 0x07, 0x99,
	0x57, 0x1f, 0xd5, 0x49, 0x47, 0x78, 0x51, 0x85, 0xee, 0x4d, 0x68, 0x91, 0xae, 0x47, 0x61, 0xb5,
	0xd8, 0xe1, 0x2d, 0xc4, 0x61, 0xd5, 0x0f, 0x6e, 0x07, 0x05, 0xae, 0xb2, 0xad, 0x10, 0xbb, 0xb7,
	0x38, 0x7c, 0xc0, 0xc3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x83, 0x61, 0xeb, 0x73, 0x19, 0x01,
	0x00, 0x00,
}
