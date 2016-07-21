package core

import (
	"os"
	"testing"
	"time"

	"github.com/tcolar/goed/assert"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CoreSuite struct {
}

var _ = Suite(&CoreSuite{})

func (cs *CoreSuite) SetUpSuite(c *C) {
	Testing = true
	InitHome(time.Now().Unix())
}

func (cs *CoreSuite) TestCountLines(t *C) {
	// CountLines
	f, err := os.Open("../test_data/file1.txt")
	assert.Nil(t, err)
	defer f.Close()
	lns, err := CountLines(f)
	assert.Nil(t, err)
	assert.Eq(t, lns, 12)
}

func (cs *CoreSuite) TestStringToRunes(t *C) {
	r := [][]rune{}
	s := RunesToString(r)
	assert.Eq(t, s, "")
	assert.DeepEq(t, StringToRunes(s), r)
	r = [][]rune{
		[]rune{'A', 'B', 'C'},
	}
	s = RunesToString(r)
	assert.Eq(t, s, "ABC")
	assert.DeepEq(t, StringToRunes(s), r)
	r = append(r, []rune{}, []rune{'1', '2'})
	s = RunesToString(r)
	assert.Eq(t, s, "ABC\n\n12")
	assert.DeepEq(t, StringToRunes(s), r)
	r = [][]rune{
		[]rune{},
		[]rune{'2'},
	}
	s = RunesToString(r)
	assert.Eq(t, s, "\n2")
	assert.DeepEq(t, StringToRunes(s), r)
	r = [][]rune{
		[]rune{'1'},
		[]rune{},
	}
	s = RunesToString(r)
	assert.Eq(t, s, "1\n")
	assert.DeepEq(t, StringToRunes(s), r)
}

func (cs *CoreSuite) TestTheme(t *C) {
	th, err := ReadTheme("../test_data/theme.toml")
	assert.Nil(t, err)
	s := NewStyle(0)
	s.UnmarshalText([]byte("99663311"))
	assert.Eq(t, th.Bg, s)
	sb := th.Statusbar
	s.UnmarshalText([]byte("EB070000"))
	s2 := NewStyle(0)
	s.UnmarshalText([]byte("EB000000"))
	sr := StyledRune{
		Rune: '❊',
		Bg:   s,
		Fg:   s2,
	}
	assert.Eq(t, sb, sr)
	s = NewStyle(0x41)
	s = s.WithAttr(Bold)
	assert.Eq(t, s, NewStyle(0x0141))
}

func (cs *CoreSuite) TestIsText(t *C) {
	assert.True(t, IsTextFile("../test_data/empty.txt"))
	assert.True(t, IsTextFile("../test_data/test.txt"))
	assert.False(t, IsTextFile("../test_data/test.bin"))
}
