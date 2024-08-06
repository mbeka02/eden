package main

import (
	"log"

	"github.com/mbeka02/eden/p2p"
	//"time"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	transportOpts := p2p.TCPTransportOpts{
		ListenAddr:  listenAddr,
		Decoder:     p2p.DefaultDecoder{},
		HandshakeFn: p2p.DefaultHandshakeFn,
		//TO DO
		//	OnPeer: func(p2p.Peer) error { return nil },
	}
	TCPTransport := p2p.NewTCPTransport(transportOpts)
	opts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASTransFunc,
		Transport:         TCPTransport,
		BootStrapNodes:    nodes,
	}
	return NewServer(opts)
}
func main() {
	nodes := []string{"172.0.0.1:3001", "172.0.0.1:3002", "172.0.0.1:3003", "172.0.0.1:3004"}
	fileServer := makeServer(":3000", nodes...)
	//test
	/*	go func() {
		time.Sleep(time.Second * 5)
		fileServer.Stop()
	}()*/

	if err := fileServer.Run(); err != nil {
		log.Fatalf("Unable to run the server : %v", err)
	}

}
