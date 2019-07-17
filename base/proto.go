package base

const (
	RTP_TYPE_VIDEOI uint8 = 0x00
	RTP_TYPE_VIDEOP uint8 = 0x10
	RTP_TYPE_VIDEOB uint8 = 0x20
	RTP_TYPE_AUDIO  uint8 = 0x30
	RTP_TYPE_RAW    uint8 = 0x40

	RTP_SEGMENT_COMPLETE uint8 = 0x00
	RTP_SEGMENT_FIRST    uint8 = 0x01
	RTP_SEGMENT_LAST     uint8 = 0x02
	RTP_SEGMENT_MID      uint8 = 0x03
)

type RTP struct {
	VPXCC              uint8
	MPT                uint8
	Serial             uint16
	SIM                string
	LogicalChannel     string
	Type               uint8
	Segment            uint8
	Timestamp          uint64
	LastIFrameInterval uint16
	LastFrameInterval  uint16
	Data               []byte
}

type Frame struct {
	SIM            string
	LogicalChannel string
	Data           []byte
}
