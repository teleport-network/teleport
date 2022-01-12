// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: xibc/apps/transfer/v1/transfer.proto

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

// FungibleTokenPacketData defines a struct for the packet payload
type FungibleTokenPacketData struct {
	SrcChain  string `protobuf:"bytes,1,opt,name=src_chain,json=srcChain,proto3" json:"src_chain,omitempty"`
	DestChain string `protobuf:"bytes,2,opt,name=dest_chain,json=destChain,proto3" json:"dest_chain,omitempty"`
	Sender    string `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver  string `protobuf:"bytes,4,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Amount    []byte `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount,omitempty"`
	Token     string `protobuf:"bytes,6,opt,name=token,proto3" json:"token,omitempty"`
	OriToken  string `protobuf:"bytes,7,opt,name=ori_token,json=oriToken,proto3" json:"ori_token,omitempty"`
}

func (m *FungibleTokenPacketData) Reset()         { *m = FungibleTokenPacketData{} }
func (m *FungibleTokenPacketData) String() string { return proto.CompactTextString(m) }
func (*FungibleTokenPacketData) ProtoMessage()    {}
func (*FungibleTokenPacketData) Descriptor() ([]byte, []int) {
	return fileDescriptor_f131d793707cfbd6, []int{0}
}
func (m *FungibleTokenPacketData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FungibleTokenPacketData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FungibleTokenPacketData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FungibleTokenPacketData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FungibleTokenPacketData.Merge(m, src)
}
func (m *FungibleTokenPacketData) XXX_Size() int {
	return m.Size()
}
func (m *FungibleTokenPacketData) XXX_DiscardUnknown() {
	xxx_messageInfo_FungibleTokenPacketData.DiscardUnknown(m)
}

var xxx_messageInfo_FungibleTokenPacketData proto.InternalMessageInfo

func (m *FungibleTokenPacketData) GetSrcChain() string {
	if m != nil {
		return m.SrcChain
	}
	return ""
}

func (m *FungibleTokenPacketData) GetDestChain() string {
	if m != nil {
		return m.DestChain
	}
	return ""
}

func (m *FungibleTokenPacketData) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *FungibleTokenPacketData) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *FungibleTokenPacketData) GetAmount() []byte {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *FungibleTokenPacketData) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *FungibleTokenPacketData) GetOriToken() string {
	if m != nil {
		return m.OriToken
	}
	return ""
}

func init() {
	proto.RegisterType((*FungibleTokenPacketData)(nil), "xibc.apps.transfer.v1.FungibleTokenPacketData")
}

func init() {
	proto.RegisterFile("xibc/apps/transfer/v1/transfer.proto", fileDescriptor_f131d793707cfbd6)
}

var fileDescriptor_f131d793707cfbd6 = []byte{
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4e, 0x02, 0x31,
	0x10, 0x86, 0xa9, 0x0a, 0x42, 0xe3, 0x69, 0x83, 0xba, 0xc1, 0xb8, 0x21, 0xc6, 0x03, 0x17, 0xb7,
	0x21, 0x3e, 0x80, 0x89, 0x1a, 0xcf, 0x86, 0x78, 0xd1, 0x0b, 0xe9, 0x96, 0x71, 0x69, 0x80, 0x76,
	0x33, 0x1d, 0x10, 0xdf, 0xc2, 0xc7, 0xf2, 0xc8, 0xc1, 0x83, 0x47, 0x03, 0x2f, 0x62, 0xda, 0x22,
	0x27, 0x6f, 0xfd, 0xe6, 0xfb, 0x67, 0x32, 0x1d, 0x7e, 0xb9, 0xd4, 0x85, 0x12, 0xb2, 0xaa, 0x9c,
	0x20, 0x94, 0xc6, 0xbd, 0x02, 0x8a, 0x45, 0x7f, 0xf7, 0xce, 0x2b, 0xb4, 0x64, 0x93, 0x63, 0x9f,
	0xca, 0x7d, 0x2a, 0xdf, 0x99, 0x45, 0xbf, 0xd3, 0x2e, 0x6d, 0x69, 0x43, 0x42, 0xf8, 0x57, 0x0c,
	0x5f, 0x7c, 0x31, 0x7e, 0xfa, 0x30, 0x37, 0xa5, 0x2e, 0xa6, 0xf0, 0x64, 0x27, 0x60, 0x1e, 0xa5,
	0x9a, 0x00, 0xdd, 0x4b, 0x92, 0xc9, 0x19, 0x6f, 0x39, 0x54, 0x43, 0x35, 0x96, 0xda, 0xa4, 0xac,
	0xcb, 0x7a, 0xad, 0x41, 0xd3, 0xa1, 0xba, 0xf3, 0x9c, 0x9c, 0x73, 0x3e, 0x02, 0x47, 0x5b, 0xbb,
	0x17, 0x6c, 0xcb, 0x57, 0xa2, 0x3e, 0xe1, 0x0d, 0x07, 0x66, 0x04, 0x98, 0xee, 0x07, 0xb5, 0xa5,
	0xa4, 0xc3, 0x9b, 0x08, 0x0a, 0xf4, 0x02, 0x30, 0x3d, 0x88, 0x23, 0xff, 0xd8, 0xf7, 0xc8, 0x99,
	0x9d, 0x1b, 0x4a, 0xeb, 0x5d, 0xd6, 0x3b, 0x1a, 0x6c, 0x29, 0x69, 0xf3, 0x3a, 0xf9, 0xd5, 0xd2,
	0x46, 0x68, 0x88, 0xe0, 0xb7, 0xb3, 0xa8, 0x87, 0xd1, 0x1c, 0xc6, 0x51, 0x16, 0x75, 0xf8, 0xc4,
	0xed, 0xf3, 0xe7, 0x3a, 0x63, 0xab, 0x75, 0xc6, 0x7e, 0xd6, 0x19, 0xfb, 0xd8, 0x64, 0xb5, 0xd5,
	0x26, 0xab, 0x7d, 0x6f, 0xb2, 0xda, 0xcb, 0x4d, 0xa9, 0x69, 0x3c, 0x2f, 0x72, 0x65, 0x67, 0x82,
	0x60, 0x0a, 0x95, 0x45, 0xba, 0x32, 0x40, 0x6f, 0x16, 0x27, 0xbb, 0x82, 0x58, 0x8a, 0x7f, 0x4e,
	0x4d, 0xef, 0x15, 0xb8, 0xa2, 0x11, 0x0e, 0x77, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x88,
	0xbb, 0x47, 0x8d, 0x01, 0x00, 0x00,
}

func (m *FungibleTokenPacketData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FungibleTokenPacketData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FungibleTokenPacketData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OriToken) > 0 {
		i -= len(m.OriToken)
		copy(dAtA[i:], m.OriToken)
		i = encodeVarintTransfer(dAtA, i, uint64(len(m.OriToken)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintTransfer(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Amount) > 0 {
		i -= len(m.Amount)
		copy(dAtA[i:], m.Amount)
		i = encodeVarintTransfer(dAtA, i, uint64(len(m.Amount)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Receiver) > 0 {
		i -= len(m.Receiver)
		copy(dAtA[i:], m.Receiver)
		i = encodeVarintTransfer(dAtA, i, uint64(len(m.Receiver)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTransfer(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.DestChain) > 0 {
		i -= len(m.DestChain)
		copy(dAtA[i:], m.DestChain)
		i = encodeVarintTransfer(dAtA, i, uint64(len(m.DestChain)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.SrcChain) > 0 {
		i -= len(m.SrcChain)
		copy(dAtA[i:], m.SrcChain)
		i = encodeVarintTransfer(dAtA, i, uint64(len(m.SrcChain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTransfer(dAtA []byte, offset int, v uint64) int {
	offset -= sovTransfer(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *FungibleTokenPacketData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SrcChain)
	if l > 0 {
		n += 1 + l + sovTransfer(uint64(l))
	}
	l = len(m.DestChain)
	if l > 0 {
		n += 1 + l + sovTransfer(uint64(l))
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTransfer(uint64(l))
	}
	l = len(m.Receiver)
	if l > 0 {
		n += 1 + l + sovTransfer(uint64(l))
	}
	l = len(m.Amount)
	if l > 0 {
		n += 1 + l + sovTransfer(uint64(l))
	}
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovTransfer(uint64(l))
	}
	l = len(m.OriToken)
	if l > 0 {
		n += 1 + l + sovTransfer(uint64(l))
	}
	return n
}

func sovTransfer(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTransfer(x uint64) (n int) {
	return sovTransfer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FungibleTokenPacketData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTransfer
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
			return fmt.Errorf("proto: FungibleTokenPacketData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FungibleTokenPacketData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrcChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransfer
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
				return ErrInvalidLengthTransfer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransfer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SrcChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransfer
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
				return ErrInvalidLengthTransfer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransfer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransfer
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
				return ErrInvalidLengthTransfer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransfer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransfer
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
				return ErrInvalidLengthTransfer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransfer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransfer
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
				return ErrInvalidLengthTransfer
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTransfer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount[:0], dAtA[iNdEx:postIndex]...)
			if m.Amount == nil {
				m.Amount = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransfer
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
				return ErrInvalidLengthTransfer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransfer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriToken", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTransfer
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
				return ErrInvalidLengthTransfer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTransfer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OriToken = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTransfer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTransfer
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
func skipTransfer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTransfer
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
					return 0, ErrIntOverflowTransfer
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
					return 0, ErrIntOverflowTransfer
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
				return 0, ErrInvalidLengthTransfer
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTransfer
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTransfer
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTransfer        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTransfer          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTransfer = fmt.Errorf("proto: unexpected end of group")
)
