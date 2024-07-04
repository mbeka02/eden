package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransportOpts struct {
	ListenAddr  string
	HandshakeFn HandshakeFn
	Decoder     Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener // interface with 3 methods  Accept() , Close() and Addr().
	mux      sync.RWMutex //lock for the peers
	peers    map[net.Addr]Peer
}

// represents the remote node over an established TCP connection
type TCPPeer struct {
	//underlying connection
	conn net.Conn
	//if we dial and retreive a conn => true
	//if we accept and retreive => false
	outboundPeer bool
}
type Temp struct{}

func NewTCPPeer(conn net.Conn, outboundPeer bool) *TCPPeer {
	return &TCPPeer{
		conn,
		outboundPeer,
	}

}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (tr *TCPTransport) ListenAndAccept() error {
	var err error
	tr.listener, err = net.Listen("tcp", tr.ListenAddr)
	if err != nil {
		return err
	}

	go tr.startAcceptLoop()

	return nil
}

func (tr *TCPTransport) startAcceptLoop() {
	for {
		conn, err := tr.listener.Accept()
		if err != nil {
			fmt.Printf("TCP  accept() error : %s", err)
		}
		go tr.handleConnection(conn)
	}
}

// TO DO

func (tr *TCPTransport) handleConnection(conn net.Conn) {

	peer := NewTCPPeer(conn, true)
	err := tr.HandshakeFn(peer)
	if err != nil {
		conn.Close()
		fmt.Printf("tcp handshake error : %v\n", err)
		return
	}
	fmt.Printf(" New connection : %+v\n", peer)

	msg := &Temp{}
	//Read loop
	for {
		err = tr.Decoder.Decode(conn, msg)
		if err != nil {

			fmt.Printf("tcp error , unable to read incoming data : %v\n", err)
			continue
		}
	}
}
