package socket_server

import (
	"bytes"
	"github.com/gansidui/gotcp"
	"github.com/giskook/vav-common/base"
	"os"
	"sync"
	"time"
)

const (
	VIDEO_TYPE int = 1
	AUDIO_TYPE int = 2
	AV_TYPE    int = 3
)

type ConnCallback interface {
	OnPrepare(*Connection, string, string) error
	OnFfmpegExit(*Connection) error
	OnClose(*Connection) error
}

type Connection struct {
	c           *gotcp.Conn
	recv_buffer *bytes.Buffer
	exit        chan struct{}
	conf        *Conf

	ID        string // for vavms sim_logicalchan_play_type
	SIM       string
	Channel   string
	PlayType  string
	TTL       int
	TimeStamp int

	pipe_a *os.File
	pipe_v *os.File

	frame_audio base.Frame
	frame_vedio base.Frame

	once_prepare      sync.Once
	once_start_ffmpeg sync.Once
	ffmpeg_run        bool
	ffmpeg_cmds       []string
	file_path_a       string
	file_path_v       string
	acodec            string
	vcodec            string
	avtype            int // 1 for v 2 for a 3 for av

	ffmpeg_name     string
	log_audio       bool
	log_video       bool
	baseline_time   int64
	elapsed_time    int64
	video_last_time int64
	audio_last_time int64
}

func NewConnection(c *gotcp.Conn, conf *Conf) *Connection {
	tcp_c := c.GetRawConn()
	tcp_c.SetReadDeadline(time.Now().Add(conf.DefaultReadLimit))
	return &Connection{
		conf:          conf,
		c:             c,
		recv_buffer:   bytes.NewBuffer([]byte{}),
		exit:          make(chan struct{}),
		ffmpeg_run:    true,
		TimeStamp:     int(time.Now().Unix()),
		log_audio:     false,
		log_video:     false,
		baseline_time: time.Now().Unix(),
		elapsed_time:  time.Now().Unix(),
		avtype:        0,
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

func (c *Connection) SetProperty(sim, channel, play_type, file_path_a, file_path_v, acodec, vcodec, ffmpeg_name string, cmds []string, ttl int) {
	c.SIM = sim
	c.Channel = channel
	c.PlayType = play_type
	c.ID = sim + "_" + channel + "_" + play_type
	c.ffmpeg_cmds = cmds
	c.file_path_a = file_path_a
	c.file_path_v = file_path_v
	c.acodec = acodec
	c.vcodec = vcodec
	c.ffmpeg_name = ffmpeg_name
	c.TTL = ttl
}

func (c *Connection) SetReadDeadline(seconds int) {
	c.c.GetRawConn().SetReadDeadline(time.Now().Add(time.Duration(seconds) * time.Second))
}

func (c *Connection) Close() {
	c.c.Close()
}

func (c *Connection) ShutDown() {
	close(c.exit)
	c.recv_buffer.Reset()
	if c.pipe_a != nil {
		time.Sleep(1 * time.Second) // wait the ffmpeg to start.
		c.pipe_a.Close()
	}
	if c.pipe_v != nil {
		time.Sleep(1 * time.Second) // wait the ffmpeg to start.
		c.pipe_v.Close()
	}
}
