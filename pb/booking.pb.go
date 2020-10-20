// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pb/booking.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type BookingRequest struct {
	// 產品ID
	ProductID int64 `protobuf:"varint,1,opt,name=ProductID,proto3" json:"ProductID,omitempty"`
	// 數量
	Count int64 `protobuf:"varint,2,opt,name=Count,proto3" json:"Count,omitempty"`
	// 注入錯誤: 測試 Rollback
	FaultInject          bool     `protobuf:"varint,3,opt,name=FaultInject,proto3" json:"FaultInject,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BookingRequest) Reset()         { *m = BookingRequest{} }
func (m *BookingRequest) String() string { return proto.CompactTextString(m) }
func (*BookingRequest) ProtoMessage()    {}
func (*BookingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c2a95a0d25f4f8, []int{0}
}
func (m *BookingRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BookingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BookingRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BookingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookingRequest.Merge(m, src)
}
func (m *BookingRequest) XXX_Size() int {
	return m.Size()
}
func (m *BookingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BookingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BookingRequest proto.InternalMessageInfo

func (m *BookingRequest) GetProductID() int64 {
	if m != nil {
		return m.ProductID
	}
	return 0
}

func (m *BookingRequest) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *BookingRequest) GetFaultInject() bool {
	if m != nil {
		return m.FaultInject
	}
	return false
}

type BookingSyncResponse struct {
	// 請求編號
	RequestID string `protobuf:"bytes,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	// 產品ID
	OrderID int64 `protobuf:"varint,2,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	// 付款ID
	PaymentID int64 `protobuf:"varint,3,opt,name=PaymentID,proto3" json:"PaymentID,omitempty"`
	// 注入錯誤: 測試 Rollback
	FaultInject          bool     `protobuf:"varint,4,opt,name=FaultInject,proto3" json:"FaultInject,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BookingSyncResponse) Reset()         { *m = BookingSyncResponse{} }
