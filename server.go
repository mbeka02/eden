package main

import (
	"fmt"
	"net/http"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc pathTransformFunc
}
type FileServer struct {
	FileServerOpts
	store *store
}

func NewServer(fileServerOptions FileServerOpts) *FileServer {

	storeOptions := storeOpts{
		pathTransformFunc: CASTransFunc,
		root:              fileServerOptions.StorageRoot,
	}

	store := newStore(storeOptions)
	return &FileServer{
		FileServerOpts: fileServerOptions,
		store:          store,
	}
}

func (f *FileServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test route")
	})
	fmt.Printf("Server is listening on port%s", f.ListenAddr)
	return http.ListenAndServe(f.ListenAddr, router)

}
