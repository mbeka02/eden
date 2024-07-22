package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

var key = "mySuperDuperPrivateKey"

func TestTransformFunc(t *testing.T) {
	pathKey := CASTransFunc(key)
	assert.NotEmpty(t, pathKey)
	expectedFilename := "0b0fb7591de06089559f2dcaac705176c90ecf00"
	expectedPathName := "0b0fb/7591d/e0608/9559f/2dcaa/c7051/76c90/ecf00"
	assert.Equal(t, expectedPathName, pathKey.PathName)
	assert.Equal(t, expectedFilename, pathKey.Filename)
}
func TestStore(t *testing.T) {
	opts := storeOpts{
		pathTransformFunc: CASTransFunc,
	}
	store := newStore(opts)
	assert.NotNil(t, store)

	//Writing
	data := []byte("random bytes")
	err := store.WriteStream(key, bytes.NewReader(data))
	assert.NoError(t, err)

	//Reading
	r, err := store.Read(key)

	assert.NoError(t, err)
	b, _ := io.ReadAll(r)
	assert.Equal(t, data, b)

	//Deleting
	err = store.Delete(key)
	assert.NoError(t, err)
}
