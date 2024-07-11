package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:  ":5173",
		Decoder:     DefaultDecoder{},
		HandshakeFn: DefaultHandshakeFn,
	}
	transport := NewTCPTransport(opts)
	assert.Equal(t, transport.ListenAddr, opts.ListenAddr)
	//assert.Equal(t, transport.Decoder, opts.Decoder)

	//server
	err := transport.ListenAndAccept()
	assert.NoError(t, err)

}
