package main

import (
	//"fmt"
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"sync"
	"time"

	//"net/http"

	"fmt"

	"github.com/mbeka02/eden/p2p"
)

type FileServerOpts struct {
	//ListenAddr        string
	StorageRoot       string
	PathTransformFunc pathTransformFunc
	Transport         p2p.Transport
	BootStrapNodes    []string
}
type FileServer struct {
	FileServerOpts               // file server config options
	store          *store        //manages file storage on disk
	quitChannel    chan struct{} //signal channel

	peerLock sync.Mutex

	peers map[string]p2p.Peer
}

/*
	type DataMessage struct {
		Key  string
		Data []byte
	}
*/
type Message struct {
	Payload any
	//From    string
}

func NewServer(fileServerOptions FileServerOpts) *FileServer {

	storeOptions := storeOpts{
		pathTransformFunc: fileServerOptions.PathTransformFunc,
		root:              fileServerOptions.StorageRoot,
	}

	store := newStore(storeOptions)
	return &FileServer{
		FileServerOpts: fileServerOptions,
		store:          store,
		quitChannel:    make(chan struct{}),

		peers: make(map[string]p2p.Peer),
	}
}

func (f *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range f.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

func (f *FileServer) BootstrapNetwork() {
	for _, addr := range f.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			fmt.Println("attempting to connect with the remote node:", addr)
			if err := f.Transport.Dial(addr); err != nil {

				log.Printf("dial error : %v\n", err)

			}

		}(addr)

	}
}

func (f *FileServer) Run() error {

	if err := f.Transport.ListenAndAccept(); err != nil {
		return err
	}
	f.BootstrapNetwork()
	f.Loop()
	return nil
}

func (f *FileServer) Loop() {
	defer func() {
		log.Println("...exiting")
		f.Transport.Close()
	}()

	for {
		select {
		case rpc := <-f.Transport.Consume():

			var message Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&message); err != nil {
				log.Printf("decoding error:%v\n", err)
			}
			/*	if err := f.handleMessage(&message); err != nil {
				log.Printf("unable to handle the message:%v", err)

			}*/
			log.Printf("received => %s\n", string(message.Payload.([]byte)))
			peer, ok := f.peers[rpc.From]
			if !ok {
				panic("peer not found")
			}
			b := make([]byte, 1024)
			_, err := peer.Read(b)
			if err != nil {
				panic(err)
			}
			fmt.Printf("printing data: %s\n", string(b))
			peer.(*p2p.TCPPeer).Wg.Done()

		case <-f.quitChannel:
			return
		}
	}
}

/*func (f *FileServer) handleMessage(msg *Message) error {
	switch v := msg.Payload.(type) {
	case *DataMessage:
		fmt.Printf("received data %+v\n", v)
	}
	return nil
}*/

func (f *FileServer) Stop() {
	close(f.quitChannel)
}
func (f *FileServer) OnPeer(p p2p.Peer) error {
	f.peerLock.Lock()
	defer f.peerLock.Unlock()
	//add peer to the map
	f.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote=> %s\n", p.RemoteAddr())
	return nil
}

func (f *FileServer) StoreData(key string, r io.Reader) error {
	//send the message type
	//send the actual payload/message
	for _, peer := range f.peers {
		buff := new(bytes.Buffer)

		msg := Message{[]byte("this is a random message")}
		if err := gob.NewEncoder(buff).Encode(&msg); err != nil {
			fmt.Println(err)
		}

		if err := peer.Send(buff.Bytes()); err != nil {
			fmt.Println(" peer send error=>", err)
		}
	}
	time.Sleep(time.Second * 3)
	payload := []byte("LARGE PAYLOAD")
	for _, peer := range f.peers {
		err := peer.Send(payload)

		if err != nil {
			log.Println("error ,unable to send the payload:", payload)
			//	continue
		}
	}
	/*buff := new(bytes.Buffer)
	teeReader := io.TeeReader(r, buff)

	if err := f.store.Write(key, teeReader); err != nil {
		return err
	}
	payload := &DataMessage{

		Key:  key,
		Data: buff.Bytes(),
	}
	//broadcast the payload to all other remote peers
	return f.broadcast(&Message{
		From:    "todo",
		Payload: payload,
	})*/
	return nil
}
