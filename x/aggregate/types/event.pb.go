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

func init() {
	proto.RegisterEnum("teleport.aggregate.v1.Status", Status_name, Status_value)
	proto.RegisterType((*EventIBCAggregate)(nil), "teleport.aggregate.v1.EventIBCAggregate")
}

func init() { proto.RegisterFile("teleport/aggregate/v1/event.proto", fileDescriptor_b70b4a642b9b3cca) }

var fileDescriptor_b70b4a642b9b3cca = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x41, 0x4b, 0x32, 0x41,
	0x18, 0xc7, 0x77, 0x7c, 0x7d, 0xad, 0x06, 0x14, 0x9d, 0x0a, 0x16, 0xa1, 0xc1, 0x82, 0x40, 0x82,
	0x76, 0xb0, 0xa8, 0xbb, 0x6e, 0x06, 0x52, 0x58, 0xb8, 0x4a, 0xd0, 0x45, 0xd6, 0xed, 0x61, 0x94,
	0x74, 0xc6, 0x76, 0x66, 0xad, 0xbe, 0x41, 0xc7, 0xbe, 0x43, 0x5f, 0xa6, 0xa3, 0x47, 0x8f, 0xa1,
	0x5f, 0x24, 0xdc, 0x75, 0x97, 0x3d, 0x74, 0x9b, 0xe7, 0x37, 0xbf, 0xff, 0x30, 0xfc, 0x1f, 0x7c,
	0xa8, 0x61, 0x0c, 0x53, 0xe9, 0x6b, 0xe6, 0x72, 0xee, 0x03, 0x77, 0x35, 0xb0, 0x59, 0x8d, 0xc1,
	0x0c, 0x84, 0xb6, 0xa6, 0xbe, 0xd4, 0x92, 0xec, 0xc7, 0x8a, 0x95, 0x28, 0xd6, 0xac, 0x56, 0xde,
	0xe3, 0x92, 0xcb, 0xd0, 0x60, 0xeb, 0x53, 0x24, 0x1f, 0x2d, 0x10, 0x2e, 0x35, 0xd7, 0xe1, 0x56,
	0xc3, 0xae, 0xc7, 0x3a, 0xb9, 0xc0, 0x39, 0xa5, 0x5d, 0x1d, 0x28, 0x13, 0x55, 0x50, 0xb5, 0x70,
	0x76, 0x60, 0xfd, 0xf9, 0xa6, 0xe5, 0x84, 0x52, 0x67, 0x23, 0x13, 0x13, 0x6f, 0x4d, 0x40, 0x29,
	0x97, 0x83, 0x99, 0xa9, 0xa0, 0xea, 0x4e, 0x27, 0x1e, 0x49, 0x19, 0x6f, 0x2b, 0x78, 0x09, 0x40,
	0x78, 0x60, 0xfe, 0xab, 0xa0, 0x6a, 0xb6, 0x93, 0xcc, 0xe4, 0x18, 0x17, 0x94, 0x0c, 0x7c, 0x0f,
	0xfa, 0xde, 0xd0, 0x15, 0x02, 0xc6, 0x66, 0x36, 0x0c, 0xe7, 0x23, 0x6a, 0x47, 0x90, 0x30, 0xbc,
	0xfb, 0x04, 0x4a, 0x8f, 0x84, 0xab, 0x47, 0x52, 0x24, 0xee, 0xff, 0xd0, 0x25, 0xa9, 0xab, 0x4d,
	0xe0, 0xa4, 0x85, 0x73, 0xd1, 0xff, 0x08, 0xc1, 0x05, 0xa7, 0x5b, 0xef, 0xf6, 0x9c, 0x7e, 0xaf,
	0x7d, 0xd3, 0xbe, 0x7b, 0x68, 0x17, 0x8d, 0x14, 0x73, 0x7a, 0xb6, 0xdd, 0x74, 0x9c, 0x22, 0x22,
	0x25, 0x9c, 0xdf, 0xb0, 0xeb, 0x7a, 0xeb, 0xb6, 0x79, 0x55, 0xcc, 0x94, 0xb3, 0x1f, 0x5f, 0xd4,
	0x68, 0xdc, 0x7f, 0x2f, 0x29, 0x9a, 0x2f, 0x29, 0xfa, 0x59, 0x52, 0xf4, 0xb9, 0xa2, 0xc6, 0x7c,
	0x45, 0x8d, 0xc5, 0x8a, 0x1a, 0x8f, 0x97, 0x7c, 0xa4, 0x87, 0xc1, 0xc0, 0xf2, 0xe4, 0x84, 0xc5,
	0x1d, 0x9d, 0x0a, 0xd0, 0xaf, 0xd2, 0x7f, 0x4e, 0x00, 0x7b, 0x4b, 0x6d, 0x4b, 0xbf, 0x4f, 0x41,
	0x0d, 0x72, 0x61, 0xfd, 0xe7, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x15, 0x46, 0xde, 0xe9, 0xd0,
	0x01, 0x00, 0x00,
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
