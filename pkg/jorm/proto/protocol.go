package proto

import (
	"io"
)

type Packet struct {
	Chksum  int32
	Size    int32
	Payload []byte
}

func ReadPacket(stream io.Reader) (*Packet, error) {
	// TODO
	return nil, nil
}
