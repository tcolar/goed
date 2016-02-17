package core

import "io"

// Backend represent the backend(data operations) of a View.
// Backend implements the low level data handling.
type Backend interface {
	// SrcLoc returns the location of the original data.
	SrcLoc() string
	// BufferLoc is the location of the "copy" the backend works directly on.
	BufferLoc() string

	Insert(line, col int, text string) error
	Append(text string) error
	Remove(line1, col1, line2, col2 int) error

	LineCount() int

	// Save saves the edited data (BufferLoc) into the original (SrcLoc)
	Save(loc string) error

	// Slice gets a region of text ("rectangle") as a runes matrix
	Slice(line1, col, line2, col2 int) *Slice

	// Close closes the backend resources.
	Close() error

	// ViewId returns the "unique" viewid given to this buffer.
	ViewId() int64

	// Completely clears the buffer text (empty document)
	Wipe()

	// Reloads the text (from SrcLoc to BufferLoc)
	Reload() error

	// return the color style at a specific location (mem backends)
	ColorAt(ln, col int) (fg, bg Style)

	//Sync() error         // sync from source ?
	//IsStale() bool       // whether the source as changed under us (fsnotify)
	//IsBufferStale() bool // whether the buffer has changed under us

	//SourceMd5 or ts?
	//BufferMd5 or ts?
	SetVtCols(cols int)
}

type Rwsc interface {
	io.Reader
	io.Writer
	io.ReaderAt
	io.WriterAt
	io.Seeker
	io.Closer
}
