package socket_server

import (
	"fmt"
	"github.com/gansidui/gotcp"
	"github.com/giskook/vav-common/base"
	"github.com/giskook/vav-common/protocol"
	"log"
	"os"
	"os/exec"
	//"runtime/debug"
	"time"
)

func Copy(c *Connection, in, out string) error {
	fin, err := os.OpenFile(in, os.O_RDONLY, 0600)
	defer fin.Close()
	if err != nil {
		return err
	}
	fout, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fout.Close()
	b := make([]byte, 1024)
	for {
		select {
		case <-c.exit:
			return nil
		default:
			n, err := fin.Read(b)
			if err != nil {
				return nil
			}
			_, err = fout.Write(b[0:n])
			if err != nil {
				return err
			}
			fout.Sync()
		}
	}
}

func (ss *SocketServer) OnConnect(c *gotcp.Conn) bool {
	connection := NewConnection(c, ss.conf)
	c.PutExtraData(connection)
	log.Printf("<CNT> %v \n", c.GetRawConn())

	return true
}

func (ss *SocketServer) OnClose(c *gotcp.Conn) {
	connection := c.GetExtraData().(*Connection)
	ss.cm.Del(connection.ID)
	log.Printf("<DIS> %v\n", c.GetRawConn())
	connection.ShutDown()
	ss.callback.OnClose(connection)
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
	c.GetRawConn().SetReadDeadline(time.Now().Add(ss.conf.DefaultReadLimit))

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
			if int(time.Now().Unix())-connection.TimeStamp > connection.TTL {
				return false
			}

			// do ffmpeg
			start_ffmpeg := func(avtype int, cmds []string) {
				if avtype != VIDEO_TYPE &&
					avtype != AUDIO_TYPE &&
					avtype != AV_TYPE {
					return
				}
				var cmd string
				if connection.avtype != avtype {
					connection.avtype = avtype
					cmd = cmds[avtype-1]
				} else {
					return
				}
				log.Printf("<INFO> %s %s %s\n", rtp.SIM, rtp.LogicalChannel, cmd)
				do_ffmpeg := func(ffmpeg_cmd string) {
					ffmpeg_killer := fmt.Sprintf(connection.conf.FFmpegKiller, connection.ffmpeg_name)
					cmd_quit := exec.Command("bash", "-c", ffmpeg_killer)
					cmd_quit.Output()
					cmd := exec.Command("bash", "-c", ffmpeg_cmd)
					_, err = cmd.Output()
					if err != nil {
						log.Printf("<INFO> run ffmpeg error %s %s err msg %s\n", rtp.SIM, rtp.LogicalChannel, err.Error())
					}
					connection.ffmpeg_run = false
					//ss.callback.OnFfmpegExit(connection)
				}
				if !ss.conf.Debug.Debug {
					do_ffmpeg(cmd)
				} else {
					if connection.SIM == ss.conf.Debug.DestID {
						if ss.conf.Debug.RecordFileA {
							Copy(connection, connection.file_path_a, "./"+connection.SIM+"_"+connection.Channel+"."+connection.acodec)
							connection.ffmpeg_run = false
						}
					} else {
						do_ffmpeg(cmd)
					}
				}
			}
			if !connection.ffmpeg_run {
				return false
			}
			t := time.Now().Unix()
			connection.elapsed_time = t
			log.Println(connection.elapsed_time, connection.baseline_time, rtp.Type)
			if connection.elapsed_time-connection.baseline_time > 1 {
				connection.baseline_time = connection.elapsed_time
				avtype := 0
				if t-connection.video_last_time < 1 {
					avtype |= VIDEO_TYPE
				}
				if t-connection.audio_last_time < 1 {
					avtype |= AUDIO_TYPE
				}
				go start_ffmpeg(avtype, connection.ffmpeg_cmds)
			}
			if rtp.Type <= base.RTP_TYPE_VIDEOB {
				log.Println("v")
				connection.video_last_time = t
				if rtp.Segment == base.RTP_SEGMENT_FIRST ||
					rtp.Segment == base.RTP_SEGMENT_MID {
					connection.frame_vedio.SIM = rtp.SIM
					connection.frame_vedio.LogicalChannel = rtp.LogicalChannel
					connection.frame_vedio.Data = append(connection.frame_vedio.Data, rtp.Data...)
				} else {
					if !connection.log_video {
						connection.log_video = true
						log.Printf("<INFO> %s %s upload video\n", rtp.SIM, rtp.LogicalChannel)
					}
					connection.frame_vedio.Data = append(connection.frame_vedio.Data, rtp.Data...)
					_, err = connection.pipe_v.Write(connection.frame_vedio.Data)
					if err != nil {
						log.Printf("<INFO> %s %s write to video fail err msg :%s \n", rtp.SIM, rtp.LogicalChannel, err.Error())
						return false
					}
					connection.frame_vedio.Data = nil
				}
			} else if rtp.Type == base.RTP_TYPE_AUDIO {
				log.Println("a")
				connection.audio_last_time = t
				if rtp.Segment == base.RTP_SEGMENT_FIRST ||
					rtp.Segment == base.RTP_SEGMENT_MID {
					connection.frame_audio.SIM = rtp.SIM
					connection.frame_audio.LogicalChannel = rtp.LogicalChannel
					connection.frame_audio.Data = append(connection.frame_audio.Data, rtp.Data...)
				} else {
					if !connection.log_audio {
						connection.log_audio = true
						log.Printf("<INFO> %s %s upload audio\n", rtp.SIM, rtp.LogicalChannel)
					}
					connection.frame_audio.Data = append(connection.frame_audio.Data, rtp.Data...)
					_, err = connection.pipe_a.Write(connection.frame_audio.Data)
					if err != nil {
						log.Printf("<INFO> %s %s write to audio fail err msg :%s \n", rtp.SIM, rtp.LogicalChannel, err.Error())
						return false
					}
					connection.frame_audio.Data = nil
				}
			}
		}
	}
}
