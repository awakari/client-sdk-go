// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.12
// source: api/grpc/writer/service.proto

package writer

import (
	context "context"
	pb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
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

type SubmitMessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgs []*pb.CloudEvent `protobuf:"bytes,1,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (x *SubmitMessagesRequest) Reset() {
	*x = SubmitMessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_writer_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitMessagesRequest) ProtoMessage() {}

func (x *SubmitMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_writer_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitMessagesRequest.ProtoReflect.Descriptor instead.
func (*SubmitMessagesRequest) Descriptor() ([]byte, []int) {
	return file_api_grpc_writer_service_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitMessagesRequest) GetMsgs() []*pb.CloudEvent {
	if x != nil {
		return x.Msgs
	}
	return nil
}

type SubmitMessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AckCount uint32 `protobuf:"varint,1,opt,name=ackCount,proto3" json:"ackCount,omitempty"`
}

func (x *SubmitMessagesResponse) Reset() {
	*x = SubmitMessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_writer_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitMessagesResponse) ProtoMessage() {}

func (x *SubmitMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_writer_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitMessagesResponse.ProtoReflect.Descriptor instead.
func (*SubmitMessagesResponse) Descriptor() ([]byte, []int) {
	return file_api_grpc_writer_service_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitMessagesResponse) GetAckCount() uint32 {
	if x != nil {
		return x.AckCount
	}
	return 0
}

var File_api_grpc_writer_service_proto protoreflect.FileDescriptor

var file_api_grpc_writer_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x77, 0x72, 0x69, 0x74, 0x65,
	0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x1a,
	0x25, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x15, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x22, 0x0a, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x6d,
	0x73, 0x67, 0x73, 0x22, 0x34, 0x0a, 0x16, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x61, 0x63, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x61, 0x63, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x6e, 0x0a, 0x07, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x63, 0x0a, 0x0e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x25, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69,
	0x2e, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e,
	0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x2e, 0x53,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x11, 0x5a, 0x0f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_grpc_writer_service_proto_rawDescOnce sync.Once
	file_api_grpc_writer_service_proto_rawDescData = file_api_grpc_writer_service_proto_rawDesc
)

func file_api_grpc_writer_service_proto_rawDescGZIP() []byte {
	file_api_grpc_writer_service_proto_rawDescOnce.Do(func() {
		file_api_grpc_writer_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_grpc_writer_service_proto_rawDescData)
	})
	return file_api_grpc_writer_service_proto_rawDescData
}

var file_api_grpc_writer_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_grpc_writer_service_proto_goTypes = []interface{}{
	(*SubmitMessagesRequest)(nil),  // 0: awakari.writer.SubmitMessagesRequest
	(*SubmitMessagesResponse)(nil), // 1: awakari.writer.SubmitMessagesResponse
	(*pb.CloudEvent)(nil),          // 2: pb.CloudEvent
}
var file_api_grpc_writer_service_proto_depIdxs = []int32{
	2, // 0: awakari.writer.SubmitMessagesRequest.msgs:type_name -> pb.CloudEvent
	0, // 1: awakari.writer.Service.SubmitMessages:input_type -> awakari.writer.SubmitMessagesRequest
	1, // 2: awakari.writer.Service.SubmitMessages:output_type -> awakari.writer.SubmitMessagesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_grpc_writer_service_proto_init() }
func file_api_grpc_writer_service_proto_init() {
	if File_api_grpc_writer_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_grpc_writer_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitMessagesRequest); i {
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
		file_api_grpc_writer_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitMessagesResponse); i {
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
			RawDescriptor: file_api_grpc_writer_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_grpc_writer_service_proto_goTypes,
		DependencyIndexes: file_api_grpc_writer_service_proto_depIdxs,
		MessageInfos:      file_api_grpc_writer_service_proto_msgTypes,
	}.Build()
	File_api_grpc_writer_service_proto = out.File
	file_api_grpc_writer_service_proto_rawDesc = nil
	file_api_grpc_writer_service_proto_goTypes = nil
	file_api_grpc_writer_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceClient interface {
	SubmitMessages(ctx context.Context, opts ...grpc.CallOption) (Service_SubmitMessagesClient, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) SubmitMessages(ctx context.Context, opts ...grpc.CallOption) (Service_SubmitMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Service_serviceDesc.Streams[0], "/awakari.writer.Service/SubmitMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceSubmitMessagesClient{stream}
	return x, nil
}

type Service_SubmitMessagesClient interface {
	Send(*SubmitMessagesRequest) error
	Recv() (*SubmitMessagesResponse, error)
	grpc.ClientStream
}

type serviceSubmitMessagesClient struct {
	grpc.ClientStream
}

func (x *serviceSubmitMessagesClient) Send(m *SubmitMessagesRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serviceSubmitMessagesClient) Recv() (*SubmitMessagesResponse, error) {
	m := new(SubmitMessagesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	SubmitMessages(Service_SubmitMessagesServer) error
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) SubmitMessages(Service_SubmitMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubmitMessages not implemented")
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_SubmitMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServiceServer).SubmitMessages(&serviceSubmitMessagesServer{stream})
}

type Service_SubmitMessagesServer interface {
	Send(*SubmitMessagesResponse) error
	Recv() (*SubmitMessagesRequest, error)
	grpc.ServerStream
}

type serviceSubmitMessagesServer struct {
	grpc.ServerStream
}

func (x *serviceSubmitMessagesServer) Send(m *SubmitMessagesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serviceSubmitMessagesServer) Recv() (*SubmitMessagesRequest, error) {
	m := new(SubmitMessagesRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "awakari.writer.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubmitMessages",
			Handler:       _Service_SubmitMessages_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/grpc/writer/service.proto",
}