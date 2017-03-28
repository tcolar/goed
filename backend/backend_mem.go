package backend

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"unicode/utf8"

	"github.com/tcolar/goed/core"
)

var _ core.Backend = (*MemBackend)(nil)

// MemBackend is a Backend implementation backed by an in memory bufer.
type MemBackend struct {
	text   [][]rune
	colors [][]*color
	file   string
	viewId int64
	lock   sync.Mutex
	vtCols int
}

// NewmemBackend creates a new in memory backend by reading a file.
func NewMemBackend(loc string, viewId int64) (*MemBackend, error) {
	m := &MemBackend{
		text:   [][]rune{[]rune{}},
		colors: [][]*color{[]*color{}},
		file:   loc,
		viewId: viewId,
		vtCols: 80,
	}
	err := m.Reload()
	return m, err
}

func (m *MemBackend) SetVtCols(cols int) {
	m.vtCols = cols
}

func (m *MemBackend) ColorAt(ln, col int) (fg, bg core.Style) {
	t := core.Ed.Theme()
	m.lock.Lock()
	defer m.lock.Unlock()
	if ln >= len(m.colors) {
		return t.Fg, t.Bg
	}
	if col >= len(m.colors[ln]) {
		return t.Fg, t.Bg
	}
	if m.colors[ln][col] == nil {
		return t.Fg, t.Bg
	}
	return m.colors[ln][col].fg, m.colors[ln][col].bg
}

func (m *MemBackend) Reload() error {
	m.Wipe()
	m.lock.Lock()
	defer m.lock.Unlock()
	if len(m.file) == 0 {
		return nil
	}
	if _, err := os.Stat(m.file); os.IsNotExist(err) {
		return nil
	}
	data, err := ioutil.ReadFile(m.file)
	if err != nil {
		return err
	}
	m.text = core.StringToRunes(string(data))
	if len(m.text) == 0 {
		m.text = append(m.text, []rune{})
	}
	if core.Ed != nil {
		v := core.Ed.ViewById(m.viewId)
		if v != nil {
			v.SetDirty(false)
		}
	}
	return nil
}

func (b *MemBackend) Save(loc string) error {
	// TODO: check dirty ?
	b.lock.Lock()
	defer b.lock.Unlock()
	if len(loc) == 0 {
		loc = b.file
	}
	if len(loc) == 0 {
		return fmt.Errorf("Save where ? Use save [path]")
	}
	_, err := os.Stat(loc)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(loc), 0750); err != nil {
			return err
		}
	}
	f, err := os.Create(loc)
	if err != nil {
		return fmt.Errorf("Saving Failed ! %v", loc)
	}
	defer f.Close()
	buf := make([]byte, 4)
	for i, l := range b.text {
		for _, c := range l {
			n := utf8.EncodeRune(buf, c)
			_, err := f.Write(buf[0:n])
			if err != nil {
				return fmt.Errorf("Saved Failed failed %v", err.Error())
			}
		}
		if i != b.LineCount() || len(l) != 0 {
			f.WriteString("\n")
		}
	}
	b.file = loc
	core.Ed.SetStatus("Saved " + b.file)
	return nil
}

func (b *MemBackend) SrcLoc() string {
	return b.file
}

func (b *MemBackend) BufferLoc() string {
	return "_MEM_" // TODO : BufferLoc for in-memory ??
}

func (b *MemBackend) Append(text string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	runes := core.StringToRunes(text)
	row := len(b.text) - 1
	for i, ln := range runes {
		if i == 0 {
			b.text[row] = append(b.text[row], ln...)
		} else {
			b.text = append(b.text, ln)
		}
	}
	return nil
}

