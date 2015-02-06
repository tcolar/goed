package main

import (
	"log"
	"testing"
)

func init() {
	Ed = &Editor{}
	Ed.initHome()
}

const bufferId = 9999

func Test(t *testing.T) {
	var b *FileBackend
	var err error
	// TODO: copy it to a temp dir
	b, err = NewFileBackend("test_data/file1.txt", bufferId)
	if err != nil {
		panic(err)
	}

	// Test some backend_file internals
	r, sz, err := b.readRune()
	if r != '1' || sz != 1 {
		log.Fatalf("Expected '1' got %s (%d);", string(r), sz)
	}
	b.seek(1, 10)
	r, sz, err = b.readRune()
	if r != '0' || sz != 1 {
		log.Fatalf("Expected '0' got %s (%d);", string(r), sz)
	}
	err = b.seek(200, 100)
	if err == nil {
		log.Fatalf("Expected error on seek")
	}
	_, _, err = b.readRune() // still at EOF
	if err == nil {
		log.Fatalf("Expected error on readRune")
	}

	// Test backend interface compliance
	b.bufferSize = 17 // use a smallish buffer to make things more intresting
	testBackend(t, b)

	err = b.Close()
	if err != nil {
		panic(err)
	}
}

// test Backend API methods
func testBackend(t *testing.T, b Backend) {
	if b.LineCount() != 10 {
		log.Fatalf("Exepected 10 lines, got %d ", b.LineCount())
	}
	if b.SrcLoc() != "test_data/file1.txt" {
		log.Fatalf("srcloc : %s", b.SrcLoc)
	}
	if b.BufferLoc() != Ed.BufferFile(bufferId) {
		log.Fatalf("bufferloc : %s", b.BufferLoc)
	}
	s := Ed.RunesToString(b.Slice(1, 1, 1, 10))
	if s != "1234567890" {
		log.Fatalf("Expected '1234567890' got %s ", s)
	}
	s = Ed.RunesToString(b.Slice(4, 5, 4, 5))
	if s != "E" {
		log.Fatalf("Expected 'E' got %s ", s)
	}
	s = Ed.RunesToString(b.Slice(3, 2, 3, 4))
	if s != "bcd" {
		log.Fatalf("Expected 'bcd' got %s ", s)
	}
	s = Ed.RunesToString(b.Slice(7, 2, 7, 6))
	if s != "βξδεφ" { // 2 columns each
		log.Fatalf("Expected 'βξδεφ' got %s ", s)
	}
	// Should be an "absolute" move.
	s = Ed.RunesToString(b.Slice(1, 1, 1, 10))
	if s != "1234567890" {
		log.Fatalf("Expected '1234567890' got %s ", s)
	}
	// actual rectangle slice
	expected := "567890\n\nefghijkl\nEFGHIJKL\n\n\nεφγηιςκλ\nΕΦΓΗΙςΚΛ"
	s = Ed.RunesToString(b.Slice(1, 5, 8, 12))
	if s != expected {
		log.Fatalf("Expected %s got %s ", expected, s)
	}
	// "backward" and mostly out of bounds slice
	s = Ed.RunesToString(b.Slice(20, 21, 10, 10))
	if s != `"wide" runes` {
		log.Fatalf("Expected '\"wide\" runes' got '%s' ", s)
	}

	insertionTests(t, b)

	// TODO: test save etc ....
}

func insertionTests(t *testing.T, b Backend) {
	lines := b.LineCount()
	whole := Ed.RunesToString(b.Slice(1, 1, 100, 100))

	// Some inserts
	err := b.Insert(3, 3, "^\n^")
	if err != nil {
		panic(err)
	}
	if b.LineCount() != lines+1 {
		log.Fatalf("Expected 2 more lines after insertion, got %d", b.LineCount()-lines)
	}
	s := Ed.RunesToString(b.Slice(3, 1, 4, 30))
	expected := "ab^\n^cdefghijklmnopqrstuvwxyz"
	if s != expected {
		log.Fatalf("Expected '%s' got '%s' ", expected, s)
	}

	err = b.Remove(3, 3, "^\n^")
	if err != nil {
		panic(err)
	}

	whole2 := Ed.RunesToString(b.Slice(1, 1, 100, 100))
	if whole2 != whole {
		log.Fatalf("File has changed :\n---\n%s\n---\n%s", whole, whole2)
	}
}
