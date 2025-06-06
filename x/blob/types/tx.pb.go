// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: celestia/blob/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgPayForBlobs pays for the inclusion of a blob in the block.
type MsgPayForBlobs struct {
	// signer is the bech32 encoded signer address. See
	// https://en.bitcoin.it/wiki/Bech32.
	Signer string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	// namespaces is a list of namespaces that the blobs are associated with. A
	// namespace is a byte slice of length 29 where the first byte is the
	// namespaceVersion and the subsequent 28 bytes are the namespaceId.
	Namespaces [][]byte `protobuf:"bytes,2,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
	// blob_sizes is a list of blob sizes (one per blob). Each size is in bytes.
	BlobSizes []uint32 `protobuf:"varint,3,rep,packed,name=blob_sizes,json=blobSizes,proto3" json:"blob_sizes,omitempty"`
	// share_commitments is a list of share commitments (one per blob).
	ShareCommitments [][]byte `protobuf:"bytes,4,rep,name=share_commitments,json=shareCommitments,proto3" json:"share_commitments,omitempty"`
	// share_versions are the versions of the share format that the blobs
	// associated with this message should use when included in a block. The
	// share_versions specified must match the share_versions used to generate the
	// share_commitment in this message.
	ShareVersions []uint32 `protobuf:"varint,8,rep,packed,name=share_versions,json=shareVersions,proto3" json:"share_versions,omitempty"`
}

func (m *MsgPayForBlobs) Reset()         { *m = MsgPayForBlobs{} }
func (m *MsgPayForBlobs) String() string { return proto.CompactTextString(m) }
func (*MsgPayForBlobs) ProtoMessage()    {}
func (*MsgPayForBlobs) Descriptor() ([]byte, []int) {
	return fileDescriptor_9157fbf3d3cd004d, []int{0}
}
func (m *MsgPayForBlobs) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgPayForBlobs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgPayForBlobs.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgPayForBlobs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgPayForBlobs.Merge(m, src)
}
func (m *MsgPayForBlobs) XXX_Size() int {
	return m.Size()
}
func (m *MsgPayForBlobs) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgPayForBlobs.DiscardUnknown(m)
}

var xxx_messageInfo_MsgPayForBlobs proto.InternalMessageInfo

func (m *MsgPayForBlobs) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *MsgPayForBlobs) GetNamespaces() [][]byte {
	if m != nil {
		return m.Namespaces
	}
	return nil
}

func (m *MsgPayForBlobs) GetBlobSizes() []uint32 {
	if m != nil {
		return m.BlobSizes
	}
	return nil
}

func (m *MsgPayForBlobs) GetShareCommitments() [][]byte {
	if m != nil {
		return m.ShareCommitments
	}
	return nil
}

func (m *MsgPayForBlobs) GetShareVersions() []uint32 {
	if m != nil {
		return m.ShareVersions
	}
	return nil
}

// MsgPayForBlobsResponse describes the response returned after the submission
// of a PayForBlobs
type MsgPayForBlobsResponse struct {
}

func (m *MsgPayForBlobsResponse) Reset()         { *m = MsgPayForBlobsResponse{} }
func (m *MsgPayForBlobsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgPayForBlobsResponse) ProtoMessage()    {}
func (*MsgPayForBlobsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9157fbf3d3cd004d, []int{1}
}
func (m *MsgPayForBlobsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgPayForBlobsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgPayForBlobsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgPayForBlobsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgPayForBlobsResponse.Merge(m, src)
}
func (m *MsgPayForBlobsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgPayForBlobsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgPayForBlobsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgPayForBlobsResponse proto.InternalMessageInfo

// MsgUpdateBlobParams defines the sdk.Msg type to update the client parameters.
type MsgUpdateBlobParams struct {
	// authority is the address of the governance account.
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// params defines the blob parameters to update.
	//
	// NOTE: All parameters must be supplied.
	Params Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgUpdateBlobParams) Reset()         { *m = MsgUpdateBlobParams{} }
func (m *MsgUpdateBlobParams) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateBlobParams) ProtoMessage()    {}
func (*MsgUpdateBlobParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_9157fbf3d3cd004d, []int{2}
}
func (m *MsgUpdateBlobParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateBlobParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateBlobParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateBlobParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateBlobParams.Merge(m, src)
}
func (m *MsgUpdateBlobParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateBlobParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateBlobParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateBlobParams proto.InternalMessageInfo

func (m *MsgUpdateBlobParams) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgUpdateBlobParams) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// MsgUpdateBlobParamsResponse defines the MsgUpdateBlobParams response type.
type MsgUpdateBlobParamsResponse struct {
}

func (m *MsgUpdateBlobParamsResponse) Reset()         { *m = MsgUpdateBlobParamsResponse{} }
func (m *MsgUpdateBlobParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateBlobParamsResponse) ProtoMessage()    {}
func (*MsgUpdateBlobParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9157fbf3d3cd004d, []int{3}
}
func (m *MsgUpdateBlobParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateBlobParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateBlobParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateBlobParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateBlobParamsResponse.Merge(m, src)
}
func (m *MsgUpdateBlobParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateBlobParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateBlobParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateBlobParamsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgPayForBlobs)(nil), "celestia.blob.v1.MsgPayForBlobs")
	proto.RegisterType((*MsgPayForBlobsResponse)(nil), "celestia.blob.v1.MsgPayForBlobsResponse")
	proto.RegisterType((*MsgUpdateBlobParams)(nil), "celestia.blob.v1.MsgUpdateBlobParams")
	proto.RegisterType((*MsgUpdateBlobParamsResponse)(nil), "celestia.blob.v1.MsgUpdateBlobParamsResponse")
}

func init() { proto.RegisterFile("celestia/blob/v1/tx.proto", fileDescriptor_9157fbf3d3cd004d) }

var fileDescriptor_9157fbf3d3cd004d = []byte{
	// 521 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x6a, 0xdb, 0x40,
	0x10, 0xf5, 0xc6, 0xa9, 0xa9, 0xd7, 0x8d, 0x71, 0x55, 0xd3, 0xca, 0x6e, 0xac, 0x18, 0x41, 0x40,
	0xb8, 0x58, 0x4a, 0x52, 0xe8, 0xc1, 0xb7, 0xba, 0xd0, 0x43, 0xc1, 0x10, 0x14, 0xda, 0x43, 0x2f,
	0x66, 0x25, 0x6f, 0xd7, 0x02, 0x4b, 0x2b, 0x76, 0x36, 0x26, 0x4e, 0x2f, 0x25, 0x5f, 0x50, 0xe8,
	0x8f, 0xe4, 0xd0, 0x8f, 0xc8, 0x31, 0xb4, 0x97, 0x9e, 0x4a, 0xb1, 0x0b, 0xb9, 0xf7, 0x0b, 0x8a,
	0xb4, 0x52, 0xec, 0xc4, 0x81, 0xf6, 0x36, 0x7a, 0xef, 0xcd, 0x9b, 0x79, 0x9a, 0xc5, 0x0d, 0x9f,
	0x4e, 0x28, 0xc8, 0x80, 0x38, 0xde, 0x84, 0x7b, 0xce, 0x74, 0xdf, 0x91, 0x27, 0x76, 0x2c, 0xb8,
	0xe4, 0x5a, 0x2d, 0xa7, 0xec, 0x84, 0xb2, 0xa7, 0xfb, 0xcd, 0xd6, 0x9a, 0x38, 0x26, 0x82, 0x84,
	0xa0, 0x1a, 0x9a, 0x75, 0xc6, 0x19, 0x4f, 0x4b, 0x27, 0xa9, 0x32, 0x74, 0x9b, 0x71, 0xce, 0x26,
	0xd4, 0x21, 0x71, 0xe0, 0x90, 0x28, 0xe2, 0x92, 0xc8, 0x80, 0x47, 0x79, 0xcf, 0x13, 0x9f, 0x43,
	0xc8, 0xc1, 0x09, 0x81, 0x25, 0x7e, 0x21, 0xb0, 0x8c, 0x68, 0x28, 0x62, 0xa8, 0xfc, 0xd4, 0x87,
	0xa2, 0xcc, 0x39, 0xc2, 0xd5, 0x01, 0xb0, 0x43, 0x32, 0x7b, 0xcd, 0x45, 0x7f, 0xc2, 0x3d, 0xd0,
	0xf6, 0x70, 0x09, 0x02, 0x16, 0x51, 0xa1, 0xa3, 0x36, 0xb2, 0xca, 0x7d, 0xfd, 0xdb, 0xd7, 0x6e,
	0x3d, 0x6b, 0x7a, 0x39, 0x1a, 0x09, 0x0a, 0x70, 0x24, 0x45, 0x10, 0x31, 0x37, 0xd3, 0x69, 0x06,
	0xc6, 0x11, 0x09, 0x29, 0xc4, 0xc4, 0xa7, 0xa0, 0x6f, 0xb4, 0x8b, 0xd6, 0x03, 0x77, 0x05, 0xd1,
	0x5a, 0x18, 0x27, 0x21, 0x87, 0x10, 0x9c, 0x52, 0xd0, 0x8b, 0xed, 0xa2, 0xb5, 0xe5, 0x96, 0x13,
	0xe4, 0x28, 0x01, 0xb4, 0x67, 0xf8, 0x21, 0x8c, 0x89, 0xa0, 0x43, 0x9f, 0x87, 0x61, 0x20, 0x43,
	0x1a, 0x49, 0xd0, 0x37, 0x53, 0x97, 0x5a, 0x4a, 0xbc, 0x5a, 0xe2, 0xda, 0x2e, 0xae, 0x2a, 0xf1,
	0x94, 0x0a, 0x48, 0xc2, 0xeb, 0xf7, 0x53, 0xbf, 0xad, 0x14, 0x7d, 0x97, 0x81, 0xbd, 0xca, 0xd9,
	0xd5, 0x79, 0x27, 0xdb, 0xcf, 0xd4, 0xf1, 0xe3, 0x9b, 0x19, 0x5d, 0x0a, 0x31, 0x8f, 0x80, 0x9a,
	0x1f, 0xf1, 0xa3, 0x01, 0xb0, 0xb7, 0xf1, 0x88, 0x48, 0x9a, 0x30, 0x87, 0xe9, 0x0d, 0xb4, 0x6d,
	0x5c, 0x26, 0xc7, 0x72, 0xcc, 0x45, 0x20, 0x67, 0xea, 0x2f, 0xb8, 0x4b, 0x40, 0x7b, 0x81, 0x4b,
	0xea, 0x56, 0xfa, 0x46, 0x1b, 0x59, 0x95, 0x03, 0xdd, 0xbe, 0x7d, 0x5d, 0x5b, 0xf9, 0xf4, 0x37,
	0x2f, 0x7e, 0xee, 0x14, 0xdc, 0x4c, 0xdd, 0xab, 0x26, 0x3b, 0x2d, 0x7d, 0xcc, 0x16, 0x7e, 0x7a,
	0xc7, 0xf0, 0x7c, 0xb7, 0x83, 0x3f, 0x08, 0x17, 0x07, 0xc0, 0xb4, 0x53, 0x5c, 0x59, 0x3d, 0x4f,
	0x7b, 0x7d, 0xda, 0xcd, 0x70, 0x4d, 0xeb, 0x5f, 0x8a, 0xeb, 0xf8, 0x3b, 0x67, 0xdf, 0x7f, 0x7f,
	0xd9, 0x68, 0xf4, 0x50, 0xc7, 0xac, 0xaf, 0xbc, 0xc3, 0xd9, 0x07, 0x2e, 0xbc, 0x74, 0xd8, 0x18,
	0xd7, 0xd6, 0x7e, 0xce, 0xee, 0x9d, 0xf6, 0xb7, 0x65, 0xcd, 0xee, 0x7f, 0xc9, 0xf2, 0x55, 0x9a,
	0xf7, 0x3e, 0x5d, 0x9d, 0x77, 0x50, 0xff, 0xcd, 0xc5, 0xdc, 0x40, 0x97, 0x73, 0x03, 0xfd, 0x9a,
	0x1b, 0xe8, 0xf3, 0xc2, 0x28, 0x5c, 0x2e, 0x8c, 0xc2, 0x8f, 0x85, 0x51, 0x78, 0xbf, 0xc7, 0x02,
	0x39, 0x3e, 0xf6, 0x6c, 0x9f, 0x87, 0x4e, 0xee, 0xcc, 0x05, 0xbb, 0xae, 0xbb, 0x24, 0x8e, 0x9d,
	0x13, 0x95, 0x42, 0xce, 0x62, 0x0a, 0x5e, 0x29, 0x7d, 0xe2, 0xcf, 0xff, 0x06, 0x00, 0x00, 0xff,
	0xff, 0xb3, 0xf0, 0x8e, 0x82, 0x98, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// PayForBlobs allows the user to pay for the inclusion of one or more blobs
	PayForBlobs(ctx context.Context, in *MsgPayForBlobs, opts ...grpc.CallOption) (*MsgPayForBlobsResponse, error)
	// UpdateBlobParams defines a rpc handler method for MsgUpdateBlobParams.
	UpdateBlobParams(ctx context.Context, in *MsgUpdateBlobParams, opts ...grpc.CallOption) (*MsgUpdateBlobParamsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) PayForBlobs(ctx context.Context, in *MsgPayForBlobs, opts ...grpc.CallOption) (*MsgPayForBlobsResponse, error) {
	out := new(MsgPayForBlobsResponse)
	err := c.cc.Invoke(ctx, "/celestia.blob.v1.Msg/PayForBlobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateBlobParams(ctx context.Context, in *MsgUpdateBlobParams, opts ...grpc.CallOption) (*MsgUpdateBlobParamsResponse, error) {
	out := new(MsgUpdateBlobParamsResponse)
	err := c.cc.Invoke(ctx, "/celestia.blob.v1.Msg/UpdateBlobParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// PayForBlobs allows the user to pay for the inclusion of one or more blobs
	PayForBlobs(context.Context, *MsgPayForBlobs) (*MsgPayForBlobsResponse, error)
	// UpdateBlobParams defines a rpc handler method for MsgUpdateBlobParams.
	UpdateBlobParams(context.Context, *MsgUpdateBlobParams) (*MsgUpdateBlobParamsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) PayForBlobs(ctx context.Context, req *MsgPayForBlobs) (*MsgPayForBlobsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PayForBlobs not implemented")
}
func (*UnimplementedMsgServer) UpdateBlobParams(ctx context.Context, req *MsgUpdateBlobParams) (*MsgUpdateBlobParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBlobParams not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_PayForBlobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPayForBlobs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PayForBlobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/celestia.blob.v1.Msg/PayForBlobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PayForBlobs(ctx, req.(*MsgPayForBlobs))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateBlobParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateBlobParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateBlobParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/celestia.blob.v1.Msg/UpdateBlobParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateBlobParams(ctx, req.(*MsgUpdateBlobParams))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_serviceDesc = _Msg_serviceDesc
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "celestia.blob.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PayForBlobs",
			Handler:    _Msg_PayForBlobs_Handler,
		},
		{
			MethodName: "UpdateBlobParams",
			Handler:    _Msg_UpdateBlobParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "celestia/blob/v1/tx.proto",
}

func (m *MsgPayForBlobs) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgPayForBlobs) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgPayForBlobs) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ShareVersions) > 0 {
		dAtA2 := make([]byte, len(m.ShareVersions)*10)
		var j1 int
		for _, num := range m.ShareVersions {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintTx(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x42
	}
	if len(m.ShareCommitments) > 0 {
		for iNdEx := len(m.ShareCommitments) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ShareCommitments[iNdEx])
			copy(dAtA[i:], m.ShareCommitments[iNdEx])
			i = encodeVarintTx(dAtA, i, uint64(len(m.ShareCommitments[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.BlobSizes) > 0 {
		dAtA4 := make([]byte, len(m.BlobSizes)*10)
		var j3 int
		for _, num := range m.BlobSizes {
			for num >= 1<<7 {
				dAtA4[j3] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j3++
			}
			dAtA4[j3] = uint8(num)
			j3++
		}
		i -= j3
		copy(dAtA[i:], dAtA4[:j3])
		i = encodeVarintTx(dAtA, i, uint64(j3))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Namespaces) > 0 {
		for iNdEx := len(m.Namespaces) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Namespaces[iNdEx])
			copy(dAtA[i:], m.Namespaces[iNdEx])
			i = encodeVarintTx(dAtA, i, uint64(len(m.Namespaces[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgPayForBlobsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgPayForBlobsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgPayForBlobsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgUpdateBlobParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateBlobParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateBlobParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateBlobParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateBlobParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateBlobParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgPayForBlobs) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(m.Namespaces) > 0 {
		for _, b := range m.Namespaces {
			l = len(b)
			n += 1 + l + sovTx(uint64(l))
		}
	}
	if len(m.BlobSizes) > 0 {
		l = 0
		for _, e := range m.BlobSizes {
			l += sovTx(uint64(e))
		}
		n += 1 + sovTx(uint64(l)) + l
	}
	if len(m.ShareCommitments) > 0 {
		for _, b := range m.ShareCommitments {
			l = len(b)
			n += 1 + l + sovTx(uint64(l))
		}
	}
	if len(m.ShareVersions) > 0 {
		l = 0
		for _, e := range m.ShareVersions {
			l += sovTx(uint64(e))
		}
		n += 1 + sovTx(uint64(l)) + l
	}
	return n
}

func (m *MsgPayForBlobsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgUpdateBlobParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Params.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgUpdateBlobParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgPayForBlobs) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgPayForBlobs: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgPayForBlobs: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespaces", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespaces = append(m.Namespaces, make([]byte, postIndex-iNdEx))
			copy(m.Namespaces[len(m.Namespaces)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v uint32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTx
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.BlobSizes = append(m.BlobSizes, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTx
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTx
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTx
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.BlobSizes) == 0 {
					m.BlobSizes = make([]uint32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTx
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.BlobSizes = append(m.BlobSizes, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field BlobSizes", wireType)
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ShareCommitments", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ShareCommitments = append(m.ShareCommitments, make([]byte, postIndex-iNdEx))
			copy(m.ShareCommitments[len(m.ShareCommitments)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType == 0 {
				var v uint32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTx
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint32(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ShareVersions = append(m.ShareVersions, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTx
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTx
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTx
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.ShareVersions) == 0 {
					m.ShareVersions = make([]uint32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTx
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ShareVersions = append(m.ShareVersions, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ShareVersions", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgPayForBlobsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgPayForBlobsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgPayForBlobsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpdateBlobParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdateBlobParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateBlobParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpdateBlobParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdateBlobParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateBlobParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
