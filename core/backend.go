package core

import "io"

type Backend interface {
	SrcLoc() string    // "original" source
	BufferLoc() string // buffer location

	Insert(row, col int, text string) error
	Append(text string) error
	Remove(row1, col1, row, col2 int) error

	LineCount() int

	Save(loc string) error

	// Get a region ("rectangle") as a runes matrix
	Slice(row, col, row2, col2 int) *Slice

	Close() error

	ViewId() int64

	// Completely clears the buffer (empty)
	Wipe()

	// Reloads from source
	Reload() error

	//Sync() error         // sync from source ?
	//IsStale() bool       // whether the source as changed under us (fsnotify)
	//IsBufferStale() bool // whether the buffer has changed under us

	//SourceMd5 or ts?
	//BufferMd5 or ts?
}

type Rwsc interface {
	io.Reader
	io.Writer
	io.ReaderAt
	io.WriterAt
	io.Seeker
	io.Closer
}
