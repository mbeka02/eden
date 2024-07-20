package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type pathTransformFunc func(string) PathKey
type storeOpts struct {
	pathTransformFunc pathTransformFunc
}
type store struct {
	storeOpts
}
type PathKey struct {
	PathName string
	Filename string
}

func CASTransFunc(key string) PathKey {

	hash := sha1.Sum([]byte(key))

	hashString := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashString) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashString,
	}
}

func (pk PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", pk.PathName, pk.Filename)

}
func defaultTransFunc(key string) string {
	return key

}

func newStore(opts storeOpts) *store {
	return &store{storeOpts: opts}
}

func (s *store) WriteStream(key string, r io.Reader) error {
	// transform key
	pathKey := s.pathTransformFunc(key)
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	//fileNameBytes := md5.Sum(buff.Bytes())
	//fileName := hex.EncodeToString(fileNameBytes[:])
	name := pathKey.FullPath()
	//create file
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	//copy buffer content
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Printf("Written (%d) bytes to disk in file : %s", n, name)
	return nil
}

// refactor
func (s *store) ReadStream(key string) (io.Reader, error) {
	pathKey := s.pathTransformFunc(key)
	f, err := os.Open(pathKey.FullPath())
	if err != nil {
		return nil, err
	}
	return f, nil
}
