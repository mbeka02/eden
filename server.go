package main

import (
	//"fmt"
	//"log"
	//"net/http"

	"fmt"

	"github.com/mbeka02/eden/p2p"
)

type FileServerOpts struct {
	//ListenAddr        string
	StorageRoot       string
	PathTransformFunc pathTransformFunc
	Transport         p2p.Transport
}
type FileServer struct {
	FileServerOpts               //config options
	store          *store        //manages file storage on disk
	quitCh         chan struct{} //signal channel
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
		quitCh:         make(chan struct{}),
	}
}

func (f *FileServer) Run() error {

	if err := f.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}

func (f *FileServer) Loop() {

	for {
		select {
		case msg := <-f.Transport.Consume():
			fmt.Printf("Msg=>%v", msg)

		case <-f.quitCh:
			return
		}
	}
}
func (f *FileServer) Stop() {
	close(f.quitCh)
}

// cmds
/*func (f *FileServer) Store(key string, r io.Reader) error {

	return f.store.Write(key, r)

}*/
