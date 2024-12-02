package main

import (
	"fmt"
	"io"
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

	// fileServer1.Store("myprivateDataKey", data)
	time.Sleep(time.Second * 2)

	// data := bytes.NewReader([]byte("random data"))
	//
	// fileServer2.Store("coolPicture.jpg", data)
	// time.Sleep(time.Millisecond * 500)

	r, err := fileServer2.Get("coolPicture.jpg")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("buffer content=>", string(b))
}
