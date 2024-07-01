// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: api/grpc/reader/service.proto

package reader

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

type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Command:
	//
	//	*ReadRequest_Start
	//	*ReadRequest_Ack
	Command isReadRequest_Command `protobuf_oneof:"command"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_reader_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_reader_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRequest.ProtoReflect.Descriptor instead.
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return file_api_grpc_reader_service_proto_rawDescGZIP(), []int{0}
}

func (m *ReadRequest) GetCommand() isReadRequest_Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func (x *ReadRequest) GetStart() *ReadCommandStart {
	if x, ok := x.GetCommand().(*ReadRequest_Start); ok {
		return x.Start
	}
	return nil
}

func (x *ReadRequest) GetAck() *ReadCommandAck {
	if x, ok := x.GetCommand().(*ReadRequest_Ack); ok {
		return x.Ack
	}
	return nil
}

type isReadRequest_Command interface {
	isReadRequest_Command()
}

type ReadRequest_Start struct {
	Start *ReadCommandStart `protobuf:"bytes,1,opt,name=start,proto3,oneof"`
}

type ReadRequest_Ack struct {
	Ack *ReadCommandAck `protobuf:"bytes,2,opt,name=ack,proto3,oneof"`
}

func (*ReadRequest_Start) isReadRequest_Command() {}

func (*ReadRequest_Ack) isReadRequest_Command() {}

type ReadCommandStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SubId     string `protobuf:"bytes,1,opt,name=subId,proto3" json:"subId,omitempty"`
	BatchSize uint32 `protobuf:"varint,2,opt,name=batchSize,proto3" json:"batchSize,omitempty"`
}

func (x *ReadCommandStart) Reset() {
	*x = ReadCommandStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_reader_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadCommandStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadCommandStart) ProtoMessage() {}

func (x *ReadCommandStart) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_reader_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadCommandStart.ProtoReflect.Descriptor instead.
func (*ReadCommandStart) Descriptor() ([]byte, []int) {
	return file_api_grpc_reader_service_proto_rawDescGZIP(), []int{1}
}

func (x *ReadCommandStart) GetSubId() string {
	if x != nil {
		return x.SubId
	}
	return ""
}

func (x *ReadCommandStart) GetBatchSize() uint32 {
	if x != nil {
		return x.BatchSize
	}
	return 0
}

type ReadCommandAck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count uint32 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ReadCommandAck) Reset() {
	*x = ReadCommandAck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_reader_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadCommandAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadCommandAck) ProtoMessage() {}

func (x *ReadCommandAck) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_reader_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadCommandAck.ProtoReflect.Descriptor instead.
func (*ReadCommandAck) Descriptor() ([]byte, []int) {
	return file_api_grpc_reader_service_proto_rawDescGZIP(), []int{2}
}

func (x *ReadCommandAck) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgs []*pb.CloudEvent `protobuf:"bytes,1,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (x *ReadResponse) Reset() {
	*x = ReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_grpc_reader_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_grpc_reader_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadResponse.ProtoReflect.Descriptor instead.
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return file_api_grpc_reader_service_proto_rawDescGZIP(), []int{3}
}

func (x *ReadResponse) GetMsgs() []*pb.CloudEvent {
	if x != nil {
		return x.Msgs
	}
	return nil
}

var File_api_grpc_reader_service_proto protoreflect.FileDescriptor

var file_api_grpc_reader_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x1a,
	0x25, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x86, 0x01, 0x0a, 0x0b, 0x52, 0x65, 0x61, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e,
	0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x53, 0x74, 0x61, 0x72, 0x74, 0x48, 0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x12, 0x32, 0x0a, 0x03, 0x61, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x52,
	0x65, 0x61, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x48, 0x00, 0x52,
	0x03, 0x61, 0x63, 0x6b, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22,
	0x46, 0x0a, 0x10, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x75, 0x62, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x73, 0x75, 0x62, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x61, 0x74,
	0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x26, 0x0a, 0x0e, 0x52, 0x65, 0x61, 0x64, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x41, 0x63, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x32, 0x0a, 0x0c, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x22, 0x0a, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x6d,
	0x73, 0x67, 0x73, 0x32, 0x50, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45,
	0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12, 0x1b, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69,
	0x2e, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x77, 0x61, 0x6b, 0x61, 0x72, 0x69, 0x2e, 0x72, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x11, 0x5a, 0x0f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x72, 0x65, 0x61, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_grpc_reader_service_proto_rawDescOnce sync.Once
	file_api_grpc_reader_service_proto_rawDescData = file_api_grpc_reader_service_proto_rawDesc
)

