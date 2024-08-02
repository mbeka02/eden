package main

import (
	//"fmt"
	//"log"
	//"net/http"

	"github.com/mbeka02/eden/p2p"
)

type FileServerOpts struct {
	//ListenAddr        string
	StorageRoot       string
	PathTransformFunc pathTransformFunc
	Transport         p2p.Transport
}
type FileServer struct {
	FileServerOpts
	store *store
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
	}
}

func (f *FileServer) Run() error {
	/*router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "test route")
	})
	fmt.Printf("Server is listening on port%s", f.ListenAddr)
	if err := http.ListenAndServe(f.ListenAddr, router); err != nil {
		log.Fatalf("Unabel to start the server:%v", err)
	}
	*/

	if err := f.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}
