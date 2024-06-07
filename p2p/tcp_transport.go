package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr string
	listener   net.Listener // interface with 3 methods  Accept() , Close() and Addr().
	mux        sync.RWMutex //lock for the peers
	peers      map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{

		listenAddr: listenAddr,
	}
}

func (tr *TCPTransport) ListenAndAccept() error {
	var err error
	tr.listener, err = net.Listen("tcp", tr.listenAddr)
	if err != nil {
		return err
	}

	go tr.startAcceptLoop()

	return nil
}

func (tr *TCPTransport) startAcceptLoop() {
	for {
		conn, err := tr.listener.Accept()
		tr.handleConnection(conn)
		if err != nil {
			fmt.Printf("TCP  accept() error : %s", err)
		}

	}
}

func (tr *TCPTransport) handleConnection(conn net.Conn) {

}
