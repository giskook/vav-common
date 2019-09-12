package socket_server

import (
	"github.com/giskook/vav-common/redis_cli"
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
		redis_cli.GetInstance().Sub(c.SIM, c.exit, c.Send)
	}()
}
