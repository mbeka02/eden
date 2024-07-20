package main

import (
	"bytes"
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
	Original string
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
		Original: hashString,
	}
}

func (pk PathKey) Filename() string {
	return fmt.Sprintf("%s/%s", pk.PathName, pk.Original)

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
	buff := new(bytes.Buffer)
	io.Copy(buff, r)

	//fileNameBytes := md5.Sum(buff.Bytes())
	//fileName := hex.EncodeToString(fileNameBytes[:])
	name := pathKey.Filename()
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
