package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const defaultRootFolderName = "root"

type (
	pathTransformFunc func(string) PathKey
	storeOpts         struct {
		// root is the folder name of the root contaning all the folders and files
		root              string
		pathTransformFunc pathTransformFunc
	}
)

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

func defaultTransFunc(key string) PathKey {
	return PathKey{
		PathName: key,
		Filename: key,
	}
}

func newStore(opts storeOpts) *store {
	if opts.root == "" {
		opts.root = defaultRootFolderName
	}
	return &store{storeOpts: opts}
}

func (s *store) writeStream(key string, r io.Reader) (int64, error) {
	// transform key
	pathKey := s.pathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.root, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	fullPath := pathKey.FullPath()
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.root, fullPath)
	// create file
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return 0, err
	}
	// copy buffer content
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.pathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.root, pathKey.FullPath())
	fileInfo, err := os.Stat(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}
	fmt.Printf("file size -> %v", fileInfo.Size())
	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}
	return fileInfo.Size(), file, nil
}

// Read() is a wrapper fn for readStream

// FIXME: Refactor this
func (s *store) Read(key string) (int64, io.Reader, error) {
	return s.readStream(key)
	// if err != nil {
	// 	return 0, nil, err
	// }
	// defer file.Close()
	// buff := new(bytes.Buffer)
	// n, err := io.Copy(buff, f)
	// return file, buff, err
}

// Write() is a wrapper fn for writeStream
func (s *store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *store) Delete(key string) error {
	pathKey := s.pathTransformFunc(key)
	pathNames := strings.Split(pathKey.PathName, "/")
	firstPathName := pathNames[0]
	firstPathNameWithRoot := fmt.Sprintf("%s/%s", s.root, firstPathName)
	defer func() {
		fmt.Printf("deleted %s from disk\n", pathKey.Filename)
	}()
	return os.RemoveAll(firstPathNameWithRoot)
}

func (s *store) Has(key string) bool {
	pathKey := s.pathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

// clear everything including the root folder
func (s *store) Clear() error {
	return os.RemoveAll(s.root)
}
