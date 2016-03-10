package StreamServer

/// implement of interface IPacket on Server
type NameServerPacket struct {
}

func (self NameServerPacket) ToBytes() []byte {
	return []byte{}
}

func (self NameServerPacket) DispatchPacket(incoming []byte) (usedByte int, packet interface{}) {
	return -1, nil
}
