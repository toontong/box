package StreamServer

import (
	"bytes"
	"net"
	"time"

	"github.com/toontong/box/libs/log"
)

const (
	Default_Timeout   = time.Second * 30
	Max_Read_Buff_Len = 1024
)

type IPacket interface {
	ToBytes() []byte
	// 对输入数据进行反序列化成某Pakcet，成功返回使用了字节数
	DispatchPacket(incoming []byte) (usedByte int)
}

type Connection struct {
	conn    net.Conn
	recvBuf chan []byte
	sendBuf chan []byte

	_remain_mem_recvBuf bytes.Buffer //未被unpack的数据

	closed bool

	closeChan chan bool

	timeout time.Duration

	packet IPacket
}

func NewConnection(c net.Conn, packet IPacket, timeout time.Duration) *Connection {
	conn := new(Connection)
	conn.conn = c

	conn.timeout = timeout

	conn.closed = false
	conn.closeChan = make(chan bool, 2)

	conn.sendBuf = make(chan []byte, 32)
	conn.recvBuf = make(chan []byte, 32)

	conn.packet = packet
	return conn
}

func (self *Connection) WritePacket(pkg IPacket) {
	self.sendBuf <- pkg.ToBytes()
}

func (self *Connection) EventLoop() {
	go self.readLoop()

	for {
		select {
		case closed := <-self.closeChan:
			self.closed = closed
			return
		case bufRecv := <-self.sendBuf:
			self.onRecv(bufRecv)
		case bufSend := <-self.sendBuf:
			self.write(bufSend)
		}
	}
}

func (self *Connection) readLoop() {

	buffer := make([]byte, Max_Read_Buff_Len)

	for !self.closed {

		self.conn.SetReadDeadline(time.Now().Add(self.timeout))

		n, err := self.conn.Read(buffer)
		if err != nil {
			log.Infof("Connection[%s] Read error [%s], will close it.",
				self.conn.RemoteAddr(), err)
			self.closeChan <- true
			return
		}

		self.recvBuf <- buffer[:n]
	}
}

func (self *Connection) write(buf []byte) {
	self.conn.SetWriteDeadline(time.Now().Add(self.timeout))

	n, err := self.conn.Write(buf)
	if err != nil || n != len(buf) {
		log.Infof("Connection[%s] Write error [%s]",
			self.conn.RemoteAddr(), err)
		self.closeChan <- true
		return
	}

}

func (self *Connection) onRecv(buf []byte) {
	var buffer []byte
	remain := self._remain_mem_recvBuf.Len()
	now := len(buf)
	if remain > 0 {
		buffer = make([]byte, remain+now)
		self._remain_mem_recvBuf.Read(buffer)
		for i := 0; i < now; i++ {
			buffer[i+remain] = buf[i]
		}
	} else {
		buffer = buf
	}

	n := self.packet.DispatchPacket(buffer)
	if n > 0 { // must > 0
		if n < len(buffer) {
			self._remain_mem_recvBuf.Write(buffer[n:])
		}

	} else {
		self._remain_mem_recvBuf.Write(buffer)
	}
}
