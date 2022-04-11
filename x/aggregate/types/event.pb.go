// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: teleport/aggregate/v1/event.proto

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

// Status enumerates the status of IBC Aggregate
type Status int32

const (
	// STATUS_UNKNOWN defines the invalid/undefined status
	STATUS_UNKNOWN Status = 0
	// STATUS_SUCCESS defines the success IBC Aggregate execute
	STATUS_SUCCESS Status = 1
	// STATUS_FAILED defines the failed IBC Aggregate execute
	STATUS_FAILED Status = 2
)

var Status_name = map[int32]string{
	0: "STATUS_UNKNOWN",
	1: "STATUS_SUCCESS",
	2: "STATUS_FAILED",
}

var Status_value = map[string]int32{
	"STATUS_UNKNOWN": 0,
	"STATUS_SUCCESS": 1,
	"STATUS_FAILED":  2,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}

func (Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b70b4a642b9b3cca, []int{0}
}

// EventIBCAggregate is emitted on IBC Aggregate
type EventIBCAggregate struct {
	Status             Status `protobuf:"varint,1,opt,name=status,proto3,enum=teleport.aggregate.v1.Status" json:"status,omitempty"`
	Message            string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Sequence           uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
	SourceChannel      string `protobuf:"bytes,4,opt,name=source_channel,json=sourceChannel,proto3" json:"source_channel,omitempty"`
	DestinationChannel string `protobuf:"bytes,5,opt,name=destination_channel,json=destinationChannel,proto3" json:"destination_channel,omitempty"`
}

func (m *EventIBCAggregate) Reset()         { *m = EventIBCAggregate{} }
func (m *EventIBCAggregate) String() string { return proto.CompactTextString(m) }
func (*EventIBCAggregate) ProtoMessage()    {}
func (*EventIBCAggregate) Descriptor() ([]byte, []int) {
	return fileDescriptor_b70b4a642b9b3cca, []int{0}
}
func (m *EventIBCAggregate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventIBCAggregate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventIBCAggregate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventIBCAggregate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventIBCAggregate.Merge(m, src)
}
func (m *EventIBCAggregate) XXX_Size() int {
	return m.Size()
}
func (m *EventIBCAggregate) XXX_DiscardUnknown() {
	xxx_messageInfo_EventIBCAggregate.DiscardUnknown(m)
}

var xxx_messageInfo_EventIBCAggregate proto.InternalMessageInfo

func (m *EventIBCAggregate) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return STATUS_UNKNOWN
}

func (m *EventIBCAggregate) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *EventIBCAggregate) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *EventIBCAggregate) GetSourceChannel() string {
	if m != nil {
		return m.SourceChannel
	}
	return ""
}

func (m *EventIBCAggregate) GetDestinationChannel() string {
	if m != nil {
		return m.DestinationChannel
	}
	return ""
}

// EventRegisterTokens is emitted on aggregate register coins
type EventRegisterTokens struct {
	Denom      []string `protobuf:"bytes,1,rep,name=denom,proto3" json:"denom,omitempty"`
	Erc20Token string   `protobuf:"bytes,2,opt,name=erc20_token,json=erc20Token,proto3" json:"erc20_token,omitempty"`
}

func (m *EventRegisterTokens) Reset()         { *m = EventRegisterTokens{} }
func (m *EventRegisterTokens) String() string { return proto.CompactTextString(m) }
func (*EventRegisterTokens) ProtoMessage()    {}
func (*EventRegisterTokens) Descriptor() ([]byte, []int) {
	return fileDescriptor_b70b4a642b9b3cca, []int{1}
}
func (m *EventRegisterTokens) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventRegisterTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventRegisterTokens.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventRegisterTokens) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventRegisterTokens.Merge(m, src)
}
func (m *EventRegisterTokens) XXX_Size() int {
	return m.Size()
}
func (m *EventRegisterTokens) XXX_DiscardUnknown() {
	xxx_messageInfo_EventRegisterTokens.DiscardUnknown(m)
}

