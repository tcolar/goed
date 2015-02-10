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

func TestBackend(t *testing.T) {
	var b *FileBackend
	var err error
	v := Ed.NewView()
	// TODO: copy it to a temp dir
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
	testBackend(t, b)

	err = b.Close()
	assert.Nil(t, err, "close")

}

// test Backend API methods
func testBackend(t *testing.T, b Backend) {
	assert.Equal(t, b.LineCount(), 12, "lineCount")
	assert.Equal(t, b.SrcLoc(), "test_data/file1.txt", "srcLoc")
	assert.Equal(t, b.BufferLoc(), Ed.BufferFile(1), "bufferLoc")

	s := Ed.RunesToString(b.Slice(1, 1, 1, 10))
	assert.Equal(t, s, "1234567890", "slice1")
	s = Ed.RunesToString(b.Slice(4, 5, 4, 5))
	assert.Equal(t, s, "E", "slice2")
	s = Ed.RunesToString(b.Slice(3, 2, 3, 4))
	assert.Equal(t, s, "bcd", "slice3")
	s = Ed.RunesToString(b.Slice(7, 2, 7, 6))
	assert.Equal(t, s, "βξδεφ", "slice4")
	// Should be an "absolute" move.
	s = Ed.RunesToString(b.Slice(1, 1, 1, 10))
	assert.Equal(t, s, "1234567890", "slice5")
	// actual rectangle slice
	expected := "567890\n\nefghijkl\nEFGHIJKL\n\n\nεφγηιςκλ\nΕΦΓΗΙςΚΛ"
	s = Ed.RunesToString(b.Slice(1, 5, 8, 12))
	assert.Equal(t, s, expected, "slice6")
	s = Ed.RunesToString(b.Slice(10, 3, 10, 4))
	assert.Equal(t, s, "ab", "slice7")
	// "backward" and mostly out of bounds slice
	s = Ed.RunesToString(b.Slice(20, 21, 12, 10))
	assert.Equal(t, s, `"wide" runes`, "slice8")

	insertionTests(t, b)

	// TODO: test save etc ....
	// Test file MD5
}

func insertionTests(t *testing.T, b Backend) {
	lines := b.LineCount()
	whole := Ed.RunesToString(b.Slice(1, 1, 100, 100))

	// Some inserts
	err := b.Insert(3, 3, "^\n^")
	assert.Nil(t, err, "insert")
	assert.Equal(t, b.LineCount(), lines+1, "lineCount")
	s := Ed.RunesToString(b.Slice(3, 1, 4, 30))
	expected := "ab^\n^cdefghijklmnopqrstuvwxyz"
	assert.Equal(t, s, expected, "slice")

	err = b.Remove(3, 3, "^\n^")
	assert.Nil(t, err, "remove")

	whole2 := Ed.RunesToString(b.Slice(1, 1, -1, -1))
	assert.Equal(t, whole2, whole, "whole")
}
