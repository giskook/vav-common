package socket_server

import (
	"log"
	"testing"
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

func TestNewSocketServer(t *testing.T) {
	my := &myserver{}
	conf := &Conf{
		TcpAddr:    ":8876",
		FifoDir:    "/tmp/",
		FFmpegBin:  "/tmp/",
		ServerType: SERVER_TYPE_VAVMS,
	}

	my.server = NewSocketServer(conf, my)
	my.server.Start()
	for {
	}
}
