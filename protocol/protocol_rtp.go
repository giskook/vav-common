package protocol

import (
	"bytes"
	"github.com/giskook/vav/base"
	"strconv"
)

func Parse(p []byte) *base.RTP {
	r := bytes.NewReader(p)
	r.Seek(4, 0)
	vpxcc := base.ReadByte(r)
	mpt := base.ReadByte(r)
	serial := base.ReadWord(r)
	sim := base.ReadBcdString(r, 6)
	logical_channel := base.ReadByte(r)
	type_seg := base.ReadByte(r)
	t := type_seg & 0xf0
	seg := (type_seg & 0x0f)
	timestamp := uint64(0)
	if t != base.RTP_TYPE_RAW {
		timestamp = base.ReadQuaWord(r)
	}
	var last_i_frame_interval, last_frame_interval uint16 = 0, 0
	if t == base.RTP_TYPE_VIDEOI ||
		t == base.RTP_TYPE_VIDEOP ||
		t == base.RTP_TYPE_VIDEOB {
		last_i_frame_interval = base.ReadWord(r)
		last_frame_interval = base.ReadWord(r)
	}

	length := base.ReadWord(r)
	data := base.ReadBytes(r, int(length))

	return &base.RTP{
		Serial:             serial,
		SIM:                sim,
		LogicalChannel:     strconv.FormatUint(uint64(logical_channel), 10),
		Type:               t,
		Segment:            seg,
		Timestamp:          timestamp,
		LastIFrameInterval: last_i_frame_interval,
		LastFrameInterval:  last_frame_interval,
		Data:               data,
		VPXCC:              vpxcc,
		MPT:                mpt,
	}
}

const (
	PROTOCOL_ILLEGAL           uint16 = 0xffff
	PROTOCOL_HALF_PACK         uint16 = 0xfffe
	PROTOCOL_RTP               uint16 = 0x0000
	PROTOCOL_FLAG              string = "01cd"
	PROTOCOL_MIN_LEN           uint16 = 30
	PROTOCOL_HEADER_COMMON_LEN uint16 = 16

	PROTOCOL_RTP_DT_FRAME_I uint8 = 0x00
	PROTOCOL_RTP_DT_FRAME_P uint8 = 0x10
	PROTOCOL_RTP_DT_FRAME_B uint8 = 0x20
	PROTOCOL_RTP_DT_FRAME_A uint8 = 0x30
	PROTOCOL_RTP_DT_FRAME_T uint8 = 0x40
	PROTOCOL_RTP_DT_MASK    uint8 = 0xf0
)

func CheckProtocol(b *bytes.Buffer) (uint16, uint16) {
	cmd := PROTOCOL_ILLEGAL
	cmd_len := uint16(0)
	bufferlen := uint16(b.Len())
	if bufferlen == 0 {
		return PROTOCOL_ILLEGAL, 0
	}
	if bufferlen >= PROTOCOL_MIN_LEN {
		if string(b.Bytes()[0:4]) != PROTOCOL_FLAG {
			b.ReadByte()
			cmd, cmd_len = CheckProtocol(b)
		}
		dt := b.Bytes()[15] & PROTOCOL_RTP_DT_MASK
		if dt >= PROTOCOL_RTP_DT_FRAME_A {
			if dt == PROTOCOL_RTP_DT_FRAME_T {
				cmd_len = PROTOCOL_HEADER_COMMON_LEN + base.GetWord(b.Bytes()[16:18])
			} else {
				cmd_len = PROTOCOL_HEADER_COMMON_LEN + 8 + base.GetWord(b.Bytes()[24:26])
			}
		} else {
			cmd_len = PROTOCOL_HEADER_COMMON_LEN + 14 + base.GetWord(b.Bytes()[28:30])
		}
		if bufferlen < cmd_len {
			return PROTOCOL_HALF_PACK, 0
		}
	} else {
		return PROTOCOL_HALF_PACK, 0
	}
	cmd = PROTOCOL_RTP

	return cmd, cmd_len
}
