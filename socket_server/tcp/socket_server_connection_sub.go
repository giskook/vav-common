package socket_server

import (
	"bytes"
	"github.com/gansidui/gotcp"
	"github.com/giskook/vav-common/base"
	"github.com/giskook/vav-common/redis_cli"
	"os"
	"sync"
	"time"
)

type raw struct {
	raw []byte
}

func (r *raw) Serialize() []byte {
	return r.raw
}

func (c *Connection) Send(content []byte) error {
	return c.c.AsyncWritePacket(&raw{
		raw: content,
	}, 0)
}

func (c *Connection) Sub() {
	go func() {
		redis_cli.GetInstance().Sub(c.sim, c.exit, c.Send)
	}()
}
