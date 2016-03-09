package StreamServer

import (
	"fmt"
	"net"
)

type ServerHandle interface {
	OnConnection(conn net.Conn)
	OnAcceptError(err error)
}

type TCPServer struct {
	addr    string
	handler ServerHandle
}

func NewTCPServer(addr string, handler ServerHandle) *TCPServer {
	if addr == "" {
		panic("Listen TCP addr can not be empty!")
	}
	if handler == nil {
		panic("IStreamServer interface instance can not be nil.")
	}
	s := new(TCPServer)
	s.addr = addr
	s.handler = handler
	return s
}

func (self *TCPServer) RunForever() {
	ln, err := net.Listen("tcp", self.addr)
	if err != nil {
		panic("TCPServer error on Listen")
	}
	fmt.Printf("Server Listen on TCP(%s)\n", self.addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			self.handler.OnAcceptError(err)
			continue
		}
		go self.handler.OnConnection(conn)
	}

}
