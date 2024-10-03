package p2p

import (
	"encoding/gob"
	"fmt"
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

	peekBuff := make([]byte, 1)
	if _, err := r.Read(peekBuff); err != nil {
		return fmt.Errorf("peeking error: %v\n", err)
	}

	stream := peekBuff[0] == IncomingStream
	// don't decode if it's a stream
	if stream {
		m.Stream = true
		return nil
	}

	buff := make([]byte, 1024)
	n, err := r.Read(buff)
	if err != nil {
		return err
	}
	m.Payload = buff[:n]
	return nil
}
