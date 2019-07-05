package base

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
