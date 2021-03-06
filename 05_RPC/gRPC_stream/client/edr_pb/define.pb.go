// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: define.proto

package edr_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// HelloRequest 请求结构
type ProcessInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	PidCur   int64  `protobuf:"varint,2,opt,name=PidCur,proto3" json:"PidCur,omitempty"`
	PidChild int64  `protobuf:"varint,3,opt,name=PidChild,proto3" json:"PidChild,omitempty"`
}

func (x *ProcessInfo) Reset() {
	*x = ProcessInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_define_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessInfo) ProtoMessage() {}

func (x *ProcessInfo) ProtoReflect() protoreflect.Message {
	mi := &file_define_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessInfo.ProtoReflect.Descriptor instead.
func (*ProcessInfo) Descriptor() ([]byte, []int) {
	return file_define_proto_rawDescGZIP(), []int{0}
}

func (x *ProcessInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProcessInfo) GetPidCur() int64 {
	if x != nil {
		return x.PidCur
	}
	return 0
}

func (x *ProcessInfo) GetPidChild() int64 {
	if x != nil {
		return x.PidChild
	}
	return 0
}

// HelloResponse 响应结构
type ServerAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ack int32 `protobuf:"varint,1,opt,name=Ack,proto3" json:"Ack,omitempty"`
}

func (x *ServerAck) Reset() {
	*x = ServerAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_define_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerAck) ProtoMessage() {}

func (x *ServerAck) ProtoReflect() protoreflect.Message {
	mi := &file_define_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerAck.ProtoReflect.Descriptor instead.
func (*ServerAck) Descriptor() ([]byte, []int) {
	return file_define_proto_rawDescGZIP(), []int{1}
}

func (x *ServerAck) GetAck() int32 {
	if x != nil {
		return x.Ack
	}
	return 0
}

var File_define_proto protoreflect.FileDescriptor

var file_define_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x55,
	0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a,
	0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x69, 0x64, 0x43, 0x75, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x50, 0x69, 0x64, 0x43, 0x75, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x69, 0x64,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x50, 0x69, 0x64,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x22, 0x1d, 0x0a, 0x09, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x41,
	0x63, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x41, 0x63, 0x6b, 0x32, 0x41, 0x0a, 0x0c, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0c, 0x2e, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0a, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x41, 0x63, 0x6b, 0x22, 0x00, 0x28, 0x01, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x2f, 0x3b, 0x65, 0x64,
	0x72, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_define_proto_rawDescOnce sync.Once
	file_define_proto_rawDescData = file_define_proto_rawDesc
)

func file_define_proto_rawDescGZIP() []byte {
	file_define_proto_rawDescOnce.Do(func() {
		file_define_proto_rawDescData = protoimpl.X.CompressGZIP(file_define_proto_rawDescData)
	})
	return file_define_proto_rawDescData
}

var file_define_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_define_proto_goTypes = []interface{}{
	(*ProcessInfo)(nil), // 0: ProcessInfo
	(*ServerAck)(nil),   // 1: ServerAck
}
var file_define_proto_depIdxs = []int32{
	0, // 0: RemoteServer.UploadProcessInfo:input_type -> ProcessInfo
	1, // 1: RemoteServer.UploadProcessInfo:output_type -> ServerAck
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_define_proto_init() }
func file_define_proto_init() {
	if File_define_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_define_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_define_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerAck); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_define_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_define_proto_goTypes,
		DependencyIndexes: file_define_proto_depIdxs,
		MessageInfos:      file_define_proto_msgTypes,
	}.Build()
	File_define_proto = out.File
	file_define_proto_rawDesc = nil
	file_define_proto_goTypes = nil
	file_define_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RemoteServerClient is the client API for RemoteServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RemoteServerClient interface {
	// 定义UploadProcessInfo方法
	UploadProcessInfo(ctx context.Context, opts ...grpc.CallOption) (RemoteServer_UploadProcessInfoClient, error)
}

type remoteServerClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoteServerClient(cc grpc.ClientConnInterface) RemoteServerClient {
	return &remoteServerClient{cc}
}

func (c *remoteServerClient) UploadProcessInfo(ctx context.Context, opts ...grpc.CallOption) (RemoteServer_UploadProcessInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RemoteServer_serviceDesc.Streams[0], "/RemoteServer/UploadProcessInfo", opts...)
	if err != nil {
		return nil, err
	}
	x := &remoteServerUploadProcessInfoClient{stream}
	return x, nil
}

type RemoteServer_UploadProcessInfoClient interface {
	Send(*ProcessInfo) error
	CloseAndRecv() (*ServerAck, error)
	grpc.ClientStream
}

type remoteServerUploadProcessInfoClient struct {
	grpc.ClientStream
}

func (x *remoteServerUploadProcessInfoClient) Send(m *ProcessInfo) error {
	return x.ClientStream.SendMsg(m)
}

func (x *remoteServerUploadProcessInfoClient) CloseAndRecv() (*ServerAck, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ServerAck)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RemoteServerServer is the server API for RemoteServer service.
type RemoteServerServer interface {
	// 定义UploadProcessInfo方法
	UploadProcessInfo(RemoteServer_UploadProcessInfoServer) error
}

// UnimplementedRemoteServerServer can be embedded to have forward compatible implementations.
type UnimplementedRemoteServerServer struct {
}

func (*UnimplementedRemoteServerServer) UploadProcessInfo(RemoteServer_UploadProcessInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadProcessInfo not implemented")
}

func RegisterRemoteServerServer(s *grpc.Server, srv RemoteServerServer) {
	s.RegisterService(&_RemoteServer_serviceDesc, srv)
}

func _RemoteServer_UploadProcessInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RemoteServerServer).UploadProcessInfo(&remoteServerUploadProcessInfoServer{stream})
}

type RemoteServer_UploadProcessInfoServer interface {
	SendAndClose(*ServerAck) error
	Recv() (*ProcessInfo, error)
	grpc.ServerStream
}

type remoteServerUploadProcessInfoServer struct {
	grpc.ServerStream
}

func (x *remoteServerUploadProcessInfoServer) SendAndClose(m *ServerAck) error {
	return x.ServerStream.SendMsg(m)
}

func (x *remoteServerUploadProcessInfoServer) Recv() (*ProcessInfo, error) {
	m := new(ProcessInfo)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _RemoteServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "RemoteServer",
	HandlerType: (*RemoteServerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadProcessInfo",
			Handler:       _RemoteServer_UploadProcessInfo_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "define.proto",
}
