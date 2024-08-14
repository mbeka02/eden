package main

import (
	//"fmt"
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"sync"

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
	FileServerOpts               //config options
	store          *store        //manages file storage on disk
	quitChannel    chan struct{} //signal channel

	peerLock sync.Mutex

	peers map[string]p2p.Peer
}

type Payload struct {
	key  string
	data []byte
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

// TODO
func (f *FileServer) broadcast(p *Payload) error {
	peers := []io.Writer{}

	for _, peer := range f.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

func (f *FileServer) BootstrapNetwork() {
	for _, addr := range f.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			fmt.Println("attempting to connect with remote:", addr)
			if err := f.Transport.Dial(addr); err != nil {

				log.Printf("dial error : %v", err)

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
		case msg := <-f.Transport.Consume():
			fmt.Printf("Msg=>%v\n", msg)

		case <-f.quitChannel:
			return
		}
	}
}
func (f *FileServer) Stop() {
	close(f.quitChannel)
}
func (f *FileServer) OnPeer(p p2p.Peer) error {
	f.peerLock.Lock()
	defer f.peerLock.Unlock()
	f.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote=> %s", p.RemoteAddr())
	return nil
}

// cmds
/*func (f *FileServer) Store(key string, r io.Reader) error {

	return f.store.Write(key, r)

}*/
func (f *FileServer) StoreData(key string, r io.Reader) error {
	if err := f.store.Write(key, r); err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	if _, err := io.Copy(buff, r); err != nil {
		return err
	}
	fmt.Println("bytes=>", buff.Bytes())
	payload := &Payload{
		data: buff.Bytes(),
		key:  key,
	}
	return f.broadcast(payload)
}
