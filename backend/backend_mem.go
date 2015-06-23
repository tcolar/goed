package backend

import (
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/tcolar/goed/core"
)

type MemBackend struct {
	text   [][]rune
	file   string
	viewId int64
}

func NewMemBackend(path string, viewId int64) (*MemBackend, error) {
	m := &MemBackend{
		text:   [][]rune{[]rune{}},
		file:   path,
		viewId: viewId,
	}
	err := m.Reload()
	return m, err
}

func (m *MemBackend) Reload() error {
	// TODO: check dirty ?
	if len(m.file) == 0 {
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
		v.SetDirty(false)
	}
	return nil
}

func (b *MemBackend) Save(loc string) error {
	if len(loc) == 0 {
		return fmt.Errorf("Save where ? Use save [path]")
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
	row-- // 1 index to 0 index
	col--
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
	if row == len(b.text) { // appending ne wline at end of file
		b.text = append(b.text, []rune{})
	}
	for i, ln := range runes {
		line := b.text[row+i]
		if i == 0 && last == 0 {
			line = append(line, ln...)           // grow line
			copy(line[col+len(ln):], line[col:]) // create hole
			copy(line[col:], ln)                 //file hole
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

	// convert to 0-index for simplicity
	row1--
	row2--
	col1--
	col2--

	if row1 < 0 {
		row1 = 0
	}
	if row1 >= len(b.text) {
		row1 = len(b.text)
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
	if col1 < 0 {
		col1 = 0
	}
	if col2 < 0 {
		col2 = 0
	}
	if col1 > len(b.text[row1]) {
		col1 = len(b.text[row1])
	}
	if col2 >= len(b.text[row2]) {
		col2 = len(b.text[row2]) - 1
	}
	b.text[row1] = append(b.text[row1][:col1], b.text[row2][col2+1:]...)
	drop := row2 - row1
	if drop > 0 {
		copy(b.text[row1+1:], b.text[row1+1+drop:])
		b.text = b.text[:len(b.text)-drop]
	}

	return nil
}

func (b *MemBackend) Slice(row, col, row2, col2 int) *core.Slice {
	slice := core.NewSlice(row, col, row2, col2, [][]rune{})
	text := slice.Text()
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
		if r > len(b.text) {
			break
		}
		if col2 == -1 {
			*text = append(*text, b.text[r-1])
		} else {
			c, c2, l := col-1, col2, len(b.text[r-1])
			if c > l {
				c = l
			}
			if c2 > l {
				c2 = l
			}
			*text = append(*text, b.text[r-1][c:c2])
		}
	}
	return slice
}

func (b *MemBackend) LineCount() int {
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
	return nil // Noop
}

func (b *MemBackend) ViewId() int64 {
	return b.viewId
}

func (b *MemBackend) Wipe() {
	b.text = [][]rune{[]rune{}}
}
