package ui

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/core"
)

func init() {
	core.Testing = true
	core.InitHome()
	core.Ed = newMockEditor()
	core.Ed.Start("")
}

func TestCountLines(t *testing.T) {
	// CountLines
	f, err := os.Open("../test_data/file1.txt")
	assert.Nil(t, err, "file")
	defer f.Close()
	lns, err := core.CountLines(f)
	assert.Nil(t, err, "CountLines")
	assert.Equal(t, lns, 12, "CountLines")
}

func TestStringToRunes(t *testing.T) {
	r := [][]rune{}
	s := core.RunesToString(r)
	assert.Equal(t, s, "", "runestostring1")
	assert.Equal(t, core.StringToRunes(s), r, "stringToRunes1")
	r = [][]rune{
		[]rune{'A', 'B', 'C'},
	}
	s = core.RunesToString(r)
	assert.Equal(t, s, "ABC", "runestostring2")
	assert.Equal(t, core.StringToRunes(s), r, "stringToRunes2")
	r = append(r, []rune{}, []rune{'1', '2'})
	s = core.RunesToString(r)
	assert.Equal(t, s, "ABC\n\n12", "runestostring3")
	assert.Equal(t, core.StringToRunes(s), r, "stringToRunes3")
	r = [][]rune{
		[]rune{},
		[]rune{'2'},
	}
	s = core.RunesToString(r)
	assert.Equal(t, s, "\n2", "runestostring4")
	assert.Equal(t, core.StringToRunes(s), r, "stringToRunes4")
	r = [][]rune{
		[]rune{'1'},
		[]rune{},
	}
	s = core.RunesToString(r)
	assert.Equal(t, s, "1\n", "runestostring5")
	assert.Equal(t, core.StringToRunes(s), r, "stringToRunes5")
}

func TestQuitCheck(t *testing.T) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView()
	v2 := Ed.NewView()
	col := Ed.NewCol(1.0, []*View{v, v2})
	Ed.Cols = []*Col{col}
	then := time.Now()
	assert.True(t, Ed.QuitCheck(), "quitcheck1")
	v2.Dirty = true
	assert.False(t, Ed.QuitCheck(), "quitcheck2")
	assert.True(t, v2.lastCloseTs.After(then), "quitcheck ts")
	assert.True(t, Ed.QuitCheck(), "quitcheck3")
}

func TestTheme(t *testing.T) {
	th, err := core.ReadTheme("../test_data/theme.toml")
	assert.Nil(t, err, "theme")
	s := core.NewStyle(0)
	s.UnmarshalText([]byte("99663311"))
	assert.Equal(t, th.Bg, s, "theme bg")
	sb := th.Statusbar
	s.UnmarshalText([]byte("EB070000"))
	s2 := core.NewStyle(0)
	s.UnmarshalText([]byte("EB000000"))
	sr := core.StyledRune{
		Rune: '❊',
		Bg:   s,
		Fg:   s2,
	}
	assert.Equal(t, sb, sr, "styled rune")
	s = core.NewStyle(0x41)
	s = s.WithAttr(Bold)
	assert.Equal(t, s, core.NewStyle(0x0141), "style attr2")
}
