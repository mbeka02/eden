package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPTransportOpts struct {
	ListenAddr  string
	HandshakeFn HandshakeFn
	Decoder     Decoder
	OnPeer      func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener // interface with 3 methods  Accept() , Close() and Addr().
	rpcChan  chan RPC
}

// represents the remote node over an established TCP connection
type TCPPeer struct {
	//underlying connection
	conn net.Conn
	//if we dial and retreive a conn => true
	//if we accept and retreive => false
	outboundPeer bool
}

func NewTCPPeer(conn net.Conn, outboundPeer bool) *TCPPeer {
	return &TCPPeer{
		conn,
		outboundPeer,
	}

}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcChan:          make(chan RPC),
	}
}

// Close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// RemoteAddr implements the Peer interface and will return the
// remote addr of its underlying connection
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

// Send implements the Peer interface
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}
func (tr *TCPTransport) ListenAndAccept() error {
	var err error
	tr.listener, err = net.Listen("tcp", tr.ListenAddr)
	if err != nil {
		return err
	}

	go tr.startAcceptLoop()
	log.Printf("TCP transport is listening on port %s\n", tr.ListenAddr)
	return nil
}

/*
Consume implements the Transport interface , which will return a
read only channel for reading the read-only
messages received from another peer on the network
*/
func (tr *TCPTransport) Consume() <-chan RPC {

	return tr.rpcChan

}

// Close implements the transport interface
func (tr *TCPTransport) Close() error {
	return tr.listener.Close()
}
func (tr *TCPTransport) startAcceptLoop() {
	for {
		conn, err := tr.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP  accept() error : %s", err)
		}
		go tr.handleConnection(conn, false)
	}
}

// Dial implements the transport interface
func (tr *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go tr.handleConnection(conn, true)
	return nil
}
func (tr *TCPTransport) handleConnection(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection : %v\n", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, outbound)
	err = tr.HandshakeFn(peer)
	if err != nil {
		return
	}
	if tr.OnPeer != nil {
		if err := tr.OnPeer(peer); err != nil {
			return
		}

	}

	//Read loop
	rpc := RPC{}
	for {
		err = tr.Decoder.Decode(conn, &rpc)
		if isNetConnClosedErr(err) {
			fmt.Printf("error the connection is  closed=>%v\n", err)
			return
		}
		if err != nil {

			fmt.Printf("tcp  read error: %v\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		tr.rpcChan <- rpc

	}
}
