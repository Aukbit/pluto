// Code generated by protoc-gen-go.
// source: examples/user/proto/user.proto
// DO NOT EDIT!

/*
Package user is a generated protocol buffer package.

It is generated from these files:
	examples/user/proto/user.proto

It has these top-level messages:
	NewUser
	User
	Users
	Filter
	Credentials
	Verification
*/
package user

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

// The request message containing the new user data.
type NewUser struct {
	Name     string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Email    string `protobuf:"bytes,2,opt,name=email" json:"email,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
}

func (m *NewUser) Reset()                    { *m = NewUser{} }
func (m *NewUser) String() string            { return proto.CompactTextString(m) }
func (*NewUser) ProtoMessage()               {}
func (*NewUser) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The response message containing the user data
type User struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Email string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// The response message containing the a users list
type Users struct {
	Data []*User `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *Users) Reset()                    { *m = Users{} }
func (m *Users) String() string            { return proto.CompactTextString(m) }
func (*Users) ProtoMessage()               {}
func (*Users) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Users) GetData() []*User {
	if m != nil {
		return m.Data
	}
	return nil
}

// The response message containing the a users list
type Filter struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *Filter) Reset()                    { *m = Filter{} }
func (m *Filter) String() string            { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()               {}
func (*Filter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// The request message containing the user basic credentials
type Credentials struct {
	Email    string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *Credentials) Reset()                    { *m = Credentials{} }
func (m *Credentials) String() string            { return proto.CompactTextString(m) }
func (*Credentials) ProtoMessage()               {}
func (*Credentials) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// The response message containing the a users list
type Verification struct {
	IsValid bool `protobuf:"varint,1,opt,name=isValid" json:"isValid,omitempty"`
}

func (m *Verification) Reset()                    { *m = Verification{} }
func (m *Verification) String() string            { return proto.CompactTextString(m) }
func (*Verification) ProtoMessage()               {}
func (*Verification) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto.RegisterType((*NewUser)(nil), "user.NewUser")
	proto.RegisterType((*User)(nil), "user.User")
	proto.RegisterType((*Users)(nil), "user.Users")
	proto.RegisterType((*Filter)(nil), "user.Filter")
	proto.RegisterType((*Credentials)(nil), "user.Credentials")
	proto.RegisterType((*Verification)(nil), "user.Verification")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UserService service

type UserServiceClient interface {
	CreateUser(ctx context.Context, in *NewUser, opts ...grpc.CallOption) (*User, error)
	ReadUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	UpdateUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	DeleteUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	FilterUsers(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Users, error)
	VerifyUser(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Verification, error)
}

type userServiceClient struct {
	cc *grpc.ClientConn
}

func NewUserServiceClient(cc *grpc.ClientConn) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *NewUser, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/user.UserService/CreateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ReadUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/user.UserService/ReadUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/user.UserService/UpdateUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/user.UserService/DeleteUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) FilterUsers(ctx context.Context, in *Filter, opts ...grpc.CallOption) (*Users, error) {
	out := new(Users)
	err := grpc.Invoke(ctx, "/user.UserService/FilterUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) VerifyUser(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Verification, error) {
	out := new(Verification)
	err := grpc.Invoke(ctx, "/user.UserService/VerifyUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceServer interface {
	CreateUser(context.Context, *NewUser) (*User, error)
	ReadUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User) (*User, error)
	DeleteUser(context.Context, *User) (*User, error)
	FilterUsers(context.Context, *Filter) (*Users, error)
	VerifyUser(context.Context, *Credentials) (*Verification, error)
}

func RegisterUserServiceServer(s *grpc.Server, srv UserServiceServer) {
	s.RegisterService(&_UserService_serviceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUser)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*NewUser))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ReadUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ReadUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/ReadUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ReadUser(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/UpdateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUser(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_FilterUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Filter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).FilterUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/FilterUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).FilterUsers(ctx, req.(*Filter))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_VerifyUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).VerifyUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/VerifyUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).VerifyUser(ctx, req.(*Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "ReadUser",
			Handler:    _UserService_ReadUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserService_DeleteUser_Handler,
		},
		{
			MethodName: "FilterUsers",
			Handler:    _UserService_FilterUsers_Handler,
		},
		{
			MethodName: "VerifyUser",
			Handler:    _UserService_VerifyUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/user/proto/user.proto",
}

func init() { proto.RegisterFile("examples/user/proto/user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 328 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x92, 0x4f, 0x4b, 0xf3, 0x40,
	0x10, 0xc6, 0x9b, 0x34, 0xfd, 0xf3, 0x4e, 0xfa, 0x0a, 0x2e, 0x1e, 0x42, 0x90, 0x52, 0x16, 0xd1,
	0xea, 0xa1, 0x85, 0x8a, 0x67, 0x85, 0x8a, 0x47, 0x85, 0x48, 0x7b, 0x5f, 0xbb, 0x23, 0x2c, 0x6c,
	0x9b, 0xb0, 0xbb, 0x5a, 0xfd, 0x18, 0x7e, 0x63, 0xc9, 0x6c, 0xda, 0x86, 0x52, 0xf5, 0x36, 0x33,
	0xcf, 0xc3, 0x6f, 0x67, 0x1e, 0x16, 0xfa, 0xf8, 0x21, 0x96, 0x85, 0x46, 0x3b, 0x7e, 0xb3, 0x68,
	0xc6, 0x85, 0xc9, 0x5d, 0x4e, 0xe5, 0x88, 0x4a, 0x16, 0x95, 0x35, 0x7f, 0x82, 0xce, 0x23, 0xae,
	0x67, 0x16, 0x0d, 0x63, 0x10, 0xad, 0xc4, 0x12, 0x93, 0x60, 0x10, 0x0c, 0xff, 0x65, 0x54, 0xb3,
	0x13, 0x68, 0xe1, 0x52, 0x28, 0x9d, 0x84, 0x34, 0xf4, 0x0d, 0x4b, 0xa1, 0x5b, 0x08, 0x6b, 0xd7,
	0xb9, 0x91, 0x49, 0x93, 0x84, 0x6d, 0xcf, 0xef, 0x20, 0x22, 0xda, 0x11, 0x84, 0x4a, 0x56, 0xac,
	0x50, 0xc9, 0x2d, 0x3d, 0x3c, 0x44, 0x6f, 0xd6, 0xe8, 0xfc, 0x02, 0x5a, 0x25, 0xc1, 0xb2, 0x3e,
	0x44, 0x52, 0x38, 0x91, 0x04, 0x83, 0xe6, 0x30, 0x9e, 0xc0, 0x88, 0x96, 0x2f, 0xa5, 0x8c, 0xe6,
	0xfc, 0x14, 0xda, 0x0f, 0x4a, 0xbb, 0xc3, 0xab, 0xf3, 0x5b, 0x88, 0xa7, 0x06, 0x25, 0xae, 0x9c,
	0x12, 0xda, 0xee, 0xde, 0x0a, 0x7e, 0xba, 0x24, 0xdc, 0xbb, 0x64, 0x08, 0xbd, 0x39, 0x1a, 0xf5,
	0xaa, 0x16, 0xc2, 0xa9, 0x7c, 0xc5, 0x12, 0xe8, 0x28, 0x3b, 0x17, 0xba, 0x3a, 0xab, 0x9b, 0x6d,
	0xda, 0xc9, 0x57, 0x08, 0x71, 0xb9, 0xd7, 0x33, 0x9a, 0x77, 0xb5, 0x40, 0x76, 0x09, 0x30, 0x35,
	0x28, 0x1c, 0x52, 0x12, 0xff, 0xfd, 0xe2, 0x55, 0xcc, 0x69, 0xed, 0x0e, 0xde, 0x60, 0x67, 0xd0,
	0xcd, 0x50, 0x48, 0x32, 0xd6, 0x94, 0x3d, 0xd7, 0x39, 0xc0, 0xac, 0x90, 0x1b, 0xe0, 0xaf, 0xbe,
	0x7b, 0xd4, 0xf8, 0xa7, 0xef, 0x0a, 0x62, 0x9f, 0x9c, 0x0f, 0xba, 0xe7, 0x45, 0x3f, 0x4a, 0xe3,
	0x9d, 0xd5, 0xf2, 0x06, 0xbb, 0x01, 0xa0, 0x18, 0x3e, 0x89, 0x79, 0xec, 0xc5, 0x5a, 0xb2, 0x29,
	0xf3, 0xa3, 0x7a, 0x56, 0xbc, 0xf1, 0xd2, 0xa6, 0x5f, 0x76, 0xfd, 0x1d, 0x00, 0x00, 0xff, 0xff,
	0x72, 0x6f, 0x14, 0xbd, 0x87, 0x02, 0x00, 0x00,
}
