package socket_server

import (
	"bytes"
	"github.com/gansidui/gotcp"
	"github.com/giskook/vav-common/base"
	"os"
	"sync"
	"time"
)

type ConnCallback interface {
	OnPrepare(*Connection, string, string) error
	OnClose(*Connection) error
}

const (
	CONN_PLAY_TYPE_LIVE      string = "live"
	CONN_PLAY_TYPE_PLAY_BACK string = "back"
)

type Connection struct {
	c           *gotcp.Conn
	recv_buffer *bytes.Buffer
	exit        chan struct{}
	conf        *Conf

	ID       string // for vavms sim_logicalchan_play_type
	SIM      string
	Channel  string
	PlayType string

	pipe_a *os.File
	pipe_v *os.File

	frame_audio base.Frame
	frame_vedio base.Frame

	once_prepare      sync.Once
	once_start_ffmpeg sync.Once
	ffmpeg_run        bool
	ffmpeg_cmd        string
}

func NewConnection(c *gotcp.Conn, conf *Conf) *Connection {
	tcp_c := c.GetRawConn()
	tcp_c.SetReadDeadline(time.Now().Add(conf.DefaultReadLimit))
	return &Connection{
		conf:        conf,
		c:           c,
		recv_buffer: bytes.NewBuffer([]byte{}),
		exit:        make(chan struct{}),
		ffmpeg_run:  true,
	}
}

func (c *Connection) OpenPipeA(pipe_a string) error {
	var err error
	c.pipe_a, err = os.OpenFile(pipe_a, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) OpenPipeV(pipe_v string) error {
	var err error
	c.pipe_v, err = os.OpenFile(pipe_v, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) SetProperty(sim, channel, play_type, cmd string) {
	c.SIM = sim
	c.Channel = channel
	c.PlayType = play_type
	c.ID = sim + "_" + channel + "_" + play_type
	c.ffmpeg_cmd = cmd
}

func (c *Connection) SetReadDeadline(minutes int) {
	c.c.GetRawConn().SetReadDeadline(time.Now().Add(time.Duration(minutes) * time.Minute))
}

func (c *Connection) Close() {
	close(c.exit)
	c.recv_buffer.Reset()
	if c.pipe_a != nil {
		c.pipe_a.Close()
	}
	if c.pipe_v != nil {
		c.pipe_v.Close()
	}
}
