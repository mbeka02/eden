package main

import (
	"fmt"
	"log"

	"github.com/mbeka02/eden/p2p"
)

func main() {
	//using the channel for blocking
	testChannel := make(chan string)
	transport := p2p.NewTCPTransport(":5173")
	fmt.Println("....starting")

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	<-testChannel
}
