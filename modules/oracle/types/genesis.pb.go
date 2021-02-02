// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: oracle/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/irisnet/irismod/modules/service/types"
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

// GenesisState defines the oracle module's genesis state.
type GenesisState struct {
	Entries []FeedEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce0b3a2b4a184fc3, []int{0}
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

func (m *GenesisState) GetEntries() []FeedEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type FeedEntry struct {
	Feed   Feed                      `protobuf:"bytes,1,opt,name=feed,proto3" json:"feed"`
	State  types.RequestContextState `protobuf:"varint,2,opt,name=state,proto3,enum=irismod.service.RequestContextState" json:"state,omitempty"`
	Values []FeedValue               `protobuf:"bytes,3,rep,name=values,proto3" json:"values"`
}

func (m *FeedEntry) Reset()         { *m = FeedEntry{} }
func (m *FeedEntry) String() string { return proto.CompactTextString(m) }
func (*FeedEntry) ProtoMessage()    {}
func (*FeedEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce0b3a2b4a184fc3, []int{1}
}
func (m *FeedEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeedEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeedEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FeedEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeedEntry.Merge(m, src)
}
func (m *FeedEntry) XXX_Size() int {
	return m.Size()
}
func (m *FeedEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_FeedEntry.DiscardUnknown(m)
}

var xxx_messageInfo_FeedEntry proto.InternalMessageInfo

func (m *FeedEntry) GetFeed() Feed {
	if m != nil {
		return m.Feed
	}
	return Feed{}
}

func (m *FeedEntry) GetState() types.RequestContextState {
	if m != nil {
		return m.State
	}
	return types.RUNNING
}

func (m *FeedEntry) GetValues() []FeedValue {
	if m != nil {
		return m.Values
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "irismod.oracle.GenesisState")
	proto.RegisterType((*FeedEntry)(nil), "irismod.oracle.FeedEntry")
}

func init() { proto.RegisterFile("oracle/genesis.proto", fileDescriptor_ce0b3a2b4a184fc3) }

var fileDescriptor_ce0b3a2b4a184fc3 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x63, 0x5a, 0x8a, 0x70, 0x51, 0x87, 0x50, 0xa4, 0xd0, 0xc1, 0x54, 0x15, 0x43, 0x26,
	0x5b, 0x0a, 0x03, 0x82, 0xb1, 0x08, 0x50, 0xd7, 0x20, 0x31, 0xb0, 0xa5, 0xc9, 0x11, 0x2c, 0x25,
	0x71, 0xb1, 0x9d, 0x8a, 0xbe, 0x05, 0x8f, 0xc2, 0x63, 0x74, 0xec, 0xc8, 0x84, 0x50, 0xf2, 0x22,
	0x28, 0xb1, 0x5b, 0x09, 0x09, 0x26, 0x9f, 0xee, 0xbe, 0xfb, 0xff, 0xdf, 0x87, 0x87, 0x42, 0x46,
	0x71, 0x06, 0x2c, 0x85, 0x02, 0x14, 0x57, 0x74, 0x21, 0x85, 0x16, 0xee, 0x80, 0x4b, 0xae, 0x72,
	0x91, 0x50, 0x33, 0x1d, 0x1d, 0x5b, 0xca, 0x3c, 0x06, 0x1a, 0x9d, 0x28, 0x90, 0x4b, 0x1e, 0x03,
	0xb3, 0xaf, 0x6d, 0x0f, 0x53, 0x91, 0x8a, 0xb6, 0x64, 0x4d, 0x65, 0xba, 0x93, 0x19, 0x3e, 0xba,
	0x37, 0x16, 0x0f, 0x3a, 0xd2, 0xe0, 0x5e, 0xe1, 0x03, 0x28, 0xb4, 0xe4, 0xa0, 0x3c, 0x34, 0xee,
	0xf8, 0xfd, 0xe0, 0x94, 0xfe, 0xf6, 0xa4, 0x77, 0x00, 0xc9, 0x6d, 0xa1, 0xe5, 0x6a, 0xda, 0x5d,
	0x7f, 0x9d, 0x39, 0xe1, 0x96, 0x9f, 0x7c, 0x20, 0x7c, 0xb8, 0x1b, 0xba, 0x14, 0x77, 0x9f, 0x01,
	0x12, 0x0f, 0x8d, 0x91, 0xdf, 0x0f, 0x86, 0x7f, 0xa9, 0x58, 0x81, 0x96, 0x73, 0xaf, 0xf1, 0xbe,
	0x6a, 0x12, 0x78, 0x7b, 0x63, 0xe4, 0x0f, 0x82, 0xf3, 0xdd, 0xc2, 0xf6, 0x17, 0x21, 0xbc, 0x96,
	0xa0, 0xf4, 0x8d, 0x28, 0x34, 0xbc, 0xe9, 0x36, 0x6d, 0x68, 0x56, 0xdc, 0x4b, 0xdc, 0x5b, 0x46,
	0x59, 0x09, 0xca, 0xeb, 0xfc, 0x9f, 0xf9, 0xb1, 0x21, 0xac, 0xa5, 0xc5, 0xa7, 0xb3, 0x75, 0x45,
	0xd0, 0xa6, 0x22, 0xe8, 0xbb, 0x22, 0xe8, 0xbd, 0x26, 0xce, 0xa6, 0x26, 0xce, 0x67, 0x4d, 0x9c,
	0x27, 0x96, 0x72, 0xfd, 0x52, 0xce, 0x69, 0x2c, 0x72, 0xd6, 0x88, 0x15, 0xa0, 0x99, 0x15, 0x65,
	0xb9, 0x48, 0xca, 0x0c, 0x94, 0xbd, 0x3a, 0xd3, 0xab, 0x05, 0xa8, 0x79, 0xaf, 0xbd, 0xe7, 0xc5,
	0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x10, 0xb1, 0x02, 0xf4, 0xb9, 0x01, 0x00, 0x00,
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
	if len(m.Entries) > 0 {
		for iNdEx := len(m.Entries) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Entries[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *FeedEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeedEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FeedEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Values) > 0 {
		for iNdEx := len(m.Values) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Values[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if m.State != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.Feed.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
	if len(m.Entries) > 0 {
		for _, e := range m.Entries {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *FeedEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Feed.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.State != 0 {
		n += 1 + sovGenesis(uint64(m.State))
	}
	if len(m.Values) > 0 {
		for _, e := range m.Values {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Entries", wireType)
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
			m.Entries = append(m.Entries, FeedEntry{})
			if err := m.Entries[len(m.Entries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *FeedEntry) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: FeedEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeedEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feed", wireType)
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
			if err := m.Feed.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= types.RequestContextState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Values", wireType)
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
			m.Values = append(m.Values, FeedValue{})
			if err := m.Values[len(m.Values)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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