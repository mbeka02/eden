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
type MessageGetFile struct {
	Key string
}

//ensures type registration happens only once

var registerOnce sync.Once

func registerTypes() {
	registerOnce.Do(func() {
		gob.Register(MessageStoreFile{})
		gob.Register(MessageGetFile{})
	})
}
func NewServer(fileServerOptions FileServerOpts) *FileServer {
	// register types when creating a new server
	// registerTypes()
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

func (f *FileServer) stream(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range f.peers {
		peers = append(peers, peer)
	}
	//duplicate the write to all the writers/peers
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

// broadcast encodes and sends the message to all the peers
func (f *FileServer) broadcast(msg *Message) error {
	messageBuffer := new(bytes.Buffer)
	if err := gob.NewEncoder(messageBuffer).Encode(msg); err != nil {
		return err
	}

	for _, peer := range f.peers {

		err := peer.Send([]byte{p2p.IncomingMessage})
		if err != nil {
			fmt.Println(err)
		}

		if err := peer.Send(messageBuffer.Bytes()); err != nil {
			fmt.Println("peer send error=>", err)
		}
	}

	return nil
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
		log.Println("...stopping the file server")
		f.Transport.Close()
	}()

	for {
		select {
		case rpc := <-f.Transport.Consume():

			var message Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&message); err != nil {
				log.Println("decoding error:", err)
			}
			if err := f.handleMessage(rpc.From, &message); err != nil {

				log.Println("handle Message error:", err)

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
	case MessageGetFile:
		return f.handleMessageGetFile(from, v)

	}
	return nil
}

func (f *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {

	if !f.store.Has(msg.Key) {
		return fmt.Errorf(" need to serve file %s but it  does not exist on disk\n", msg.Key)
	}

	r, err := f.store.Read(msg.Key)
	if err != nil {
		return err
	}
	//check if the peer exists in the map
	peer, ok := f.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) does not exist", from)
	}
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}
	log.Printf("written (%d) bytes over the network to %s", n, from)
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
	log.Printf(" [%s],Written (%d) bytes to disk   \n", f.Transport.Addr(), n)

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

func (f *FileServer) Get(key string) (io.Reader, error) {
	if f.store.Has(key) {
		return f.store.Read(key)
	}

	fmt.Printf("%s is not stored locally, fetching from the network\n", key)

	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	err := f.broadcast(&msg)
	if err != nil {
		return nil, err
	}
	for _, peer := range f.peers {
		fmt.Println("INSIDE THE LOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOP")
		fileBuffer := new(bytes.Buffer)
		_, err := io.Copy(fileBuffer, peer)
		fmt.Println("checkpoint 1")
		if err != nil {
			return nil, err
		}
		fmt.Println("checkpoint 2")
		fmt.Println("received:", string(fileBuffer.Bytes()))
	}
	select {}
	// r := bytes.NewReader([]byte{})
	//return nil , nil
}

// This method stores the file on disk then broadcasts it to all the other nodes
func (f *FileServer) Store(key string, r io.Reader) error {
	var (
		fileBuffer = new(bytes.Buffer)
		teeReader  = io.TeeReader(r, fileBuffer)
	)
	// store the file on disk (locally)
	size, err := f.store.Write(key, teeReader)
	if err != nil {
		return err
	}
	fmt.Printf("bytes written to disk locally=>%v\n", size)

	// broadcast a 'store file' message with the key and size of the payload
	msg := Message{Payload: MessageStoreFile{
		Key:  key,
		Size: size,
	}}
	if err := f.broadcast(&msg); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 500)
	// send and write the file to all the peers in the list

	peers := []io.Writer{}
	for _, peer := range f.peers {
		//send the peek buffer that indicates a stream is incoming
		peer.Send([]byte{p2p.IncomingStream})

		peers = append(peers, peer)
		/*n, err := io.Copy(peer, fileBuffer)

		if err != nil {
			return err
		}
		fmt.Printf("received and written %v bytes to disk\n", n)
		*/
	}
	mw := io.MultiWriter(peers...)
	mw.Write(fileBuffer.Bytes())
	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
	// gob.Register(MessageGetFile{})
}
