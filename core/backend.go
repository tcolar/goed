package core

import "io"

type Backend interface {
	SrcLoc() string    // "original" source
	BufferLoc() string // buffer location

	Insert(line, col int, text string) error
	Append(text string) error
	Remove(line1, col1, line2, col2 int) error

	LineCount() int

	Save(loc string) error

	// Get a region ("rectangle") as a runes matrix
	Slice(line1, col, line2, col2 int) *Slice

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
