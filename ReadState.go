package goxdr

type ReadState interface {
	Update([]byte) (int, bool)
	EndPacket() error
}

type RequestReadState interface {
	ReadState
	ResponsePacket() Packet
}
