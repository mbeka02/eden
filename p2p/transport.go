package p2p

import "net"

/*Peer represents the remote node*/
type Peer interface {
	Send([]byte) error
	//RemoteAddr() net.Addr
	//Close() error
	net.Conn
}

/*
	Transport is anything  that handles comms between the nodes in the network

It can be TCP , UDP , web sockets etc
*/
type Transport interface {
	Addr() string
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
