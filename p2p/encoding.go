package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct {
}
type DefaultDecoder struct {
}

func (dec GOBDecoder) Decode(r io.Reader, m *RPC) error {
	return gob.NewDecoder(r).Decode(m)
}

func (dec DefaultDecoder) Decode(r io.Reader, m *RPC) error {
	buff := make([]byte, 1024)
	n, err := r.Read(buff)
	if err != nil {
		return err
	}
	m.Payload = buff[:n]
	return nil
}
