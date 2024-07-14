package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

type pathTransformFunc func(string) string
type storeOpts struct {
	pathTransformFunc pathTransformFunc
}
type store struct {
	storeOpts
}

func CASTransFunc(key string) string {

	hash := sha1.Sum([]byte(key))

	hashString := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashString) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}
	return strings.Join(paths, "/")
}
func defaultTransFunc(key string) string {
	return key

}

func newStore(opts storeOpts) *store {
	return &store{storeOpts: opts}
}

func (s *store) WriteStream(key string, r io.Reader) error {
	// transform key
	pathName := s.pathTransformFunc(key)
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}
	buff := new(bytes.Buffer)
	io.Copy(buff, r)

	fileNameBytes := md5.Sum(buff.Bytes())
	fileName := hex.EncodeToString(fileNameBytes[:])
	name := pathName + "/" + fileName
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, buff)
	if err != nil {
		return err
	}
	log.Printf("Written (%d) bytes to disk in file : %s", n, name)
	return nil
}
