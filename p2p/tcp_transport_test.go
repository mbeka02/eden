package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":8080"
	transport := NewTCPTransport(listenAddr)
	assert.Equal(t, transport.listenAddr, listenAddr)
}