func (m *BookingSyncResponse) String() string { return proto.CompactTextString(m) }
func (*BookingSyncResponse) ProtoMessage()    {}
func (*BookingSyncResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c2a95a0d25f4f8, []int{1}
}
func (m *BookingSyncResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BookingSyncResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BookingSyncResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BookingSyncResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookingSyncResponse.Merge(m, src)
}
func (m *BookingSyncResponse) XXX_Size() int {
	return m.Size()
}
func (m *BookingSyncResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BookingSyncResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BookingSyncResponse proto.InternalMessageInfo

func (m *BookingSyncResponse) GetRequestID() string {
	if m != nil {
		return m.RequestID
	}
	return ""
}

func (m *BookingSyncResponse) GetOrderID() int64 {
	if m != nil {
		return m.OrderID
	}
	return 0
}

func (m *BookingSyncResponse) GetPaymentID() int64 {
	if m != nil {
		return m.PaymentID
	}
	return 0
}

func (m *BookingSyncResponse) GetFaultInject() bool {
	if m != nil {
		return m.FaultInject
	}
	return false
}

type BookingASyncResponse struct {
	// 請求編號
	RequestID string `protobuf:"bytes,1,opt,name=RequestID,proto3" json:"RequestID,omitempty"`
	// 注入錯誤: 測試 Rollback
	FaultInject          bool     `protobuf:"varint,2,opt,name=FaultInject,proto3" json:"FaultInject,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BookingASyncResponse) Reset()         { *m = BookingASyncResponse{} }
func (m *BookingASyncResponse) String() string { return proto.CompactTextString(m) }
func (*BookingASyncResponse) ProtoMessage()    {}
func (*BookingASyncResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c2a95a0d25f4f8, []int{2}
}
func (m *BookingASyncResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BookingASyncResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BookingASyncResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BookingASyncResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BookingASyncResponse.Merge(m, src)
}
func (m *BookingASyncResponse) XXX_Size() int {
	return m.Size()
}
func (m *BookingASyncResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BookingASyncResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BookingASyncResponse proto.InternalMessageInfo

func (m *BookingASyncResponse) GetRequestID() string {
	if m != nil {
		return m.RequestID
	}
	return ""
}

func (m *BookingASyncResponse) GetFaultInject() bool {
	if m != nil {
		return m.FaultInject
	}
	return false
}

func init() {
	proto.RegisterType((*BookingRequest)(nil), "pb.BookingRequest")
	proto.RegisterType((*BookingSyncResponse)(nil), "pb.BookingSyncResponse")
	proto.RegisterType((*BookingASyncResponse)(nil), "pb.BookingASyncResponse")
}

func init() { proto.RegisterFile("pb/booking.proto", fileDescriptor_f9c2a95a0d25f4f8) }

var fileDescriptor_f9c2a95a0d25f4f8 = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x48, 0xd2, 0x4f,
	0xca, 0xcf, 0xcf, 0xce, 0xcc, 0x4b, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48,
	0x52, 0x4a, 0xe3, 0xe2, 0x73, 0x82, 0x08, 0x06, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0xc9,
	0x70, 0x71, 0x06, 0x14, 0xe5, 0xa7, 0x94, 0x26, 0x97, 0x78, 0xba, 0x48, 0x30, 0x2a, 0x30, 0x6a,
	0x30, 0x07, 0x21, 0x04, 0x84, 0x44, 0xb8, 0x58, 0x9d, 0xf3, 0x4b, 0xf3, 0x4a, 0x24, 0x98, 0xc0,
	0x32, 0x10, 0x8e, 0x90, 0x02, 0x17, 0xb7, 0x5b, 0x62, 0x69, 0x4e, 0x89, 0x67, 0x5e, 0x56, 0x6a,
	0x72, 0x89, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x47, 0x10, 0xb2, 0x90, 0x52, 0x2f, 0x23, 0x97, 0x30,
	0xd4, 0xa2, 0xe0, 0xca, 0xbc, 0xe4, 0xa0, 0xd4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x90, 0x6d,
	0x50, 0x8b, 0xa1, 0xb6, 0x71, 0x06, 0x21, 0x04, 0x84, 0x24, 0xb8, 0xd8, 0xfd, 0x8b, 0x52, 0x52,
	0x8b, 0x3c, 0x5d, 0xa0, 0xf6, 0xc1, 0xb8, 0x60, 0x57, 0x26, 0x56, 0xe6, 0xa6, 0xe6, 0x81, 0xf4,
	0x31, 0x43, 0x5d, 0x09, 0x13, 0x40, 0x77, 0x0f, 0x0b, 0xa6, 0x7b, 0xc2, 0xb8, 0x44, 0xa0, 0xce,
	0x71, 0x24, 0xc1, 0x3d, 0x68, 0xe6, 0x32, 0x61, 0x98, 0x6b, 0x74, 0x8f, 0x11, 0x1e, 0xa0, 0xc1,
	0xa9, 0x45, 0x65, 0x99, 0xc9, 0xa9, 0x42, 0x0e, 0x5c, 0x82, 0x1e, 0x89, 0x79, 0x29, 0x39, 0xa9,
	0x20, 0x8b, 0xa0, 0x72, 0x42, 0x42, 0x7a, 0x05, 0x49, 0x7a, 0xa8, 0x21, 0x2f, 0x25, 0x8e, 0x24,
	0x86, 0xe2, 0x28, 0x27, 0x2e, 0x21, 0x88, 0x09, 0x8e, 0xc5, 0x04, 0x8c, 0x90, 0x40, 0x12, 0x43,
	0xf5, 0x98, 0x1b, 0x97, 0x38, 0xc4, 0x8c, 0x90, 0x8c, 0xa2, 0xfc, 0x92, 0x92, 0x9c, 0xcc, 0xbc,
	0x74, 0x72, 0xdc, 0xe2, 0x24, 0x70, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e,
	0xc9, 0x31, 0xce, 0x78, 0x2c, 0xc7, 0x90, 0xc4, 0x06, 0x4e, 0x4d, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x90, 0x1a, 0xca, 0x0d, 0x61, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BookingServiceClient is the client API for BookingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BookingServiceClient interface {
	// 同步的 Booking 請求
	HandleSyncBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*BookingSyncResponse, error)
	// 異步的 Booking 請求
	HandleAsyncBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*BookingASyncResponse, error)
	// 節流的 Booking 請求
	HandleThrottlingBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*BookingSyncResponse, error)
}

type bookingServiceClient struct {
	cc *grpc.ClientConn
}

func NewBookingServiceClient(cc *grpc.ClientConn) BookingServiceClient {
	return &bookingServiceClient{cc}
}

func (c *bookingServiceClient) HandleSyncBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*BookingSyncResponse, error) {
	out := new(BookingSyncResponse)
	err := c.cc.Invoke(ctx, "/pb.BookingService/HandleSyncBooking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) HandleAsyncBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*BookingASyncResponse, error) {
	out := new(BookingASyncResponse)
	err := c.cc.Invoke(ctx, "/pb.BookingService/HandleAsyncBooking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) HandleThrottlingBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*BookingSyncResponse, error) {
	out := new(BookingSyncResponse)
	err := c.cc.Invoke(ctx, "/pb.BookingService/HandleThrottlingBooking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookingServiceServer is the server API for BookingService service.
type BookingServiceServer interface {
	// 同步的 Booking 請求
	HandleSyncBooking(context.Context, *BookingRequest) (*BookingSyncResponse, error)
	// 異步的 Booking 請求
	HandleAsyncBooking(context.Context, *BookingRequest) (*BookingASyncResponse, error)
	// 節流的 Booking 請求
	HandleThrottlingBooking(context.Context, *BookingRequest) (*BookingSyncResponse, error)
}

// UnimplementedBookingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedBookingServiceServer struct {
}

func (*UnimplementedBookingServiceServer) HandleSyncBooking(ctx context.Context, req *BookingRequest) (*BookingSyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleSyncBooking not implemented")
}
func (*UnimplementedBookingServiceServer) HandleAsyncBooking(ctx context.Context, req *BookingRequest) (*BookingASyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleAsyncBooking not implemented")
}
func (*UnimplementedBookingServiceServer) HandleThrottlingBooking(ctx context.Context, req *BookingRequest) (*BookingSyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleThrottlingBooking not implemented")
}

func RegisterBookingServiceServer(s *grpc.Server, srv BookingServiceServer) {
	s.RegisterService(&_BookingService_serviceDesc, srv)
}

func _BookingService_HandleSyncBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).HandleSyncBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BookingService/HandleSyncBooking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).HandleSyncBooking(ctx, req.(*BookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_HandleAsyncBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).HandleAsyncBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BookingService/HandleAsyncBooking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).HandleAsyncBooking(ctx, req.(*BookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_HandleThrottlingBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).HandleThrottlingBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.BookingService/HandleThrottlingBooking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).HandleThrottlingBooking(ctx, req.(*BookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BookingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BookingService",
	HandlerType: (*BookingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleSyncBooking",
			Handler:    _BookingService_HandleSyncBooking_Handler,
		},
		{
			MethodName: "HandleAsyncBooking",
			Handler:    _BookingService_HandleAsyncBooking_Handler,
		},
		{
			MethodName: "HandleThrottlingBooking",
			Handler:    _BookingService_HandleThrottlingBooking_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/booking.proto",
}

func (m *BookingRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BookingRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BookingRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.FaultInject {
		i--
		if m.FaultInject {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.Count != 0 {
		i = encodeVarintBooking(dAtA, i, uint64(m.Count))
		i--
		dAtA[i] = 0x10
	}
	if m.ProductID != 0 {
		i = encodeVarintBooking(dAtA, i, uint64(m.ProductID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BookingSyncResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BookingSyncResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BookingSyncResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.FaultInject {
		i--
		if m.FaultInject {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.PaymentID != 0 {
		i = encodeVarintBooking(dAtA, i, uint64(m.PaymentID))
		i--
		dAtA[i] = 0x18
	}
	if m.OrderID != 0 {
		i = encodeVarintBooking(dAtA, i, uint64(m.OrderID))
		i--
		dAtA[i] = 0x10
	}
	if len(m.RequestID) > 0 {
		i -= len(m.RequestID)
		copy(dAtA[i:], m.RequestID)
		i = encodeVarintBooking(dAtA, i, uint64(len(m.RequestID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BookingASyncResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BookingASyncResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BookingASyncResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.FaultInject {
		i--
		if m.FaultInject {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.RequestID) > 0 {
		i -= len(m.RequestID)
		copy(dAtA[i:], m.RequestID)
		i = encodeVarintBooking(dAtA, i, uint64(len(m.RequestID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintBooking(dAtA []byte, offset int, v uint64) int {
	offset -= sovBooking(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BookingRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ProductID != 0 {
		n += 1 + sovBooking(uint64(m.ProductID))
	}
	if m.Count != 0 {
		n += 1 + sovBooking(uint64(m.Count))
	}
	if m.FaultInject {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *BookingSyncResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RequestID)
	if l > 0 {
		n += 1 + l + sovBooking(uint64(l))
	}
	if m.OrderID != 0 {
		n += 1 + sovBooking(uint64(m.OrderID))
	}
	if m.PaymentID != 0 {
		n += 1 + sovBooking(uint64(m.PaymentID))
	}
	if m.FaultInject {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *BookingASyncResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RequestID)
	if l > 0 {
		n += 1 + l + sovBooking(uint64(l))
	}
	if m.FaultInject {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovBooking(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBooking(x uint64) (n int) {
	return sovBooking(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BookingRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBooking
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BookingRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BookingRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProductID", wireType)
			}
			m.ProductID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProductID |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FaultInject", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FaultInject = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipBooking(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBooking
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBooking
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BookingSyncResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBooking
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BookingSyncResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BookingSyncResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBooking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBooking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RequestID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderID", wireType)
			}
			m.OrderID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderID |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PaymentID", wireType)
			}
			m.PaymentID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PaymentID |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FaultInject", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FaultInject = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipBooking(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBooking
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBooking
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BookingASyncResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBooking
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BookingASyncResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BookingASyncResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthBooking
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBooking
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RequestID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FaultInject", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FaultInject = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipBooking(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBooking
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBooking
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipBooking(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBooking
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowBooking
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthBooking
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBooking
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBooking
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBooking        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBooking          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBooking = fmt.Errorf("proto: unexpected end of group")
)
