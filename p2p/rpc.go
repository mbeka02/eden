package p2p

const (
	IncomingStream  = 0x2
	IncomingMessage = 0x1
)

/*
RPC holds any arbritary data
that is being sent
over each  transport
between two nodes in the network
*/
type RPC struct {
	From    string //network addr of a node
	Payload []byte //the data
	Stream  bool   //data streaming flag , indicates whether data streaming is active or not
}
