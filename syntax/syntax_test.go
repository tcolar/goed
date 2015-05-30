package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/core"
)

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

func TestSyntax(t *testing.T) {
	hs := Highlights{}
	hs.Update(core.StringToRunes(testKw), ".go")
	lns := hs.Lines
	assert.Equal(t, len(lns), 1, "Kw test nb of lines")
	assert.Equal(t, len(lns[0]), 1, "Kw test line1, nb of highlights")
	checkHl(t, lns[0][0], StyleKw1, 0, 2, "Kw test line 0")
	hs.Update(core.StringToRunes(testSymb), ".go")
	lns = hs.Lines
	assert.Equal(t, len(lns), 1, "Symb test nb of lines")
	assert.Equal(t, len(lns[0]), 5, "Symb test line0, nb of highlights")
	checkHl(t, lns[0][0], StyleSymb1, 2, 3, "Symb test ':='")
	checkHl(t, lns[0][1], StyleSep1, 5, 5, "Symb test '('")
	checkHl(t, lns[0][2], StyleSymb3, 7, 7, "Symb test '*'")
	checkHl(t, lns[0][3], StyleSep1, 9, 9, "Symb test ')'")
	checkHl(t, lns[0][4], StyleSymb3, 10, 11, "Symb test '>>'")
	hs.Update(core.StringToRunes(testEsc), ".go")
	lns = hs.Lines
	assert.Equal(t, len(lns), 1, "Esc test nb of lines")
	assert.Equal(t, len(lns[0]), 2, "Esc test line0, nb of highlights")
	checkHl(t, lns[0][0], StyleSymb1, 4, 5, "Esc test ':='")
	checkHl(t, lns[0][1], StyleString, 7, 16, "Esc test string")
	hs.Update(core.StringToRunes(testSrc), ".go")
	lns = hs.Lines
	assert.Equal(t, len(lns), 11, "Src test nb of lines")
	checkHl(t, lns[0][0], StyleKw1, 0, 6, "Src test 'package'")
	checkHl(t, lns[1][0], StyleComment, 0, 12, "Src test comment")
	checkHl(t, lns[2][0], StyleKw1, 0, 3, "Src test 'type'")
	checkHl(t, lns[4][0], StyleComment, 0, 4, "Src test ML comment")
	checkHl(t, lns[5][0], StyleComment, 0, 2, "Src test ML comment")
	checkHl(t, lns[6][0], StyleComment, 0, 4, "Src test ML comment")
	assert.Equal(t, len(lns[7]), 4, "Src test len(line7)")
	assert.Equal(t, len(lns[8]), 8, "Src test len(line8)")
	checkHl(t, lns[9][0], StyleSymb1, 8, 9, "Src test ':='")
	checkHl(t, lns[10][0], StyleSep1, 0, 0, "Src test '}'")
	hs.Update(core.StringToRunes(testTop), ".go")
	lns = hs.Lines
	assert.Equal(t, len(lns), 3, "Top test nb of lines")
	checkHl(t, lns[0][0], StyleComment, 0, 2, "Top test comment 1")
	checkHl(t, lns[1][0], StyleComment, 0, 3, "Top test comment 2")
	checkHl(t, lns[2][0], StyleComment, 0, 2, "Top test comment 3")
	checkHl(t, lns[2][1], StyleSymb1, 4, 5, "Top test symb1")
}

func checkHl(t *testing.T, h Highlight, id StyleId, from, to int, msg string) {
	assert.Equal(t, h.Style, id, msg+" (style)")
	assert.Equal(t, h.ColFrom, from, msg+" (from)")
	assert.Equal(t, h.ColTo, to, msg+" (to)")
}
