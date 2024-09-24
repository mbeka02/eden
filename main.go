package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/mbeka02/eden/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	transportOpts := p2p.TCPTransportOpts{
		ListenAddr:  listenAddr,
		Decoder:     p2p.DefaultDecoder{},
		HandshakeFn: p2p.DefaultHandshakeFn,
	}
	TCPTransport := p2p.NewTCPTransport(transportOpts)
	opts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASTransFunc,
		Transport:         TCPTransport,
		BootStrapNodes:    nodes,
	}
	server := NewServer(opts)

	TCPTransport.OnPeer = server.OnPeer
	return server
}
func main() {
	fileServer1 := makeServer(":3000", "")
	fileServer2 := makeServer(":4000", "127.0.0.1:3000")
	go func() {
		log.Fatal(fileServer1.Run())
	}()
	time.Sleep(time.Second * 2)

	go fileServer2.Run()
	time.Sleep(time.Second * 4)

	data := bytes.NewReader([]byte("some random data"))
	fileServer2.Store("myprivateDataKey", data)
	r, err := fileServer2.Get("myprivateDataKey")

	if err != nil {
		fmt.Println(err)
	}
	buff := make([]byte, 1024)
	n, err := r.Read(buff)

	if err != nil {
		fmt.Println(err)
	}
	content := buff[:n]
	fmt.Println("buffer content=>", string(content))

	select {}
}
