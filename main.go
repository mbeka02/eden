package main

import (
	"log"

	"github.com/mbeka02/eden/p2p"
)

//	"fmt"
//
// "log"
// "github.com/mbeka02/eden/p2p"

/*func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("....on peer logic")
	return nil
}

func main() {
	blockingChannel := make(chan string)
	opts := p2p.TCPTransportOpts{
		ListenAddr:  ":5173",
		Decoder:     p2p.DefaultDecoder{},
		HandshakeFn: p2p.DefaultHandshakeFn,
		OnPeer:      OnPeer,
	}
	transport := p2p.NewTCPTransport(opts)
	fmt.Println("....starting service")

	go func() {
		for {
			msg := <-transport.Consume()
			fmt.Printf("message=>%v\n", msg)
		}
	}()

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	<-blockingChannel
}*/

func main() {
	transportOpts := p2p.TCPTransportOpts{
		ListenAddr:  ":3000",
		Decoder:     p2p.DefaultDecoder{},
		HandshakeFn: p2p.DefaultHandshakeFn,
		//TO DO
		//	OnPeer: func(p2p.Peer) error { return nil },
	}
	TCPTransport := p2p.NewTCPTransport(transportOpts)
	opts := FileServerOpts{
		StorageRoot:       "home",
		PathTransformFunc: CASTransFunc,
		Transport:         TCPTransport,
	}
	fileServer := NewServer(opts)
	if err := fileServer.Run(); err != nil {
		log.Fatalf("Unable to run the server : %v", err)
	}
	select {}
}
