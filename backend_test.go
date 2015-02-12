package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	Ed = &Editor{
		testing: true,
	}
	Ed.initHome()
}

func TestFileBackend(t *testing.T) {
	var b *FileBackend
	var err error
	v := Ed.NewView()
	// TODO: copy it to a temp dir for save test
	b, err = Ed.NewFileBackend("test_data/file1.txt", v.Id)
	assert.Nil(t, err, "newFileBackend")

	// Test some backend_file internals
	r, sz, err := b.readRune()
	assert.Nil(t, err, "readrune")
	assert.True(t, r == '1' && sz == 1, "=1r, sz=1")

	b.seek(1, 10)
	r, sz, err = b.readRune()
	assert.Nil(t, err, "readrune")
	assert.True(t, r == '0' && sz == 1, "r=0, sz=1")

	err = b.seek(200, 100)
	assert.NotNil(t, err, "Invalid seek")

	_, _, err = b.readRune() // still at EOF
	assert.NotNil(t, err, "Invalid offset read")

	// Test backend interface compliance
	b.bufferSize = 17 // use a smallish buffer to make things more interesting
	assert.Equal(t, b.BufferLoc(), Ed.BufferFile(v.Id), "bufferLoc")
	testBackend(t, b, v.Id)

	err = b.Close()
	assert.Nil(t, err, "close")

}

func TestMemBackend(t *testing.T) {
	v2 := Ed.NewView()
	b2, err := Ed.NewMemBackend("test_data/file1.txt", v2.Id)
	testBackend(t, b2, v2.Id)
	err = b2.Close()
	assert.Nil(t, err, "close")
}

// test Backend API methods
func testBackend(t *testing.T, b Backend, id int) {
	assert.Equal(t, b.LineCount(), 12, "lineCount")
	assert.Equal(t, b.SrcLoc(), "test_data/file1.txt", "srcLoc")

	s1 := b.Slice(1, 1, 1, 10)
	s := Ed.RunesToString(s1.text)
	assert.Equal(t, s, "1234567890", "slice1")
	s = Ed.RunesToString(b.Slice(4, 5, 4, 5).text)
	assert.Equal(t, s, "E", "slice2")
	s = Ed.RunesToString(b.Slice(3, 2, 3, 4).text)
	assert.Equal(t, s, "bcd", "slice3")
	s = Ed.RunesToString(b.Slice(7, 2, 7, 6).text)
	assert.Equal(t, s, "βξδεφ", "slice4")
	// Should be an "absolute" move.
	s = Ed.RunesToString(b.Slice(1, 1, 1, 10).text)
	assert.Equal(t, s, "1234567890", "slice5")
	// actual rectangle slice
	expected := "567890\n\nefghijkl\nEFGHIJKL\n\n\nεφγηιςκλ\nΕΦΓΗΙςΚΛ"
	s6 := b.Slice(1, 5, 8, 12)
	s = Ed.RunesToString(s6.text)
	assert.Equal(t, s6.r1, 1, "slice6.r1")
	assert.Equal(t, s6.c1, 5, "slice6.c1")
	assert.Equal(t, s6.r2, 8, "slice6.r2")
	assert.Equal(t, s6.c2, 12, "slice6.c2")
	assert.Equal(t, s, expected, "slice6")
	s = Ed.RunesToString(b.Slice(10, 3, 10, 4).text)
	assert.Equal(t, s, "ab", "slice7")
	// "backward" and mostly out of bounds slice
	s = Ed.RunesToString(b.Slice(12, 21, 12, 10).text)
	assert.Equal(t, s, `"wide" runes`, "slice8")

	insertionTests(t, b)

	// TODO: test save etc ....
	// Test file MD5
}

func insertionTests(t *testing.T, b Backend) {
	whole := Ed.RunesToString(b.Slice(1, 1, -1, -1).text)

	// Some inserts
	testInsertRm(t, b, "$", 0, "ab$cdefghijklmnopqrstuvwxyz")
	//testInsertRm(t, b, "\n", 1, "ab\ncdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "@\n", 1, "ab@\ncdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "\n@", 1, "ab\n@cdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "^\n^", 1, "ab^\n^cdefghijklmnopqrstuvwxyz")
	testInsertRm(t, b, "#\n##\n\n#", 3, "ab#\n##\n\n#cdefghijklmnopqrstuvwxyz")

	whole2 := Ed.RunesToString(b.Slice(1, 1, -1, -1).text)
	assert.Equal(t, whole2, whole, "whole")
}

const testLine3 = "abcdefghijklmnopqrstuvwxyz"

func testInsertRm(t *testing.T, b Backend, add string, lns int, expected string) {
	lines := b.LineCount()
	err := b.Insert(3, 3, add)
	assert.Nil(t, err, "insert "+add)
	s := Ed.RunesToString(b.Slice(3, 1, 3+lns, 30).text)
	assert.Equal(t, s, expected, "slice "+add)
	assert.Equal(t, b.LineCount(), lines+lns, "lineCount "+add)
	err = b.Remove(3, 3, add)
	assert.Nil(t, err, "remove "+add)
	s = Ed.RunesToString(b.Slice(3, 1, 3, 30).text)
	assert.Equal(t, s, testLine3, "rm "+add)
	assert.Equal(t, b.LineCount(), lines, "count "+add)
}
