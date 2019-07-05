package socket_server

import (
	"errors"
	"github.com/gansidui/gotcp"
	"log"
	"net"
)

var (
	ErrPeerClosed = errors.New("peer closed")
)

type Raw struct {
	raw []byte
}

func (r *Raw) Serialize() []byte {
	return r.raw
}

func (ss *SocketServer) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	data := make([]byte, 4096)
	length, err := conn.Read(data)
	if err != nil {
		log.Printf("<ERR> %s\n", err.Error())
		return nil, err
	}

	if length == 0 {
		log.Println("<ERR> peer error")
		return nil, ErrPeerClosed
	}

	return &Raw{
		raw: data[0:length],
	}, nil
}
