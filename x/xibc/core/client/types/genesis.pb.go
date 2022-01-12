// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: xibc/core/client/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// GenesisState defines the xibc client submodule's genesis state.
type GenesisState struct {
	// client states with their corresponding identifiers
	Clients IdentifiedClientStates `protobuf:"bytes,1,rep,name=clients,proto3,castrepeated=IdentifiedClientStates" json:"clients"`
	// consensus states from each client
	ClientsConsensus ClientsConsensusStates `protobuf:"bytes,2,rep,name=clients_consensus,json=clientsConsensus,proto3,castrepeated=ClientsConsensusStates" json:"clients_consensus"`
	// metadata from each client
	ClientsMetadata []IdentifiedGenesisMetadata `protobuf:"bytes,3,rep,name=clients_metadata,json=clientsMetadata,proto3" json:"clients_metadata"`
	// the chain name of the current chain
	NativeChainName string `protobuf:"bytes,4,opt,name=native_chain_name,json=nativeChainName,proto3" json:"native_chain_name,omitempty"`
	// IdentifiedRelayer defines a list of authorized relayers for the specified
	// client.
	Relayers []IdentifiedRelayers `protobuf:"bytes,5,rep,name=relayers,proto3" json:"relayers"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_81a93b5746096393, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetClients() IdentifiedClientStates {
	if m != nil {
		return m.Clients
	}
	return nil
}

func (m *GenesisState) GetClientsConsensus() ClientsConsensusStates {
	if m != nil {
		return m.ClientsConsensus
	}
	return nil
}

func (m *GenesisState) GetClientsMetadata() []IdentifiedGenesisMetadata {
	if m != nil {
		return m.ClientsMetadata
	}
	return nil
}

func (m *GenesisState) GetNativeChainName() string {
	if m != nil {
		return m.NativeChainName
	}
	return ""
}

func (m *GenesisState) GetRelayers() []IdentifiedRelayers {
	if m != nil {
		return m.Relayers
	}
	return nil
}

// GenesisMetadata defines the genesis type for metadata that clients may return
// with ExportMetadata
type GenesisMetadata struct {
	// store key of metadata without chainName-prefix
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// metadata value
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *GenesisMetadata) Reset()         { *m = GenesisMetadata{} }
func (m *GenesisMetadata) String() string { return proto.CompactTextString(m) }
func (*GenesisMetadata) ProtoMessage()    {}
func (*GenesisMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_81a93b5746096393, []int{1}
}
func (m *GenesisMetadata) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisMetadata.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisMetadata.Merge(m, src)
}
func (m *GenesisMetadata) XXX_Size() int {
	return m.Size()
}
func (m *GenesisMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisMetadata proto.InternalMessageInfo

// IdentifiedGenesisMetadata has the client metadata with the corresponding
// chain name.
type IdentifiedGenesisMetadata struct {
	ChainName string            `protobuf:"bytes,1,opt,name=chain_name,json=chainName,proto3" json:"chain_name,omitempty"`
	Metadata  []GenesisMetadata `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata"`
}

func (m *IdentifiedGenesisMetadata) Reset()         { *m = IdentifiedGenesisMetadata{} }
func (m *IdentifiedGenesisMetadata) String() string { return proto.CompactTextString(m) }
func (*IdentifiedGenesisMetadata) ProtoMessage()    {}
func (*IdentifiedGenesisMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_81a93b5746096393, []int{2}
}
func (m *IdentifiedGenesisMetadata) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IdentifiedGenesisMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IdentifiedGenesisMetadata.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IdentifiedGenesisMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentifiedGenesisMetadata.Merge(m, src)
}
func (m *IdentifiedGenesisMetadata) XXX_Size() int {
	return m.Size()
}
func (m *IdentifiedGenesisMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentifiedGenesisMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_IdentifiedGenesisMetadata proto.InternalMessageInfo

func (m *IdentifiedGenesisMetadata) GetChainName() string {
	if m != nil {
		return m.ChainName
	}
	return ""
}

func (m *IdentifiedGenesisMetadata) GetMetadata() []GenesisMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "xibc.core.client.v1.GenesisState")
	proto.RegisterType((*GenesisMetadata)(nil), "xibc.core.client.v1.GenesisMetadata")
	proto.RegisterType((*IdentifiedGenesisMetadata)(nil), "xibc.core.client.v1.IdentifiedGenesisMetadata")
}

