package core

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	Testing = true
	InitHome(time.Now().Unix())
}

func TestCountLines(t *testing.T) {
	// CountLines
	f, err := os.Open("../test_data/file1.txt")
	assert.Nil(t, err, "file")
	defer f.Close()
	lns, err := CountLines(f)
	assert.Nil(t, err, "CountLines")
	assert.Equal(t, lns, 12, "CountLines")
}

func TestStringToRunes(t *testing.T) {
	r := [][]rune{}
	s := RunesToString(r)
	assert.Equal(t, s, "", "runestostring1")
	assert.Equal(t, StringToRunes(s), r, "stringToRunes1")
	r = [][]rune{
		[]rune{'A', 'B', 'C'},
	}
	s = RunesToString(r)
	assert.Equal(t, s, "ABC", "runestostring2")
	assert.Equal(t, StringToRunes(s), r, "stringToRunes2")
	r = append(r, []rune{}, []rune{'1', '2'})
	s = RunesToString(r)
	assert.Equal(t, s, "ABC\n\n12", "runestostring3")
	assert.Equal(t, StringToRunes(s), r, "stringToRunes3")
	r = [][]rune{
		[]rune{},
		[]rune{'2'},
	}
	s = RunesToString(r)
	assert.Equal(t, s, "\n2", "runestostring4")
	assert.Equal(t, StringToRunes(s), r, "stringToRunes4")
	r = [][]rune{
		[]rune{'1'},
		[]rune{},
	}
	s = RunesToString(r)
	assert.Equal(t, s, "1\n", "runestostring5")
	assert.Equal(t, StringToRunes(s), r, "stringToRunes5")
}

func TestTheme(t *testing.T) {
	th, err := ReadTheme("../test_data/theme.toml")
	assert.Nil(t, err, "theme")
	s := NewStyle(0)
	s.UnmarshalText([]byte("99663311"))
	assert.Equal(t, th.Bg, s, "theme bg")
	sb := th.Statusbar
	s.UnmarshalText([]byte("EB070000"))
	s2 := NewStyle(0)
	s.UnmarshalText([]byte("EB000000"))
	sr := StyledRune{
		Rune: '❊',
		Bg:   s,
		Fg:   s2,
	}
	assert.Equal(t, sb, sr, "styled rune")
	s = NewStyle(0x41)
	s = s.WithAttr(Bold)
	assert.Equal(t, s, NewStyle(0x0141), "style attr2")
}