var xxx_messageInfo_EventRegisterTokens proto.InternalMessageInfo

func (m *EventRegisterTokens) GetDenom() []string {
	if m != nil {
		return m.Denom
	}
	return nil
}

func (m *EventRegisterTokens) GetErc20Token() string {
	if m != nil {
		return m.Erc20Token
	}
	return ""
}

func init() {
	proto.RegisterEnum("teleport.aggregate.v1.Status", Status_name, Status_value)
	proto.RegisterType((*EventIBCAggregate)(nil), "teleport.aggregate.v1.EventIBCAggregate")
	proto.RegisterType((*EventRegisterTokens)(nil), "teleport.aggregate.v1.EventRegisterTokens")
}

func init() { proto.RegisterFile("teleport/aggregate/v1/event.proto", fileDescriptor_b70b4a642b9b3cca) }

var fileDescriptor_b70b4a642b9b3cca = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcd, 0xce, 0xd2, 0x40,
	0x14, 0x86, 0x3b, 0xdf, 0x07, 0x28, 0x63, 0x20, 0x30, 0x60, 0xd2, 0x90, 0x58, 0x2b, 0x89, 0x49,
	0x63, 0x62, 0x2b, 0x18, 0xdd, 0x43, 0xc5, 0x84, 0x48, 0xd0, 0xb4, 0x10, 0x13, 0x37, 0xa4, 0x94,
	0x93, 0xa1, 0x01, 0x66, 0xb0, 0x33, 0x45, 0xbd, 0x03, 0x97, 0xde, 0x83, 0x37, 0xe3, 0x92, 0x25,
	0x4b, 0x03, 0x37, 0x62, 0x98, 0xfe, 0x84, 0x85, 0xbb, 0xbe, 0x4f, 0x9f, 0x77, 0x32, 0x67, 0x0e,
	0x7e, 0x26, 0x61, 0x0b, 0x7b, 0x1e, 0x4b, 0x27, 0xa0, 0x34, 0x06, 0x1a, 0x48, 0x70, 0x0e, 0x3d,
	0x07, 0x0e, 0xc0, 0xa4, 0xbd, 0x8f, 0xb9, 0xe4, 0xe4, 0x71, 0xae, 0xd8, 0x85, 0x62, 0x1f, 0x7a,
	0x9d, 0x36, 0xe5, 0x94, 0x2b, 0xc3, 0xb9, 0x7e, 0xa5, 0x72, 0xf7, 0x84, 0x70, 0x73, 0x74, 0x2d,
	0x8f, 0x87, 0xee, 0x20, 0xd7, 0xc9, 0x1b, 0x5c, 0x11, 0x32, 0x90, 0x89, 0xd0, 0x91, 0x89, 0xac,
	0x7a, 0xff, 0x89, 0xfd, 0xdf, 0x33, 0x6d, 0x5f, 0x49, 0x5e, 0x26, 0x13, 0x1d, 0x3f, 0xd8, 0x81,
	0x10, 0x01, 0x05, 0xfd, 0xce, 0x44, 0x56, 0xd5, 0xcb, 0x23, 0xe9, 0xe0, 0x87, 0x02, 0xbe, 0x26,
	0xc0, 0x42, 0xd0, 0xef, 0x4d, 0x64, 0x95, 0xbc, 0x22, 0x93, 0xe7, 0xb8, 0x2e, 0x78, 0x12, 0x87,
	0xb0, 0x08, 0xd7, 0x01, 0x63, 0xb0, 0xd5, 0x4b, 0xaa, 0x5c, 0x4b, 0xa9, 0x9b, 0x42, 0xe2, 0xe0,
	0xd6, 0x0a, 0x84, 0x8c, 0x58, 0x20, 0x23, 0xce, 0x0a, 0xb7, 0xac, 0x5c, 0x72, 0xf3, 0x2b, 0x2b,
	0x74, 0x27, 0xb8, 0xa5, 0x26, 0xf3, 0x80, 0x46, 0x42, 0x42, 0x3c, 0xe3, 0x1b, 0x60, 0x82, 0xb4,
	0x71, 0x79, 0x05, 0x8c, 0xef, 0x74, 0x64, 0xde, 0x5b, 0x55, 0x2f, 0x0d, 0xe4, 0x29, 0x7e, 0x04,
	0x71, 0xd8, 0x7f, 0xb5, 0x90, 0x57, 0x2b, 0xbb, 0x3e, 0x56, 0x48, 0xf5, 0x5e, 0x8c, 0x71, 0x25,
	0x9d, 0x96, 0x10, 0x5c, 0xf7, 0x67, 0x83, 0xd9, 0xdc, 0x5f, 0xcc, 0xa7, 0x1f, 0xa6, 0x1f, 0x3f,
	0x4f, 0x1b, 0xda, 0x0d, 0xf3, 0xe7, 0xae, 0x3b, 0xf2, 0xfd, 0x06, 0x22, 0x4d, 0x5c, 0xcb, 0xd8,
	0xfb, 0xc1, 0x78, 0x32, 0x7a, 0xd7, 0xb8, 0xeb, 0x94, 0x7e, 0xfe, 0x36, 0xb4, 0xe1, 0xa7, 0x3f,
	0x67, 0x03, 0x1d, 0xcf, 0x06, 0xfa, 0x7b, 0x36, 0xd0, 0xaf, 0x8b, 0xa1, 0x1d, 0x2f, 0x86, 0x76,
	0xba, 0x18, 0xda, 0x97, 0xb7, 0x34, 0x92, 0xeb, 0x64, 0x69, 0x87, 0x7c, 0xe7, 0xe4, 0x2f, 0xfe,
	0x92, 0x81, 0xfc, 0xc6, 0xe3, 0x4d, 0x01, 0x9c, 0xef, 0x37, 0xbb, 0x97, 0x3f, 0xf6, 0x20, 0x96,
	0x15, 0xb5, 0xcc, 0xd7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x02, 0x60, 0xaf, 0x63, 0x1e, 0x02,
	0x00, 0x00,
}

