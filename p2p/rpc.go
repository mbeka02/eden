package p2p

/*
RPC holds any arbritary data
that is being sent
over each  transport
between two nodes in the network
*/
type RPC struct {
	From    string
	Payload []byte
}
