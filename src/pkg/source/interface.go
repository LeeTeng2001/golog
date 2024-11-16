package source

import "io"

var AllReaders = map[string]Reader{}

type Reader interface {
	io.ReadSeeker
	OnlyStream() bool
	Init(source string) error
}

// should be called during init phase
func Register(id string, r Reader) {
	if _, ok := AllReaders[id]; ok {
		panic("duplicate reader id: " + id)
	}
	AllReaders[id] = r
}

// TODO: remote ssh / stdin / stream
