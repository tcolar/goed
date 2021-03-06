package backend

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/encoding/unicode"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

var _ core.Backend = (*FileBackend)(nil)

// FileBackend is a backend implemetation that uses a plain unbuffered file as
// its buffer.
type FileBackend struct {
	srcLoc    string
	bufferLoc string
	file      core.Rwsc //ReaderWriterSeekerCloser
	viewId    int64
	textInfo  *core.TextInfo

	bufferSize int64 // Internal buffer size for file ops

	// file state
	ln, col, prevCol int
	offset           int64
	lnCount          int
	length           int64
	lock             sync.Mutex
}

// NewFileBackend creates a backend from a copy of the file in the buffer dir.
// Note that very large files(>100Mb) might be edited in place.
func NewFileBackend(loc string, viewId int64) (*FileBackend, error) {
	b := &FileBackend{
		viewId:     viewId,
		srcLoc:     loc,
		ln:         0,
		lnCount:    1,
		col:        0,
		prevCol:    0,
		bufferSize: 65536,
		textInfo:   core.CrLfTextInfo(unicode.UTF8, false),
	}
	err := b.Reload()
	return b, err
}

func (f *FileBackend) SrcLoc() string {
	return f.srcLoc
}

func (f *FileBackend) BufferLoc() string {
	return f.bufferLoc
}

func (f *FileBackend) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func (f *FileBackend) ColorAt(ln, col int) (fg, bg core.Style) {
	return core.Ed.Theme().Fg, core.Ed.Theme().Bg
}

func (b *FileBackend) Reload() error {
	b.lock.Lock()
	defer b.lock.Unlock()
	// TODO : check dirty
	b.Close()
	fb := BufferFile(b.viewId)
	b.bufferLoc = fb
	if fb != b.srcLoc {
		// unless we are opening the buffer directly,
		// make sure there is no existing buffer content
		os.Remove(fb)
	}
	newFile := false
	if _, err := os.Stat(b.srcLoc); os.IsNotExist(err) {
		newFile = true
	}
	if !newFile && len(b.srcLoc) > 0 && fb != b.srcLoc {
		usesCrLf := core.UsesCrLf(b.srcLoc)
		f, err := os.Open(b.srcLoc)
		if err != nil {
			return err
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil {
			return err
		}
		b.length = stat.Size()
		b.textInfo = core.ReadTextInfo(b.srcLoc, usesCrLf)
		if b.textInfo == nil {
			return fmt.Errorf("Unsupported encoding ? Binary file ? %s", b.srcLoc)
		}
		if b.length > 10000000 {
			b.bufferLoc = b.srcLoc
			if core.Ed != nil {
				core.Ed.SetStatusErr("EDITING IN PLACE ! (Large file)")
			}
		} else {
			err = core.CopyToUTF8(b.srcLoc, b.bufferLoc, b.textInfo.Enc)
			if err != nil {
				return err
			}
		}
	}
	var err error
	// TODO: is sync necessary or better to call it selectively ??
	b.file, err = os.OpenFile(b.bufferLoc, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0640)
	if err != nil {
		return err
	}

	// Ensure file end with '\n', POSIX rule and simplifies usage.
	_, err = b.file.Seek(-1, 2)
	r := make([]byte, 1)
	if err == nil {
		_, err = b.file.Read(r)
	}
	// Note: err likely would mean empty file
	if err != nil || r[0] != '\n' {
		b.file.Write([]byte{'\n'})
	}
	b.file.Seek(0, 0)

	// get base line count
	b.lnCount, _ = core.CountLines(b.file)
	if b.lnCount == 0 {
		b.lnCount = 1
	}

	b.reset()
	if core.Ed != nil {
		v := core.Ed.ViewById(b.viewId)
		if v != nil {
			v.SetDirty(false)
		}
	}
	return nil
}

func (f *FileBackend) Append(text string) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	b := []byte(text)
	_, err := f.file.WriteAt(b, f.length)
	if err != nil {
		return err
	}
	// Update line Count
	f.lnCount += bytes.Count(b, core.LineSep)
	return nil
}

func (f *FileBackend) Insert(row, col int, text string) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	b := []byte(text)
	ln := int64(len(b))
	err := f.seek(row, col)
	if err != nil {
		return err
	}
	// Make a hole in the file by shifting content after insertion point to the right
	err = f.shiftFileBytes(ln)
	if err != nil {
		return err
	}
	// insert in the hole created
	_, err = f.file.WriteAt(b, f.offset)
	if err != nil {
		return err
	}
	// Update line Count
	f.lnCount += bytes.Count(b, core.LineSep)
	return nil
}

