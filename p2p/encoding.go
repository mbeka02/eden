package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	//temp fix
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct {
}
type NOPDecoder struct {
}

func (dec GOBDecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

func (dec NOPDecoder) Decode(r io.Reader, msg *Message) error {
	buff := make([]byte, 1024)
	n, err := r.Read(buff)
	if err != nil {
		return err
	}
	msg.Payload = buff[:n]
	return nil
}
