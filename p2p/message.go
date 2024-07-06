package p2p

/*
	Message represents represents any arbritary data

that is being sent between two nodes in the network
*/
type Message struct {
	Payload []byte
}