func (m *EventIBCAggregate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventIBCAggregate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventIBCAggregate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DestinationChannel) > 0 {
		i -= len(m.DestinationChannel)
		copy(dAtA[i:], m.DestinationChannel)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.DestinationChannel)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.SourceChannel) > 0 {
		i -= len(m.SourceChannel)
		copy(dAtA[i:], m.SourceChannel)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.SourceChannel)))
		i--
		dAtA[i] = 0x22
	}
	if m.Sequence != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Message) > 0 {
		i -= len(m.Message)
		copy(dAtA[i:], m.Message)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Message)))
		i--
		dAtA[i] = 0x12
	}
	if m.Status != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *EventRegisterTokens) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventRegisterTokens) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventRegisterTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Erc20Token) > 0 {
		i -= len(m.Erc20Token)
		copy(dAtA[i:], m.Erc20Token)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Erc20Token)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Denom) > 0 {
		for iNdEx := len(m.Denom) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Denom[iNdEx])
			copy(dAtA[i:], m.Denom[iNdEx])
			i = encodeVarintEvent(dAtA, i, uint64(len(m.Denom[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventIBCAggregate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != 0 {
		n += 1 + sovEvent(uint64(m.Status))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.Sequence != 0 {
		n += 1 + sovEvent(uint64(m.Sequence))
	}
	l = len(m.SourceChannel)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.DestinationChannel)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func (m *EventRegisterTokens) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Denom) > 0 {
		for _, s := range m.Denom {
			l = len(s)
			n += 1 + l + sovEvent(uint64(l))
		}
	}
	l = len(m.Erc20Token)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventIBCAggregate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventIBCAggregate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventIBCAggregate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= Status(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func (m *EventRegisterTokens) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: EventRegisterTokens: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventRegisterTokens: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = append(m.Denom, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Erc20Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Erc20Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
