package main

import (
	//"fmt"
	"log"
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
	}
}
func (f *FileServer) BootstrapNetwork() {
	for _, addr := range f.BootStrapNodes {

		go func(addr string) {
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

// cmds
/*func (f *FileServer) Store(key string, r io.Reader) error {

	return f.store.Write(key, r)

}*/
