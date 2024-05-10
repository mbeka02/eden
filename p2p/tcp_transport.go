package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr string
	listener   net.Listener // interface with 3 methods  Accept() , Close() and Addr().
	mux        sync.RWMutex //lock for the peers
	peers      map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) Transport {
	return &TCPTransport{

		listenAddr: listenAddr,
	}
}
