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
	OnClose(*Connection) bool
}

type Connection struct {
	c           *gotcp.Conn
	recv_buffer *bytes.Buffer
	exit        chan struct{}
	status      uint8
	conf        *Conf

	id string // format servertype.sim.logicalchan

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

func (c *Connection) SetFfmpegCmd(cmd string) {
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
