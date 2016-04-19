package syntax

import (
	"testing"

	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SyntaxSuite struct {
}

var _ = Suite(&SyntaxSuite{})

var testKw = "var gop"

var testSymb = "a := (5*3)>>2"

var testEsc = `foo := "bar\"bar"`

var testSrc = `package dummy
// comment 1
type foo string

/*111
222
333*/
func test(a int) {
	go doStuff([]string{"bar"})
	gopher := 3
}`

var testTop = `aaa
bbbb
x*/a:=3`

func (ss *SyntaxSuite) TestSyntax(t *C) {
	hs := Highlights{}
	hs.Update(core.StringToRunes(testKw), ".go")
	lns := hs.Lines
	assert.Eq(t, len(lns), 1)
	assert.Eq(t, len(lns[0]), 1)
	ss.checkHl(t, lns[0][0], StyleKw1, 0, 2)
	hs.Update(core.StringToRunes(testSymb), ".go")
	lns = hs.Lines
	assert.Eq(t, len(lns), 1)
	assert.Eq(t, len(lns[0]), 5)
	ss.checkHl(t, lns[0][0], StyleSymb1, 2, 3)
	ss.checkHl(t, lns[0][1], StyleSep1, 5, 5)
	ss.checkHl(t, lns[0][2], StyleSymb3, 7, 7)
	ss.checkHl(t, lns[0][3], StyleSep1, 9, 9)
	ss.checkHl(t, lns[0][4], StyleSymb3, 10, 11)
	hs.Update(core.StringToRunes(testEsc), ".go")
	lns = hs.Lines
	assert.Eq(t, len(lns), 1)
	assert.Eq(t, len(lns[0]), 2)
	ss.checkHl(t, lns[0][0], StyleSymb1, 4, 5)
	ss.checkHl(t, lns[0][1], StyleString, 7, 16)
	hs.Update(core.StringToRunes(testSrc), ".go")
	lns = hs.Lines
	assert.Eq(t, len(lns), 11)
	ss.checkHl(t, lns[0][0], StyleKw1, 0, 6)
	ss.checkHl(t, lns[1][0], StyleComment, 0, 12)
	ss.checkHl(t, lns[2][0], StyleKw1, 0, 3)
	ss.checkHl(t, lns[4][0], StyleComment, 0, 4)
	ss.checkHl(t, lns[5][0], StyleComment, 0, 2)
	ss.checkHl(t, lns[6][0], StyleComment, 0, 4)
	assert.Eq(t, len(lns[7]), 4)
	assert.Eq(t, len(lns[8]), 8)
	ss.checkHl(t, lns[9][0], StyleSymb1, 8, 9)
	ss.checkHl(t, lns[10][0], StyleSep1, 0, 0)
	hs.Update(core.StringToRunes(testTop), ".go")
	lns = hs.Lines
	assert.Eq(t, len(lns), 3)
	ss.checkHl(t, lns[0][0], StyleComment, 0, 2)
	ss.checkHl(t, lns[1][0], StyleComment, 0, 3)
	ss.checkHl(t, lns[2][0], StyleComment, 0, 2)
	ss.checkHl(t, lns[2][1], StyleSymb1, 4, 5)
}

func (ss *SyntaxSuite) checkHl(t *C, h Highlight, id StyleId, from, to int) {
	assert.Eq(t, h.Style, id)
	assert.Eq(t, h.ColFrom, from)
	assert.Eq(t, h.ColTo, to)
}