func (f *FileBackend) Remove(row1, col1, row2, col2 int) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	err := f.seek(row2, col2) // beginning of last char to be removed
	if err != nil {
		return err
	}
	f.readRune() // move to end of last char (rune could be more than 1 byte - UTF8)
	end := f.offset
	if end > f.length {
		end = f.length
	}
	err = f.seek(row1, col1) // start of first character
	if err != nil {
		return err
	}
	ln := end - f.offset // read bytes to be removed (to count newlines, could optimize)
	if ln <= 0 {
		return nil
	}
	buf := make([]byte, ln)
	n := 0
	n, err = f.file.ReadAt(buf, f.offset)
	if err != nil {
		return err
	}
	err = f.shiftFileBytes(-int64(n)) // delete the bytes to be removed
	if err != nil {
		return err
	}
	f.lnCount -= bytes.Count(buf[:n], core.LineSep)
	return nil
}

func (f *FileBackend) LineCount() int {
	return f.lnCount
}

// Slice returns the runes that are in the given rectangle.
// line2 / col2 maybe -1, meaning all lines / whole lines
func (f *FileBackend) Slice(line1, col, line2, col2 int) *core.Slice {
	f.lock.Lock()
	defer f.lock.Unlock()
	slice := core.NewSlice(line1, col, line2, col2, [][]rune{})
	text := slice.Text()
	if line1 < 0 || col < 0 {
		return slice
	}
	l := slice.R1
	for ; slice.R2 == -1 || l <= slice.R2; l++ {
		err := f.seek(l, slice.C1)
		if err != nil {
			if err != io.EOF {
				core.Ed.SetStatusErr(err.Error())
				return slice
			}
		}
		ln := []rune{}
		for slice.C2 == -1 || f.col <= slice.C2 {
			rune, _, err := f.readRune()
			if err != nil {
				if err == io.EOF && len(ln) > 0 {
					*text = append(*text, ln)
				}
				return slice
			}
			if rune == '\n' {
				break
			}
			ln = append(ln, rune)
		}
		*text = append(*text, ln)
	}
	return slice
}

func (f *FileBackend) Save(loc string) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	if loc == f.bufferLoc {
		return nil // editing in place
	}
	_, err := os.Stat(loc)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(loc), 0750); err != nil {
			return err
		}
	}
	f.srcLoc = loc

	err = core.CopyFromUTF8(f.bufferLoc, loc, f.textInfo.Enc)
	// some sort of rsync would be nice ?
	if err != nil {
		return err
	}
	// temporary hack for go format
	// should hook-up through action/eventing
	if strings.HasSuffix(loc, ".go") {
		go actions.ExecScript("goimports.sh")
	}

	return nil
}

func (f *FileBackend) ViewId() int64 {
	return f.viewId
}

func (f *FileBackend) Wipe() {
	f.lock.Lock()
	defer f.lock.Unlock()
	os.Truncate(f.bufferLoc, 0)
	f.file.Seek(0, 0)
	f.offset = 0
	f.ln = 0
	f.col = 0
	f.prevCol = 0
	f.lnCount = 1
	f.length = 0
}

// seek moves the offest to the given line/col
func (f *FileBackend) seek(line, col int) error {
	// absolute move likely more efficient if we are looking back 500 lines or more
	if line < f.ln && f.ln-line > 500 {
		// Seek to the beginning of the right line
		if err := f.reset(); err != nil {
			return err
		}
	}
	if err := f.seekLine(line); err != nil {
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

// seek to the beginning of a line
func (f *FileBackend) seekLine(line int) error {
	var err error
	if line > f.ln {
		err = f.seekLineFwd(line)
	} else {
		err = f.seekLineRwd(line)
	}
	f.col = 0
	return err
}

func (f *FileBackend) seekLineFwd(line int) error {
	buf := make([]byte, 8192)
outer:
	for {
		c, err := f.file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		for _, b := range buf[:c] {
			f.offset++
			f.col++
			if b == '\n' { // TODO: is that good enough with UTF emcoding ??
				f.ln++
				f.col = 0
				if f.ln >= line {
					break outer
				}
			}
		}
		if c < 8192 {
			break // no more to read
		}
	}
	f.file.Seek(f.offset, 0)
	return nil
}

func (f *FileBackend) seekLineRwd(line int) error {
	// Move to the end of the previous line
	for f.ln >= line {
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
		f.col = 0
		f.ln++
	} else {
		f.col++
	}
	return
}

// read a rune whithout moving offsets
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
	f.ln = 0
	f.col = 0
	_, err := f.file.Seek(0, 0)
	return err
}

func (f *FileBackend) size() int64 {
	return f.length
}

// Shift the bytes in a file by "shift" value (can be negative).
// Will either expand the file and create a hole, or shrink it.
func (f *FileBackend) shiftFileBytes(shift int64) error {
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

func (b *FileBackend) SetVtCols(cols int) { // N/A
}

func (b *FileBackend) SendBytes(data []byte) {}

func (b *FileBackend) OnActivate() {}

func (b *FileBackend) OffsetAt(ln, col int) int64 {
	b.lock.Lock()
	defer b.lock.Unlock()
	prevLn, prevCol := b.ln, b.col
	b.seek(ln, col)
	offset := b.offset
	b.seek(prevLn, prevCol)
	return offset
}
