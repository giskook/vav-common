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
	OnDataAudio(*Connection) bool
	OnDataVideo(*Connection) bool
	OnClose(*Connection) bool
}

type PrepareFunc func(string, string) error

type Connection struct {
	c           *gotcp.Conn
	recv_buffer *bytes.Buffer
	exit        chan struct{}
	status      uint8
	conf        *Conf

	id string // format servertype.sim.logicalchan

	pipe_aw     *os.File // current audio
	pipe_vw     *os.File // current video
	pipe_aw_his *os.File // history audio
	pipe_vw_his *os.File // history video
	pipe_ar     *os.File // for two way intercom

	frame_audio base.Frame
	frame_vedio base.Frame

	once_prepare sync.Once
	func_prepare PrepareFunc
}

func NewConnection(c *gotcp.Conn, conf *Conf, prepare_func PrepareFunc) *Connection {
	tcp_c := c.GetRawConn()
	tcp_c.SetReadDeadline(time.Now().Add(conf.DefautReadLimit))
	return &Connection{
		conf:         conf,
		c:            c,
		func_prepare: prepare_func,
		recv_buffer:  bytes.NewBuffer([]byte{}),
		exit:         make(chan struct{}),
	}
}

func (c *Connection) SetReadDeadline(minutes int) {
	c.c.GetRawConn().SetReadDeadline(time.Now().Add(time.Duration(minutes) * time.Minute))
}

func (c *Connection) Close() {
	close(c.exit)
	c.recv_buffer.Reset()
	if c.pipe_aw != nil {
		c.pipe_aw.Close()
	}
	if c.pipe_vw != nil {
		c.pipe_vw.Close()
	}
	if c.pipe_aw_his != nil {
		c.pipe_aw_his.Close()
	}
	if c.pipe_vw_his != nil {
		c.pipe_vw_his.Close()
	}
	if c.pipe_ar != nil {
		c.pipe_ar.Close()
	}
}
