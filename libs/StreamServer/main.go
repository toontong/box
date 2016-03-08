package main

import (
	"net"

	"github.com/toontong/box/libs/StreamServer"
	"github.com/toontong/box/libs/log"
)

type NameServer struct {
}

func (self *NameServer) OnConnection(conn net.Conn) {
	log.Infof("New connect to NameServer[%s] From RemoteAddr[%s].",
		conn.LocalAddr(), conn.RemoteAddr())

	c := NewConnection(conn, Default_Timeout)

	c.EventLoop()
	log.Infof("End OnConnection of [%s]", conn.RemoteAddr())
}

func (self NameServer) OnAcceptError(err error) {
	log.Infof("NameServer OnAcceptError=[%s]", err)
}

func main() {
	nSvr := new(NameServer)
	tcpSvr := StreamServer.NewTCPServer(":20000", nSvr)
	tcpSvr.RunForever()
}
