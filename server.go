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

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64 //file size
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
			if err := f.handleMessage(rpc.From, &message); err != nil {

				log.Println(err)

			}
		case <-f.quitChannel:
			return
		}
	}
}

func (f *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return f.handleMessageStoreFile(from, v)
	}
	return nil
}

func (f *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	//check if the peer exists in the map
	peer, ok := f.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) does not exist", from)
	}

	n, err := f.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}
	log.Printf("Written (%d) bytes to disk \n", n)

	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

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

// StoreData stores the file on disk then broadcasts it to all the other nodes
func (f *FileServer) StoreData(key string, r io.Reader) error {
	var (
		fileBuffer = new(bytes.Buffer)
		teeReader  = io.TeeReader(r, fileBuffer)
	)

	size, err := f.store.Write(key, teeReader)
	if err != nil {
		return err
	}
	fmt.Printf("bytes written=>%v\n", size)
	msg := Message{Payload: MessageStoreFile{
		Key:  key,
		Size: size,
	}}
	messageBuffer := new(bytes.Buffer)
	if err := gob.NewEncoder(messageBuffer).Encode(&msg); err != nil {
		fmt.Println(err)

	}

	for _, peer := range f.peers {
		if err := peer.Send(messageBuffer.Bytes()); err != nil {
			fmt.Println("peer send error=>", err)
		}
	}
	time.Sleep(time.Second * 2)
	for _, peer := range f.peers {
		n, err := io.Copy(peer, fileBuffer)

		if err != nil {
			return err
		}
		fmt.Printf("received and written %v bytes to disk\n", n)
	}
	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
}
