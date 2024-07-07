package main

import (
	"fmt"
	"log"

	"github.com/mbeka02/eden/p2p"
)

func main() {
	blockingChannel := make(chan string)
	opts := p2p.TCPTransportOpts{
		ListenAddr:  ":5173",
		Decoder:     p2p.DefaultDecoder{},
		HandshakeFn: p2p.DefaultHandshakeFn,
	}
	transport := p2p.NewTCPTransport(opts)
	fmt.Println("....starting application")

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	<-blockingChannel
}
