package socket_server

import (
	"github.com/gansidui/gotcp"
	mybase "github.com/giskook/go/base"
	"github.com/giskook/vav-common/base"
	"github.com/giskook/vav-common/protocol"
	"log"
	"os/exec"
	//"runtime/debug"
	//"time"
)

func (ss *SocketServer) OnConnect(c *gotcp.Conn) bool {
	connection := NewConnection(c, ss.conf)
	c.PutExtraData(connection)
	log.Printf("<CNT> %v \n", c.GetRawConn())

	return true
}

func (ss *SocketServer) OnClose(c *gotcp.Conn) {
	connection := c.GetExtraData().(*Connection)
	ss.cm.Del(connection.ID)
	connection.Close()
	log.Printf("<DIS> %v\n", c.GetRawConn())
	err := ss.callback.OnClose(connection)
	if err != nil {
		mybase.ErrorCheckPlus(err, connection.ID)
	}
	//debug.PrintStack()
}

func (ss *SocketServer) prepare(c *Connection, id, channel string) error {
	var err error
	c.once_prepare.Do(func() {
		err = ss.callback.OnPrepare(c, id, channel)
		ss.cm.Put(c.ID, c)
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
				return false
			}

			// do ffmpeg
			go func() {
				connection.once_start_ffmpeg.Do(func() {
					log.Printf("<INFO> %s %s %s\n", rtp.SIM, rtp.LogicalChannel, connection.ffmpeg_cmd)
					cmd := exec.Command("bash", "-c", connection.ffmpeg_cmd)
					_, err := cmd.Output()
					if err != nil {
						log.Printf("<INFO> run ffmpeg error %s %s err msg %s\n", rtp.SIM, rtp.LogicalChannel, err.Error())
					}
					connection.ffmpeg_run = false
				})
			}()
			if !connection.ffmpeg_run {
				return false
			}
			if rtp.Type <= base.RTP_TYPE_VIDEOB {
				//_, err = connection.pipe_v.Write(rtp.Data)
				//if err != nil {
				//	log.Printf("<INFO> %s %s write to video fail err msg :%s \n", rtp.SIM, rtp.LogicalChannel, err.Error())
				//	return false
				//}
				if rtp.Segment == base.RTP_SEGMENT_FIRST ||
					rtp.Segment == base.RTP_SEGMENT_MID {
					connection.frame_vedio.SIM = rtp.SIM
					connection.frame_vedio.LogicalChannel = rtp.LogicalChannel
					connection.frame_vedio.Data = append(connection.frame_vedio.Data, rtp.Data...)
				} else {
					connection.frame_vedio.Data = append(connection.frame_vedio.Data, rtp.Data...)
					_, err = connection.pipe_v.Write(connection.frame_vedio.Data)
					if err != nil {
						log.Printf("<INFO> %s %s write to video fail err msg :%s \n", rtp.SIM, rtp.LogicalChannel, err.Error())
						return false
					}
					connection.frame_vedio.Data = nil
				}
			} else if rtp.Type == base.RTP_TYPE_AUDIO {
				if rtp.Segment == base.RTP_SEGMENT_FIRST ||
					rtp.Segment == base.RTP_SEGMENT_MID {
					connection.frame_audio.SIM = rtp.SIM
					connection.frame_audio.LogicalChannel = rtp.LogicalChannel
					connection.frame_audio.Data = append(connection.frame_audio.Data, rtp.Data...)
				} else {
					connection.frame_audio.Data = append(connection.frame_audio.Data, rtp.Data...)
					_, err = connection.pipe_a.Write(connection.frame_audio.Data)
					if err != nil {
						log.Printf("<INFO> %s %s write to audio fail err msg :%s \n", rtp.SIM, rtp.LogicalChannel, err.Error())
					}
					connection.frame_audio.Data = nil
				}
			}
		}
	}
}
