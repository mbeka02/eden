package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
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
	// this is the underlying connection, in this case a TCP connection
	net.Conn
	// if we dial and retreive a conn => true
	// if we accept and retreive => false
	outboundPeer bool
	wg           *sync.WaitGroup
}

func NewTCPPeer(Conn net.Conn, outboundPeer bool) *TCPPeer {
	wg := &sync.WaitGroup{}
	return &TCPPeer{
		Conn,
		outboundPeer,
		wg,
	}
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcChan:          make(chan RPC, 1024),
	}
}

// Close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.Conn.Close()
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

// Send implements the Peer interface
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

// Addr() implements the Transport interface , returns the addr the transport is accepting connections
func (tr *TCPTransport) Addr() string {
	return tr.ListenAddr
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
	// handle the connection  in a separate go-routine

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

	// Read loop - decodes  messages and sends them to the server via the channel
	for {

		rpc := RPC{}

		err = tr.Decoder.Decode(conn, &rpc)

		if isNetConnClosedErr(err) {
			fmt.Printf("error the connection is  closed=>%v\n", err)
			return
		}
		if err != nil {

			fmt.Printf("tcp  read error: %v\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr().String()

		if rpc.Stream {
			peer.wg.Add(1)
			fmt.Println("... incoming  data stream from:", conn.RemoteAddr().String())
			peer.wg.Wait()
			fmt.Println("...done streaming the data , resuming the normal read loop")
			continue
		}

		tr.rpcChan <- rpc
	}
}
