package main

import "io"

// TODO: flush this out + File based impl
// TODO: use bufio
type Backend interface {
	SrcLoc() string    // "original" source
	BufferLoc() string // buffer location

	Insert(row, col int, text string) error
	Remove(row, col int, text string) error

	LineCount() int

	Save(loc string) error

	// Get a region ("rectangle") as a runes matrix
	Slice(row, col, line2, row2 int) [][]rune

	Close() error

	//Sync() error         // sync from source ?
	//IsStale() bool       // whether the source as changed under us (fsnotify)
	//IsBufferStale() bool // whether the buffer has changed under us

	//SourceMd5 or ts?
	//BufferMd5 or ts ?
}

type Rwsc interface {
	io.Reader
	io.Writer
	io.ReaderAt
	io.WriterAt
	io.Seeker
	io.Closer
}
