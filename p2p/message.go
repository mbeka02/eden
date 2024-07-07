package p2p

import "net"

/*
	Message represents represents any arbritary data

that is being sent between two nodes in the network
*/
type Message struct {
	From    net.Addr
	Payload []byte
}
