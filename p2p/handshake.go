package p2p

//

type HandshakeFn func(Peer) error

func DefaultHandshakeFn(Peer) error {
	return nil
}