func (b *MemBackend) Insert(row, col int, text string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	runes := core.StringToRunes(text)
	if len(runes) == 0 {
		return nil
	}
	var tail []rune
	last := len(runes) - 1
	// Create a "hole" for the new lines to be inserted
	if len(runes) > 1 {
		for i := 1; i < len(runes); i++ {
			b.text = append(b.text, []rune{})
		}
		copy(b.text[row+last:], b.text[row:])
	}
	if row == len(b.text) { // appending new line at end of file
		b.text = append(b.text, []rune{})
	}
	for i, ln := range runes {
		line := b.text[row+i]
		if i == 0 && last == 0 {
			line = append(line, ln...)           // grow line
			copy(line[col+len(ln):], line[col:]) // create hole
			copy(line[col:], ln)                 //fill hole
		} else if i == 0 {
			tail = make([]rune, len(line)-col)
			copy(tail, line[col:])
			line = append(line[:col], ln...)
		} else if i == last {
			line = append(ln, tail...)
		} else {
			line = ln
		}
		b.text[row+i] = line
	}

	return nil
}

func (b *MemBackend) Remove(row1, col1, row2, col2 int) error {
	h := core.SyntaxHighlighting()
	b.lock.Lock()
	defer b.lock.Unlock()
	if row1 < 0 {
		row1 = 0
	}
	if row1 >= len(b.text) {
		row1 = len(b.text) - 1
	}

	if col1 < 0 {
		col1 = 0
	}
	if col2 < -1 { // -1 maybe use to remove up to before the line (after prev line EOL)
		col2 = -1
	}
	if row2 < len(b.text) && col2 >= len(b.text[row2]) {
		// that is, if at end of line, then start from beginning of next row
		row2++
		col2 = -1
	}
	if row2 < 0 {
		row2 = 0
	}
	if row2 >= len(b.text) {
		row2 = len(b.text) - 1
	}
	if col1 > len(b.text[row1]) {
		col1 = len(b.text[row1])
	}
	if col2 >= len(b.text[row2]) {
		col2 = len(b.text[row2]) - 1
	}
	b.text[row1] = append(b.text[row1][:col1], b.text[row2][col2+1:]...)
	if h {
		b.colors[row1] = append(b.colors[row1][:col1], b.colors[row2][col2+1:]...)
	}
	drop := row2 - row1
	if drop > 0 {
		copy(b.text[row1+1:], b.text[row1+1+drop:])
		b.text = b.text[:len(b.text)-drop]
		if h {
			copy(b.colors[row1+1:], b.colors[row1+1+drop:])
			b.colors = b.colors[:len(b.colors)-drop]
		}
	}

	return nil
}

func (b *MemBackend) sliceNoLock(row, col, row2, col2 int) *core.Slice {
	slice := core.NewSlice(row, col, row2, col2, [][]rune{})
	text := slice.Text()
	if row < 0 || col < 0 {
		return slice
	}
	r := slice.R1
	for ; slice.R2 == -1 || r <= slice.R2; r++ {
		if r < 0 || r >= len(b.text) {
			break
		}
		if slice.C2 == -1 {
			*text = append(*text, b.text[r])
		} else {
			c, c2, l := slice.C1, slice.C2+1, len(b.text[r])
			if c > l {
				c = l
			}
			if c2 > l {
				c2 = l
			}
			*text = append(*text, b.text[r][c:c2])
		}
	}
	return slice
}

func (b *MemBackend) Slice(row, col, row2, col2 int) *core.Slice {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.sliceNoLock(row, col, row2, col2)
}

func (b *MemBackend) LineCount() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	count := len(b.text)
	if count > 0 {
		last := len(b.text) - 1
		if last > 0 && len(b.text[last]) == 0 {
			count--
		}
	}
	return count
}

func (b *MemBackend) Close() error {
	b.Wipe()
	return nil
}

func (b *MemBackend) ViewId() int64 {
	return b.viewId
}

func (b *MemBackend) Wipe() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.text = [][]rune{[]rune{}}
	b.colors = [][]*color{[]*color{}}
}

func (m *MemBackend) SendBytes(data []byte) {}

func (m *MemBackend) OnActivate() {}

func (m *MemBackend) BufferOffset() int64 { return 0 /*TODO*/ }

type color struct {
	fg, bg core.Style
}
