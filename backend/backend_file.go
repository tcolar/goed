package backend

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/tcolar/goed/core"
)

// File backend uses a plain unbeffered file as its buffer.
type FileBackend struct {
	srcLoc    string
	bufferLoc string
	file      core.Rwsc //ReaderWriterSeekerCloser
	viewId    int64

	bufferSize int64 // Internal buffer size for file ops

	// file state
	ln, col, prevCol int
	offset           int64
	lnCount          int
	length           int64
}

// File backend, the source is a file, the buffer is a copy of the file in the buffer dir.
func NewFileBackend(loc string, viewId int64) (*FileBackend, error) {
	b := &FileBackend{
		viewId:     viewId,
		srcLoc:     loc,
		ln:         0,
		lnCount:    1,
		col:        0,
		prevCol:    0,
		bufferSize: 65536,
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

func (b *FileBackend) Reload() error {
	// TODO : check dirty
	b.Close()
	fb := BufferFile(b.viewId)
	b.bufferLoc = fb
	if fb != b.srcLoc {
		// unless we are opening the buffer directly,
		// make sure there is no existing buffer content
		os.Remove(fb)
	}
	if len(b.srcLoc) > 0 && fb != b.srcLoc {
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
		if b.length > 10000000 {
			b.bufferLoc = b.srcLoc
			if core.Ed != nil {
				core.Ed.SetStatusErr("EDITING IN PLACE ! (Large file)")
			}
		} else {
			err = core.CopyFile(b.srcLoc, b.bufferLoc)
			if err != nil {
				return err
			}
		}
	}
	var err error
	// TODO: is sync necessary or better to call it selectively ??
	b.file, err = os.OpenFile(b.bufferLoc, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
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
	f.lnCount += bytes.Count(b, core.LineSep)
	return nil
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
	if ln < 0 {
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
	f.lnCount -= bytes.Count(buf[:n], core.LineSep)
	return nil
}

func (f *FileBackend) LineCount() int {
	return f.lnCount
}

// Slice returns the runes that are in the given rectangle.
// row2 / col2 maybe -1, meaning all lines / whole lines
func (f *FileBackend) Slice(row, col, row2, col2 int) *core.Slice {
	slice := core.NewSlice(row, col, row2, col2, [][]rune{})
	text := slice.Text()
	if row < 0 || col < 0 {
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
				core.Ed.SetStatusErr(err.Error())
				return slice
			}
		}
		ln := []rune{}
		for col2 == -1 || f.col <= col2 {
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
	if loc == f.bufferLoc {
		return nil // editing in place
	}
	f.srcLoc = loc
	err := core.CopyFile(f.bufferLoc, loc)
	// some sort of rsync would be nice ?
	if err != nil {
		return err
	}
	// temporary test hack for go format
	// this should eventually go trough eventing
	if strings.HasSuffix(loc, ".go") {
		e := core.Ed
		// TODO: make this configurable, ie: gofmt, goimports etc ....
		// TODO: generalize error panel
		out, _ := exec.Command("goimports", "-w", loc).CombinedOutput()
		fp := path.Join(core.Home, "errors.txt")
		if len(out) > 0 {
			file, _ := os.Create(fp)
			file.Write(out)
			file.Close()
			v := e.ViewByLoc(fp)
			v, err = e.Open(fp, v, "Errors")
			if err != nil {
				return err
			}
			return errors.New("goimports failed")
		}
		if err != nil {
			return errors.New(err.Error())
		}
		v := e.ViewByLoc(fp)
		if v != nil {
			e.DelView(v, true)
		}
		f.Reload()
	}
	return err
}

func (f *FileBackend) ViewId() int64 {
	return f.viewId
}

// seek moves the offest to the given row/col
func (f *FileBackend) seek(row, col int) error {
	// absolute move likely more efficient if we are looking back 500 lines or more
	if row < f.ln && f.ln-row > 500 {
		// Seek to the beginning of the right row
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
	f.col = 0
	return err
}

func (f *FileBackend) seekRowFwd(row int) error {
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
				if f.ln >= row {
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

func (f *FileBackend) seekRowRwd(row int) error {
	// Move to the end of the previous line
	for f.ln >= row {
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
	f.ln = 0
	f.col = 0
	f.prevCol = 0
	f.lnCount = 1
	f.length = 0
}
