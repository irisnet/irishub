// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: htlc/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the HTLC module's genesis state
type GenesisState struct {
	Params            Params        `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	PendingHtlcs      []HTLC        `protobuf:"bytes,2,rep,name=pending_htlcs,json=pendingHtlcs,proto3" json:"pending_htlcs" yaml:"pending_htlcs"`
	Supplies          []AssetSupply `protobuf:"bytes,3,rep,name=supplies,proto3" json:"supplies"`
	PreviousBlockTime time.Time     `protobuf:"bytes,4,opt,name=previous_block_time,json=previousBlockTime,proto3,stdtime" json:"previous_block_time" yaml:"previous_block_time"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ebc20432ba713fe, []int{0}
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

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetPendingHtlcs() []HTLC {
	if m != nil {
		return m.PendingHtlcs
	}
	return nil
}

func (m *GenesisState) GetSupplies() []AssetSupply {
	if m != nil {
		return m.Supplies
	}
	return nil
}

func (m *GenesisState) GetPreviousBlockTime() time.Time {
	if m != nil {
		return m.PreviousBlockTime
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "irismod.htlc.GenesisState")
}

func init() { proto.RegisterFile("htlc/genesis.proto", fileDescriptor_0ebc20432ba713fe) }

var fileDescriptor_0ebc20432ba713fe = []byte{
	// 358 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x51, 0xb1, 0x6e, 0xea, 0x30,
	0x14, 0x4d, 0x00, 0xa1, 0xa7, 0xc0, 0xd3, 0xd3, 0x4b, 0x19, 0xd2, 0xa8, 0x4a, 0x50, 0x86, 0x8a,
	0xa5, 0xb6, 0x44, 0xb7, 0x76, 0x6a, 0x3a, 0xc0, 0xd0, 0xa1, 0x02, 0xba, 0x74, 0x41, 0x09, 0xb8,
	0xc6, 0x6a, 0x1c, 0x5b, 0xb1, 0x53, 0x89, 0xbf, 0xe0, 0xb3, 0x50, 0x27, 0xc6, 0x4e, 0xb4, 0x82,
	0x3f, 0xe8, 0x17, 0x54, 0x76, 0x4c, 0x55, 0xd4, 0xc5, 0xf2, 0xf5, 0x39, 0xe7, 0xde, 0x73, 0xae,
	0x1d, 0x77, 0x21, 0xb3, 0x19, 0xc4, 0x28, 0x47, 0x82, 0x08, 0xc0, 0x0b, 0x26, 0x99, 0xdb, 0x26,
	0x05, 0x11, 0x94, 0xcd, 0x81, 0xc2, 0xfc, 0x0e, 0x66, 0x98, 0x69, 0x00, 0xaa, 0x5b, 0xc5, 0xf1,
	0xff, 0x69, 0x9d, 0x3a, 0xcc, 0x43, 0x88, 0x19, 0xc3, 0x19, 0x82, 0xba, 0x4a, 0xcb, 0x27, 0x28,
	0x09, 0x45, 0x42, 0x26, 0x94, 0x57, 0x84, 0xe8, 0xb5, 0xe6, 0xb4, 0x07, 0xd5, 0x9c, 0xb1, 0x4c,
	0x24, 0x72, 0xfb, 0x4e, 0x93, 0x27, 0x45, 0x42, 0x85, 0x67, 0x77, 0xed, 0x5e, 0xab, 0xdf, 0x01,
	0x3f, 0xe7, 0x82, 0x7b, 0x8d, 0xc5, 0x8d, 0xf5, 0x36, 0xb4, 0x46, 0x86, 0xe9, 0x3e, 0x38, 0x7f,
	0x39, 0xca, 0xe7, 0x24, 0xc7, 0x53, 0x45, 0x12, 0x5e, 0xad, 0x5b, 0xef, 0xb5, 0xfa, 0xee, 0xb1,
	0x74, 0x38, 0xb9, 0xbb, 0x8d, 0xcf, 0x94, 0xf0, 0x73, 0x1b, 0x76, 0x96, 0x09, 0xcd, 0xae, 0xa2,
	0x23, 0x59, 0x34, 0x6a, 0x9b, 0x7a, 0xa8, 0x4a, 0xf7, 0xda, 0xf9, 0x23, 0x4a, 0xce, 0x33, 0x82,
	0x84, 0x57, 0xd7, 0x1d, 0x4f, 0x8f, 0x3b, 0xde, 0x08, 0x81, 0xe4, 0x58, 0x51, 0x96, 0xc6, 0xd1,
	0xb7, 0xc0, 0x2d, 0x9c, 0x13, 0x5e, 0xa0, 0x17, 0xc2, 0x4a, 0x31, 0x4d, 0x33, 0x36, 0x7b, 0x9e,
	0xaa, 0xe8, 0x5e, 0x43, 0x87, 0xf2, 0x41, 0xb5, 0x17, 0x70, 0xd8, 0x0b, 0x98, 0x1c, 0xf6, 0x12,
	0x9f, 0x1b, 0x87, 0xbe, 0x71, 0xf8, 0xbb, 0x49, 0xb4, 0x7a, 0x0f, 0xed, 0xd1, 0xff, 0x03, 0x12,
	0x2b, 0x40, 0xe9, 0xe3, 0xc1, 0x7a, 0x17, 0xd8, 0x9b, 0x5d, 0x60, 0x7f, 0xec, 0x02, 0x7b, 0xb5,
	0x0f, 0xac, 0xcd, 0x3e, 0xb0, 0xde, 0xf6, 0x81, 0xf5, 0x78, 0x81, 0x89, 0x5c, 0x94, 0x29, 0x98,
	0x31, 0x0a, 0x55, 0x84, 0x1c, 0x49, 0x68, 0xa2, 0x40, 0xca, 0xe6, 0x65, 0x86, 0x84, 0xfe, 0x36,
	0x28, 0x97, 0x1c, 0x89, 0xb4, 0xa9, 0x7d, 0x5d, 0x7e, 0x05, 0x00, 0x00, 0xff, 0xff, 0x14, 0xc0,
	0x55, 0xa7, 0x08, 0x02, 0x00, 0x00,
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
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.PreviousBlockTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.PreviousBlockTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintGenesis(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	if len(m.Supplies) > 0 {
		for iNdEx := len(m.Supplies) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Supplies[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.PendingHtlcs) > 0 {
		for iNdEx := len(m.PendingHtlcs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PendingHtlcs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
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
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.PendingHtlcs) > 0 {
		for _, e := range m.PendingHtlcs {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Supplies) > 0 {
		for _, e := range m.Supplies {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.PreviousBlockTime)
	n += 1 + l + sovGenesis(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
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
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PendingHtlcs", wireType)
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
			m.PendingHtlcs = append(m.PendingHtlcs, HTLC{})
			if err := m.PendingHtlcs[len(m.PendingHtlcs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Supplies", wireType)
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
			m.Supplies = append(m.Supplies, AssetSupply{})
			if err := m.Supplies[len(m.Supplies)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousBlockTime", wireType)
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
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.PreviousBlockTime, dAtA[iNdEx:postIndex]); err != nil {
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