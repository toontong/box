package StreamServer

/// implement of interface IPacket on Server
type NameServerPacket struct {
}

func (self NameServerPacket) ToBytes() []byte {

}

func (self NameServerPacket) DispatchPacket(incoming []byte) (usedByte int, packet interface{}) {

}
