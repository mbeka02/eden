package main

import (
	"fmt"
	"log"

	"github.com/mbeka02/eden/p2p"
)

func OnPeer(peer p2p.Peer) error {
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
}
