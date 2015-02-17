package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

// File backend uses a plain unbeffered file as its buffer.
type FileBackend struct {
	srcLoc    string
	bufferLoc string
	file      Rwsc //ReaderWriterSeekerCloser
	viewId    int

	bufferSize int64 // Internal buffer size for file ops

	// file state
	ln, col, prevCol int
	offset           int64
	lnCount          int
	length           int64
}

// File backend, the source is a file, the buffer is a copy of the file in the buffer dir.
func (e *Editor) NewFileBackend(loc string, viewId int) (*FileBackend, error) {
	b := &FileBackend{
		viewId:     viewId,
		srcLoc:     loc,
		ln:         1,
		lnCount:    1,
		col:        1,
		prevCol:    1,
		bufferSize: 65536,
	}
	fb := Ed.BufferFile(viewId)
	b.bufferLoc = fb
	if fb != loc {
		// unless we are opening the buffer directly,
		// make sure there is no existing buffer content
		os.Remove(fb)
	}
	if len(loc) > 0 && fb != loc {
		f, err := os.Open(loc)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		b.length = stat.Size()
		if b.length > 10000000 {
			b.bufferLoc = loc
			Ed.SetStatusErr("EDITING IN PLACE ! (Large file)")
		} else {
			err = CopyFile(b.srcLoc, b.bufferLoc)
			if err != nil {
				return nil, err
			}
		}
	}
	var err error
	// TODO: is sync necessary or better to call it selectively ??
	b.file, err = os.OpenFile(b.bufferLoc, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		return nil, err
	}

	// get base line count
	b.lnCount, _ = CountLines(b.file)
	if b.lnCount == 0 {
		b.lnCount = 1
	}

	b.reset()
	return b, nil
}

func (f *FileBackend) SrcLoc() string {
	return f.srcLoc
}

func (f *FileBackend) BufferLoc() string {
	return f.bufferLoc
}

func (f *FileBackend) Close() error {
	return f.file.Close()
}

func (f *FileBackend) Insert(row, col int, text string) error {
	b := []byte(text)
	ln := int64(len(b))
	err := f.seek(row, col)
	if err != nil {
		return err
	}
	// Make a hole in the file by shifting content after insertion point to the right
	err = f.shiftFileBits(ln)
	if err != nil {
		return err
	}
	// insert in the hole created
	_, err = f.file.WriteAt(b, f.offset)
	if err != nil {
		return err
	}
	// Update line Count
	f.lnCount += bytes.Count(b, LineSep)
	return err
}

func (f *FileBackend) Remove(row1, col1, row2, col2 int) error {
	err := f.seek(row2, col2)
	if err != nil {
		return err
	}
	end := f.offset
	if end >= f.length {
		return nil
	}
	err = f.seek(row1, col1)
	if err != nil {
		return err
	}
	ln := end - f.offset + 1
	if ln <= 0 {
		return nil
	}
	buf := make([]byte, ln)
	n := 0
	n, err = f.file.ReadAt(buf, f.offset)
	if err != nil {
		return err
	}
	err = f.shiftFileBits(-int64(n))
	if err != nil {
		return err
	}
	f.lnCount -= bytes.Count(buf[:n], LineSep)
	return nil
}

func (f *FileBackend) LineCount() int {
	return f.lnCount
}

// Slice returns the runes that are in the given rectangle.
// row2 / col2 maybe -1, meaning all lines / whole lines
func (f *FileBackend) Slice(row, col, row2, col2 int) *Slice {
	slice := &Slice{
		text: [][]rune{},
		r1:   row,
		c1:   col,
		r2:   row2,
		c2:   col2,
	}
	if row < 1 || col < 1 {
		return slice
	}
	if row2 != -1 && row > row2 {
		row, row2 = row2, row
	}
	if col2 != -1 && col > col2 {
		col, col2 = col2, col
	}
	r := row
	for ; row2 == -1 || r <= row2; r++ {
		err := f.seek(r, col)
		if err != nil {
			if err != io.EOF {
				panic(err)
			} else {
				return slice
			}
		}
		ln := []rune{}
		for col2 == -1 || f.col <= col2 {
			rune, _, err := f.readRune()
			if err != nil {
				if err == io.EOF && len(ln) > 0 {
					slice.text = append(slice.text, ln)
				}
				return slice
			}
			if rune == '\n' {
				break
			}
			ln = append(ln, rune)
		}
		slice.text = append(slice.text, ln)
	}
	return slice
}

func (f *FileBackend) Save(loc string) error {
	if loc == f.bufferLoc {
		return nil // editing in place
	}
	f.srcLoc = loc
	err := CopyFile(f.bufferLoc, loc)
	// some sort of rsync would be nice ?
	if err != nil {
		return err
	}
	// temporary test hack for go format
	// this should eventually go trough eventing
	if strings.HasSuffix(loc, ".go") {
		err := exec.Command("goimports", "-w", loc).Run()
		// ignore if it fails for now
		if err == nil {
			v := Ed.ViewById(f.viewId)
			x, y := v.CurCol(), v.CurLine()
			Ed.Open(loc, v, "")
			v.MoveCursor(x, y)
		}
	}
	return nil
}