func file_api_grpc_reader_service_proto_rawDescGZIP() []byte {
	file_api_grpc_reader_service_proto_rawDescOnce.Do(func() {
		file_api_grpc_reader_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_grpc_reader_service_proto_rawDescData)
	})
	return file_api_grpc_reader_service_proto_rawDescData
}

var file_api_grpc_reader_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_grpc_reader_service_proto_goTypes = []interface{}{
	(*ReadRequest)(nil),      // 0: awakari.reader.ReadRequest
	(*ReadCommandStart)(nil), // 1: awakari.reader.ReadCommandStart
	(*ReadCommandAck)(nil),   // 2: awakari.reader.ReadCommandAck
	(*ReadResponse)(nil),     // 3: awakari.reader.ReadResponse
	(*pb.CloudEvent)(nil),    // 4: pb.CloudEvent
}
var file_api_grpc_reader_service_proto_depIdxs = []int32{
	1, // 0: awakari.reader.ReadRequest.start:type_name -> awakari.reader.ReadCommandStart
	2, // 1: awakari.reader.ReadRequest.ack:type_name -> awakari.reader.ReadCommandAck
	4, // 2: awakari.reader.ReadResponse.msgs:type_name -> pb.CloudEvent
	0, // 3: awakari.reader.Service.Read:input_type -> awakari.reader.ReadRequest
	3, // 4: awakari.reader.Service.Read:output_type -> awakari.reader.ReadResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_grpc_reader_service_proto_init() }
func file_api_grpc_reader_service_proto_init() {
	if File_api_grpc_reader_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_grpc_reader_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadRequest); i {
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
		file_api_grpc_reader_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadCommandStart); i {
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
		file_api_grpc_reader_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadCommandAck); i {
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
		file_api_grpc_reader_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadResponse); i {
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
	file_api_grpc_reader_service_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ReadRequest_Start)(nil),
		(*ReadRequest_Ack)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_grpc_reader_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_grpc_reader_service_proto_goTypes,
		DependencyIndexes: file_api_grpc_reader_service_proto_depIdxs,
		MessageInfos:      file_api_grpc_reader_service_proto_msgTypes,
	}.Build()
	File_api_grpc_reader_service_proto = out.File
	file_api_grpc_reader_service_proto_rawDesc = nil
	file_api_grpc_reader_service_proto_goTypes = nil
	file_api_grpc_reader_service_proto_depIdxs = nil
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
	// Start reading messages for a certain subscription id.
	// For every response, a client should sent the acknowledged messages count.
	Read(ctx context.Context, opts ...grpc.CallOption) (Service_ReadClient, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Read(ctx context.Context, opts ...grpc.CallOption) (Service_ReadClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Service_serviceDesc.Streams[0], "/awakari.reader.Service/Read", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceReadClient{stream}
	return x, nil
}

type Service_ReadClient interface {
	Send(*ReadRequest) error
	Recv() (*ReadResponse, error)
	grpc.ClientStream
}

type serviceReadClient struct {
	grpc.ClientStream
}

func (x *serviceReadClient) Send(m *ReadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *serviceReadClient) Recv() (*ReadResponse, error) {
	m := new(ReadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	// Start reading messages for a certain subscription id.
	// For every response, a client should sent the acknowledged messages count.
	Read(Service_ReadServer) error
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) Read(Service_ReadServer) error {
	return status.Errorf(codes.Unimplemented, "method Read not implemented")
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_Read_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServiceServer).Read(&serviceReadServer{stream})
}

type Service_ReadServer interface {
	Send(*ReadResponse) error
	Recv() (*ReadRequest, error)
	grpc.ServerStream
}

type serviceReadServer struct {
	grpc.ServerStream
}

func (x *serviceReadServer) Send(m *ReadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *serviceReadServer) Recv() (*ReadRequest, error) {
	m := new(ReadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "awakari.reader.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Read",
			Handler:       _Service_Read_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/grpc/reader/service.proto",
}
