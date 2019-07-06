package socket_server

import (
	"github.com/gansidui/gotcp"
	"log"
	"net"
	"time"
)

const (
	SERVER_TYPE_VAVMS string = "v"
	SERVER_TYPE_TWIS  string = "a"
)

type Conf struct {
	TcpAddr         string
	ServerType      string
	DefautReadLimit time.Duration
}

type SocketServer struct {
	conf         *Conf
	srv          *gotcp.Server
	cm           *ConnMgr
	callback     ConnCallback
	func_prepare PrepareFunc
	exit         chan struct{}
}

func NewSocketServer(conf *Conf, callback ConnCallback, prepare_func PrepareFunc) *SocketServer {
	return &SocketServer{
		conf:         conf,
		callback:     callback,
		func_prepare: prepare_func,
		cm:           NewConnMgr(),
		exit:         make(chan struct{}),
	}
}

func (ss *SocketServer) Start() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ss.conf.TcpAddr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}

	ss.srv = gotcp.NewServer(config, ss, ss)

	go ss.srv.Start(listener, time.Second)
	log.Println("<INFO> socket listening:", listener.Addr())

	return nil
}

func (ss *SocketServer) Stop() {
	close(ss.exit)
	ss.srv.Stop()
}
