package socket_server

import (
	"github.com/gansidui/gotcp"
	"github.com/giskook/vav-common/protocol"
	"log"
	//"runtime/debug"
	//"time"
)

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

func (ss *SocketServer) OnConnect(c *gotcp.Conn) bool {
	connection := NewConnection(c, ss.conf, ss.func_prepare)
	c.PutExtraData(connection)
	log.Printf("<CNT> %v \n", c.GetRawConn())

	return true
}

func (ss *SocketServer) OnClose(c *gotcp.Conn) {
	connection := c.GetExtraData().(*Connection)
	//	err := redis_cli.GetInstance().VehicleChannelSet(connection.term.ID, connection.term.LogicalChannel, "0", "")
	//	if err != nil {
	//		log.Printf(ERR_SS_UUID, connection.term.ID, connection.term.LogicalChannel, err.Error())
	//	}
	ss.cm.Del(connection)
	connection.Close()
	log.Printf("<DIS> %v\n", c.GetRawConn())
	ss.callback.OnClose(connection)

	//debug.PrintStack()
}

func (ss *SocketServer) prepare(c *Connection, id, channel string) error {
	var err error
	c.once_prepare.Do(func() {
		err = c.func_prepare(c, id, channel)
	})

	return err
}

func (ss *SocketServer) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	connection := c.GetExtraData().(*Connection)
	connection.recv_buffer.Write(p.Serialize())
	for {
		protocol_id, protocol_length := protocol.CheckProtocol(connection.recv_buffer)
		buf := make([]byte, protocol_length)
		connection.recv_buffer.Read(buf)
		switch protocol_id {
		case protocol.PROTOCOL_HALF_PACK:
			return true
		case protocol.PROTOCOL_ILLEGAL:
			return true
		case protocol.PROTOCOL_RTP:
			rtp := protocol.Parse(buf)
			err := ss.prepare(connection, rtp.SIM, rtp.LogicalChannel)
			if err != nil {
				log.Printf("<INFO> %s %s prepare error %s\n", rtp.SIM, rtp.LogicalChannel, err.Error())
				return false
			}
			if rtp.Type <= RTP_TYPE_VIDEOB {
				if rtp.Segment == RTP_SEGMENT_FIRST ||
					rtp.Segment == RTP_SEGMENT_MID {
					connection.frame_vedio.SIM = rtp.SIM
					connection.frame_vedio.LogicalChannel = rtp.LogicalChannel
					connection.frame_vedio.Data = append(connection.frame_vedio.Data, rtp.Data...)
				} else {
					ok := ss.callback.OnDataVideo(connection)
					if !ok {
						return false
					}
					connection.frame_vedio.Data = nil
				}
			} else if rtp.Type == RTP_TYPE_AUDIO {
				if rtp.Segment == RTP_SEGMENT_FIRST ||
					rtp.Segment == RTP_SEGMENT_MID {
					connection.frame_audio.SIM = rtp.SIM
					connection.frame_audio.LogicalChannel = rtp.LogicalChannel
					connection.frame_audio.Data = append(connection.frame_audio.Data, rtp.Data...)
				} else {
					ok := ss.callback.OnDataAudio(connection)
					if !ok {
						return false
					}
					connection.frame_audio.Data = nil
				}
			}
		}
	}
}