func init() { proto.RegisterFile("xibc/core/client/v1/genesis.proto", fileDescriptor_81a93b5746096393) }

var fileDescriptor_81a93b5746096393 = []byte{
	// 445 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0x41, 0x6b, 0x13, 0x41,
	0x14, 0xde, 0x69, 0x52, 0x6d, 0xc7, 0x42, 0xda, 0xb1, 0xc8, 0x5a, 0x70, 0x13, 0x83, 0x60, 0x28,
	0x38, 0x43, 0xf5, 0x26, 0x88, 0x90, 0x80, 0xd2, 0x83, 0x1e, 0xd6, 0x83, 0xe0, 0x65, 0x99, 0x4c,
	0x9e, 0xe9, 0xd0, 0x64, 0x26, 0xec, 0x4c, 0xb6, 0xcd, 0xd5, 0x93, 0x47, 0x7f, 0x82, 0x67, 0xff,
	0x87, 0xd0, 0x63, 0x8f, 0x9e, 0x54, 0x92, 0x3f, 0x22, 0x3b, 0xb3, 0xbb, 0x4a, 0x5c, 0xcd, 0xed,
	0xed, 0xf7, 0xde, 0xf7, 0x7d, 0xef, 0x7b, 0xbb, 0x8b, 0xef, 0x5f, 0xca, 0xa1, 0x60, 0x42, 0xa7,
	0xc0, 0xc4, 0x44, 0x82, 0xb2, 0x2c, 0x3b, 0x61, 0x63, 0x50, 0x60, 0xa4, 0xa1, 0xb3, 0x54, 0x5b,
	0x4d, 0x6e, 0xe7, 0x23, 0x34, 0x1f, 0xa1, 0x7e, 0x84, 0x66, 0x27, 0x47, 0x9d, 0x3a, 0x5e, 0xd1,
	0x76, 0xb4, 0xa3, 0xc3, 0xb1, 0x1e, 0x6b, 0x57, 0xb2, 0xbc, 0xf2, 0x68, 0xf7, 0x6b, 0x03, 0xef,
	0xbd, 0xf4, 0xf2, 0x6f, 0x2c, 0xb7, 0x40, 0x46, 0xf8, 0xa6, 0xa7, 0x99, 0x10, 0x75, 0x1a, 0xbd,
	0x5b, 0x8f, 0x8f, 0x69, 0x8d, 0x1f, 0x3d, 0x1d, 0x81, 0xb2, 0xf2, 0xbd, 0x84, 0xd1, 0xc0, 0x61,
	0x8e, 0xdc, 0x8f, 0xae, 0xbe, 0xb7, 0x83, 0x2f, 0x3f, 0xda, 0x77, 0x6a, 0xdb, 0x26, 0x2e, 0xa5,
	0xc9, 0x05, 0x3e, 0x28, 0xca, 0x44, 0x68, 0x65, 0x40, 0x99, 0xb9, 0x09, 0xb7, 0xfe, 0xe3, 0xe7,
	0x65, 0x06, 0xe5, 0xac, 0xd7, 0xfb, 0xed, 0xe7, 0xdb, 0x66, 0xad, 0x1f, 0xef, 0x8b, 0x35, 0x9c,
	0x24, 0xb8, 0xc4, 0x92, 0x29, 0x58, 0x3e, 0xe2, 0x96, 0x87, 0x0d, 0xe7, 0x4b, 0x37, 0xe4, 0x2c,
	0xae, 0xf4, 0xaa, 0x60, 0xf5, 0x9b, 0xb9, 0x77, 0xdc, 0x2a, 0xd4, 0x4a, 0x98, 0x1c, 0xe3, 0x03,
	0xc5, 0xad, 0xcc, 0x20, 0x11, 0x67, 0x5c, 0xaa, 0x44, 0xf1, 0x29, 0x84, 0xcd, 0x0e, 0xea, 0xed,
	0xc6, 0x2d, 0xdf, 0x18, 0xe4, 0xf8, 0x6b, 0x3e, 0x05, 0x72, 0x8a, 0x77, 0x52, 0x98, 0xf0, 0x05,
	0xa4, 0x26, 0xdc, 0x76, 0x4b, 0x3c, 0xdc, 0xb0, 0x44, 0x5c, 0x8c, 0x17, 0xee, 0x15, 0xbd, 0xfb,
	0x1c, 0xb7, 0xd6, 0x16, 0x24, 0xfb, 0xb8, 0x71, 0x0e, 0x8b, 0x10, 0x75, 0x50, 0x6f, 0x2f, 0xce,
	0x4b, 0x72, 0x88, 0xb7, 0x33, 0x3e, 0x99, 0x43, 0xb8, 0xe5, 0x30, 0xff, 0xf0, 0xb4, 0xf9, 0xf1,
	0x73, 0x3b, 0xe8, 0x7e, 0x40, 0xf8, 0xee, 0x3f, 0xc3, 0x92, 0x7b, 0x18, 0xff, 0x11, 0x07, 0xb9,
	0x38, 0xbb, 0xa2, 0x0a, 0xf2, 0x02, 0xef, 0x54, 0xd7, 0xf4, 0x6f, 0xf1, 0x41, 0x6d, 0x90, 0xfa,
	0x1b, 0x56, 0xdc, 0xfe, 0xdb, 0xab, 0x65, 0x84, 0xae, 0x97, 0x11, 0xfa, 0xb9, 0x8c, 0xd0, 0xa7,
	0x55, 0x14, 0x5c, 0xaf, 0xa2, 0xe0, 0xdb, 0x2a, 0x0a, 0xde, 0x3d, 0x1b, 0x4b, 0x7b, 0x36, 0x1f,
	0x52, 0xa1, 0xa7, 0xcc, 0xc2, 0x04, 0x66, 0x3a, 0xb5, 0x8f, 0x14, 0xd8, 0x0b, 0x9d, 0x9e, 0x57,
	0x00, 0xbb, 0x64, 0x7f, 0xfd, 0x06, 0x76, 0x31, 0x03, 0x33, 0xbc, 0xe1, 0xbe, 0xf6, 0x27, 0xbf,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x89, 0x96, 0xaf, 0x14, 0x5f, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Relayers) > 0 {
		for iNdEx := len(m.Relayers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Relayers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.NativeChainName) > 0 {
		i -= len(m.NativeChainName)
		copy(dAtA[i:], m.NativeChainName)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.NativeChainName)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ClientsMetadata) > 0 {
		for iNdEx := len(m.ClientsMetadata) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ClientsMetadata[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.ClientsConsensus) > 0 {
		for iNdEx := len(m.ClientsConsensus) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ClientsConsensus[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Clients) > 0 {
		for iNdEx := len(m.Clients) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Clients[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *GenesisMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisMetadata) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisMetadata) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IdentifiedGenesisMetadata) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IdentifiedGenesisMetadata) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IdentifiedGenesisMetadata) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Metadata) > 0 {
		for iNdEx := len(m.Metadata) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Metadata[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.ChainName) > 0 {
		i -= len(m.ChainName)
		copy(dAtA[i:], m.ChainName)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.ChainName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Clients) > 0 {
		for _, e := range m.Clients {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ClientsConsensus) > 0 {
		for _, e := range m.ClientsConsensus {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ClientsMetadata) > 0 {
		for _, e := range m.ClientsMetadata {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = len(m.NativeChainName)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Relayers) > 0 {
		for _, e := range m.Relayers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *GenesisMetadata) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func (m *IdentifiedGenesisMetadata) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainName)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Metadata) > 0 {
		for _, e := range m.Metadata {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Clients", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Clients = append(m.Clients, IdentifiedClientState{})
			if err := m.Clients[len(m.Clients)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientsConsensus", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientsConsensus = append(m.ClientsConsensus, ClientConsensusStates{})
			if err := m.ClientsConsensus[len(m.ClientsConsensus)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientsMetadata", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientsMetadata = append(m.ClientsMetadata, IdentifiedGenesisMetadata{})
			if err := m.ClientsMetadata[len(m.ClientsMetadata)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NativeChainName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NativeChainName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Relayers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Relayers = append(m.Relayers, IdentifiedRelayers{})
			if err := m.Relayers[len(m.Relayers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = append(m.Key[:0], dAtA[iNdEx:postIndex]...)
			if m.Key == nil {
				m.Key = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = append(m.Value[:0], dAtA[iNdEx:postIndex]...)
			if m.Value == nil {
				m.Value = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *IdentifiedGenesisMetadata) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: IdentifiedGenesisMetadata: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IdentifiedGenesisMetadata: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Metadata = append(m.Metadata, GenesisMetadata{})
			if err := m.Metadata[len(m.Metadata)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
