package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr  string
	listener    net.Listener // interface with 3 methods  Accept() , Close() and Addr().
	mux         sync.RWMutex //lock for the peers
	handshakeFn HandshakeFn
	decoder     Decoder
	peers       map[net.Addr]Peer
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

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		//placeholder
		handshakeFn: DefaultHandshakeFn,
		listenAddr:  listenAddr,
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
		if err != nil {
			fmt.Printf("TCP  accept() error : %s", err)
		}
		go tr.handleConnection(conn)
	}
}

// TO DO

func (tr *TCPTransport) handleConnection(conn net.Conn) {

	peer := NewTCPPeer(conn, true)
	err := tr.handshakeFn(peer)
	if err != nil {

	}
	fmt.Printf(" New connection : %+v\n", peer)

	msg := &Temp{}
	lenDecodeErr := 0
	//Read loop
	for {
		err = tr.decoder.Decode(conn, msg)
		if err != nil {
			if lenDecodeErr == 6 {

				log.Fatalf("tcp error : unable to read incoming data : %v", err)
			}
		}
	}
}
