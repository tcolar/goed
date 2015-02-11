package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWm(t *testing.T) {
	// TODO: test window manager
}

func TestCountLines(t *testing.T) {
	// CountLines
	f, err := os.Open("test_data/file1.txt")
	assert.Nil(t, err, "file")
	defer f.Close()
	lns, err := CountLines(f)
	assert.Nil(t, err, "CountLines")
	assert.Equal(t, lns, 12, "CountLines")
}

func TestStringToRunes(t *testing.T) {
	r := [][]rune{}
	s := Ed.RunesToString(r)
	assert.Equal(t, s, "", "runestostring1")
	assert.Equal(t, Ed.StringToRunes(s), r, "stringToRunes1")
	r = [][]rune{
		[]rune{'A', 'B', 'C'},
	}
	s = Ed.RunesToString(r)
	assert.Equal(t, s, "ABC", "runestostring2")
	assert.Equal(t, Ed.StringToRunes(s), r, "stringToRunes2")
	r = append(r, []rune{}, []rune{'1', '2'})
	s = Ed.RunesToString(r)
	assert.Equal(t, s, "ABC\n\n12", "runestostring3")
	assert.Equal(t, Ed.StringToRunes(s), r, "stringToRunes3")
	r = [][]rune{
		[]rune{},
		[]rune{'2'},
	}
	s = Ed.RunesToString(r)
	assert.Equal(t, s, "\n2", "runestostring4")
	assert.Equal(t, Ed.StringToRunes(s), r, "stringToRunes4")
	r = [][]rune{
		[]rune{'1'},
		[]rune{},
	}
	s = Ed.RunesToString(r)
	assert.Equal(t, s, "1\n", "runestostring5")
	assert.Equal(t, Ed.StringToRunes(s), r, "stringToRunes5")
}

func TestQuitCheck(t *testing.T) {
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
