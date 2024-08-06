package p2p

/*Peer represents the remote node*/
type Peer interface {
	Close() error
}

/*
	Transport is anything  that handles comms between the nodes in the network

It can be TCP , UDP , web sockets etc
*/
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
