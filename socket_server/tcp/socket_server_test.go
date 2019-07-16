package socket_server

import (
	"log"
	"testing"
	"time"
)

type myserver struct {
	server *SocketServer
}

func (s *myserver) OnDataAudio(conn *Connection) bool {
	log.Println("OnDataAudio")
	return true
}
func (s *myserver) OnDataVideo(conn *Connection) bool {
	log.Println("OnDataVideo")
	return true
}
func (s *myserver) OnClose(conn *Connection) bool {
	log.Println("TestOnClose")
	return true
}

func (s *myserver) Prepare(conn *Connection, sim, channel string) error {
	log.Println("TestPrepare", conn, sim, channel)
	return nil
}

func TestNewSocketServer(t *testing.T) {
	my := &myserver{}
	conf := &Conf{
		TcpAddr:         ":8876",
		ServerType:      SERVER_TYPE_VAVMS,
		DefautReadLimit: time.Duration(1) * time.Minute,
	}

	my.server = NewSocketServer(conf, my, my.Prepare)
	my.server.Start()
	for {
	}
}
