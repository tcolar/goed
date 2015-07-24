package backend

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

var id int64 = 9999
var id2 int64 = 9998

func init() {
	core.Testing = true
	core.InitHome(time.Now().Unix())
	core.Bus = actions.NewActionBus()
	go core.Bus.Start()
}

func TestFileBackend(t *testing.T) {
	var b *FileBackend
	var err error
	// TODO: copy it to a temp dir for save test
	b, err = NewFileBackend("../test_data/file1.txt", id)
	assert.Nil(t, err, "newFileBackend")

	// Test some backend_file internals
	r, sz, err := b.readRune()
	assert.Nil(t, err, "readrune")
	assert.True(t, r == '1' && sz == 1, "r=1, sz=1")

	b.seek(0, 9)
	r, sz, err = b.readRune()
	assert.Nil(t, err, "readrune")
	assert.True(t, r == '0' && sz == 1, "r=0, sz=1")

	err = b.seek(200, 100)
	assert.NotNil(t, err, "Invalid seek")

	_, _, err = b.readRune() // still at EOF
	assert.NotNil(t, err, "Invalid offset read")

	// Test backend interface compliance
	b.bufferSize = 17 // use a smallish buffer to make things more interesting
	assert.Equal(t, b.BufferLoc(), BufferFile(id), "bufferLoc")
	testBackend(t, b, id)

	err = b.Close()
	assert.Nil(t, err, "close")

}

func TestMemBackend(t *testing.T) {
	b2, err := NewMemBackend("../test_data/file1.txt", id2)
	testBackend(t, b2, id2)
	err = b2.Close()
	assert.Nil(t, err, "close")
}

// test Backend API methods
func testBackend(t *testing.T, b core.Backend, id int64) {
	assert.Equal(t, b.LineCount(), 12, "lineCount")
	assert.Equal(t, b.SrcLoc(), "../test_data/file1.txt", "srcLoc")

	s1 := b.Slice(0, 0, 0, 9)
	s := core.RunesToString(*s1.Text())
	assert.Equal(t, s, "1234567890", "slice1")
	s = core.RunesToString(*b.Slice(3, 4, 3, 4).Text())
	assert.Equal(t, s, "E", "slice2")
	s = core.RunesToString(*b.Slice(2, 1, 2, 3).Text())
	assert.Equal(t, s, "bcd", "slice3")
	s = core.RunesToString(*b.Slice(6, 1, 6, 5).Text())
	assert.Equal(t, s, "βξδεφ", "slice4")
	// Should be an "absolute" move.
	s = core.RunesToString(*b.Slice(0, 0, 0, 9).Text())
	assert.Equal(t, s, "1234567890", "slice5")
	// actual rectangle slice
	expected := "567890\n\nefghijkl\nEFGHIJKL\n\n\nεφγηιςκλ\nΕΦΓΗΙςΚΛ"
	s6 := b.Slice(0, 4, 7, 11)
	s = core.RunesToString(*s6.Text())
	assert.Equal(t, s6.R1, 0, "slice6.R1")
	assert.Equal(t, s6.C1, 4, "slice6.C1")
	assert.Equal(t, s6.R2, 7, "slice6.R2")
	assert.Equal(t, s6.C2, 11, "slice6.C2")
	assert.Equal(t, s, expected, "slice6")
	s = core.RunesToString(*b.Slice(9, 2, 9, 3).Text())
	assert.Equal(t, s, "ab", "slice7")
	// "backward" and mostly out of bounds slice
	s = core.RunesToString(*b.Slice(11, 20, 11, 9).Text())
	assert.Equal(t, s, `"wide" runes`, "slice8")

	insertionTests(t, b)
	// TODO: test save etc ....
	// Test file MD5
}

func insertionTests(t *testing.T, b core.Backend) {
	whole := core.RunesToString(*b.Slice(0, 0, -1, -1).Text())

	// Some inserts
	testInsertRm(t, b, "$", 0, 2, 2, "ab$cdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "\n", 1, 2, 2, "ab\ncdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "@\n", 1, 2, 3, "ab@\ncdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "\n*", 1, 3, 0, "ab\n*cdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "!\n!\n", 2, 3, 1, "ab!\n!\ncdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "\n-\n-", 2, 4, 0, "ab\n-\n-cdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "^\n^", 1, 3, 0, "ab^\n^cdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "#\n##\n\n#", 3, 5, 0, "ab#\n##\n\n#cdefghijklmnopqrstuvwxyz")
	whole2 := core.RunesToString(*b.Slice(0, 0, -1, -1).Text())
	assert.Equal(t, whole2, whole, "whole")
}

const testLine3 = "abcdefghijklmnopqrstuvwxyz"

func testInsertRm(t *testing.T, b core.Backend, add string, lns, rl, rc int, expected string) {
	lines := b.LineCount()
	err := b.Insert(2, 2, add)
	assert.Nil(t, err, "insert "+add)
	s := core.RunesToString(*b.Slice(2, 0, 2+lns, 30).Text())
	assert.Equal(t, s, expected, "slice "+add)
	assert.Equal(t, b.LineCount(), lines+lns, "lineCount "+add)
	err = b.Remove(2, 2, rl, rc)
	assert.Nil(t, err, "remove "+add)
	s = core.RunesToString(*b.Slice(2, 0, 2, 30).Text())
	assert.Equal(t, s, testLine3, "rm "+add)
	assert.Equal(t, b.LineCount(), lines, "count "+add)
}
