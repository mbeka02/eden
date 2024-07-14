package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformFunc(t *testing.T) {
	key := "mySuperDuperPrivateKey"
	pathname := CASTransFunc(key)
	assert.NotEmpty(t, pathname)
	expectedPathname := "0b0fb/7591d/e0608/9559f/2dcaa/c7051/76c90/ecf00"
	assert.Equal(t, expectedPathname, pathname)
}
func TestStore(t *testing.T) {
	opts := storeOpts{
		pathTransformFunc: CASTransFunc,
	}
	store := newStore(opts)
	assert.NotNil(t, store)
	data := bytes.NewReader([]byte("random bytes"))
	err := store.WriteStream("pics", data)
	assert.NoError(t, err)
}
