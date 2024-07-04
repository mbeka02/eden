package main

import (
	"fmt"
	"log"

	"github.com/mbeka02/eden/p2p"
)

func main() {
	//using the channel for blocking
	testChannel := make(chan string)
	opts := p2p.TCPTransportOpts{
		ListenAddr:  ":5173",
		Decoder:     p2p.GOBDecoder{},
		HandshakeFn: p2p.DefaultHandshakeFn,
	}
	transport := p2p.NewTCPTransport(opts)
	fmt.Println("....starting")

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	<-testChannel
}
