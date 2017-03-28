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

func (cs *CoreSuite) TestIsText(t *C) {
	assert.NotNil(t, ReadTextInfo("../test_data/empty.txt", false))
	assert.NotNil(t, ReadTextInfo("../test_data/test.txt", false))
}