func (f *FileBackend) ViewId() int {
	return f.viewId
}

// seek moves the offest to the given row/col
func (f *FileBackend) seek(row, col int) error {
	if row == f.ln && col == f.col {
		return nil
	}
	// Seek to the beginning of the right row
	if row < f.ln/2 {
		// absolute move more efficient
		if err := f.reset(); err != nil {
			return err
		}
	}
	if err := f.seekRow(row); err != nil {
		return err
	}
	// now seek to the right col
	for f.col < col {
		r, _, err := f.readRune()
		if err != nil {
			return err
		}
		// If we are given a column passed EOL, then stop at EOL.
		if r == '\n' {
			f.col = f.prevCol
			f.ln--
			f.offset--
			_, err := f.file.Seek(-1, 1)
			return err
		}
	}
	return nil
}

// seek to the beginning of a row
func (f *FileBackend) seekRow(row int) error {
	var err error
	if row > f.ln {
		err = f.seekRowFwd(row)
	} else {
		err = f.seekRowRwd(row)
	}
	f.col = 1
	return err
}

func (f *FileBackend) seekRowFwd(row int) error {
	for f.ln != row {
		_, _, err := f.readRune()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileBackend) seekRowRwd(row int) error {
	// Move to the end of the previous line
	for f.ln != row-1 {
		if f.offset == 0 {
			return nil // start of file
		}
		f.offset--
		f.col--
		_, err := f.file.Seek(-1, 1)
		if err != nil {
			return err
		}
		r, _, err := f.peekRune()
		if err != nil {
			return err
		}
		if r == '\n' {
			f.ln--
		}
	}
	// we are just before the '\n' of the previous line, so consume it
	// and we will be wehere we want
	_, _, err := f.readRune()
	return err
}

func (f *FileBackend) readRune() (r rune, size int, err error) {
	b := []byte{0}
	_, err = f.file.Read(b)
	if err != nil {
		return
	}
	size = 1
	r = rune(b[0])
	if r >= 0x80 {
		// This is a UTF rune that uses more than 1 byte (up to 4)
		b = []byte{b[0], 0, 0, 0}
		n := 0
		n, err = f.file.Read(b[1:])
		if err != nil {
			return
		}
		r, size = utf8.DecodeRune(b[:n+1])
		f.file.Seek(int64(size-n-1), 1)
	}
	// adjust offset, ln, col
	f.offset += int64(size)
	if r == '\n' {
		f.prevCol = f.col
		f.col = 1
		f.ln++
	} else {
		f.col++
	}
	return
}

// read a rune without leaving offsets in place
func (f *FileBackend) peekRune() (r rune, size int, err error) {
	b := []byte{4}
	_, err = f.file.ReadAt(b, f.offset)
	if err != nil {
		return
	}
	size = 1
	n := 0
	r = rune(b[0])
	if r >= 0x80 {
		// This is a UTF rune that uses more than 1 byte (up to 4)
		b = []byte{b[0], 0, 0, 0}
		n, err = f.file.ReadAt(b[1:], f.offset+1)
		if err != nil {
			return
		}
		r, size = utf8.DecodeRune(b[:n])
	}
	return
}

// reset puts us back at the beginning of the file
func (f *FileBackend) reset() error {
	f.offset = 0
	f.ln = 1
	f.col = 1
	_, err := f.file.Seek(0, 0)
	return err
}

func (f *FileBackend) size() int64 {
	return f.length
}

// Shift the bits in a file by "shift" value (can be negative).
// Will either expand the file and create a hole, or shrink it.
func (f *FileBackend) shiftFileBits(shift int64) error {
	size := f.size()
	buf := make([]byte, f.bufferSize)
	end := size
	if shift < 0 {
		// Shrinking
		from := f.offset - shift
		for from <= size {
			end = from + f.bufferSize
			if end > size {
				end = size
			}
			n, err := f.file.ReadAt(buf[:end-from], from)
			if err != nil {
				return err
			}
			n, err = f.file.WriteAt(buf[:n], from+shift)
			if err != nil {
				return err
			}
			from += f.bufferSize
		}
		f.length += shift
		os.Truncate(f.bufferLoc, f.length)
		return nil
	}
	os.Truncate(f.bufferLoc, f.length)
	// Expanding
	for end > f.offset {
		from := end - f.bufferSize
		if from < f.offset {
			from = f.offset
		}
		n, err := f.file.ReadAt(buf[:end-from], from)
		if err != nil {
			return err
		}
		n, err = f.file.WriteAt(buf[:n], from+int64(shift))
		if err != nil {
			return err
		}
		end -= f.bufferSize
	}
	f.length += shift
	return nil
}

func (f *FileBackend) Wipe() {
	os.Truncate(f.bufferLoc, 0)
	f.file.Seek(0, 0)
	f.offset = 0
	f.ln = 1
	f.col = 1
	f.prevCol = 1
	f.lnCount = 1
	f.length = 0
}
