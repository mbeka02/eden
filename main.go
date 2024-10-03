package main

import (
	"bytes"
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

	data := bytes.NewReader([]byte(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="utf-8">
	<title>My test page</title>
	<link href="http://fonts.googleapis.com/css?family=Open+Sans" rel="stylesheet" type="text/css">
	<link href="styles/style.css" rel="stylesheet" type="text/css">
	</head>
	<body>
	<h1>Mozilla is cool</h1>
	<img src="images/firefox-icon.png" alt="The Firefox logo: a flaming fox surrounding the Earth.">
	<p>At Mozilla, we’re a global community of</p>
	<ul> <!-- changed to list in the tutorial -->
	<li>technologists</li>
	<li>thinkers</li>
	<li>builders</li>
	</ul>
	<p>working together to keep the Internet alive and accessible, so people worldwide can be informed contributors and creators of the Web. We believe this act of human collaboration across an open platform is essential to individual growth and our collective future.</p>
	</body>
	</html>
	  `))

	//fileServer1.Store("myprivateDataKey", data)
	time.Sleep(time.Second * 2)

	fileServer2.Store("myprivateDataKey", data)
	//r, err := fileServer2.Get("myprivateDataKey")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// b, err := io.ReadAll(r)
	//
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("buffer content=>", string(b))
	select {}
}
