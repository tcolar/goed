package backend

import (
	"testing"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

var id int64 = 9999
var id2 int64 = 9998

func Test(t *testing.T) { TestingT(t) }

type BackendSuite struct {
}

var _ = Suite(&BackendSuite{})

func (bs *BackendSuite) SetUpSuite(t *C) {
	core.Testing = true
	core.InitHome(time.Now().Unix())
	core.Bus = actions.NewActionBus()
	go core.Bus.Start()
}

func (bs *BackendSuite) TestFileBackend(t *C) {
	var b *FileBackend
	var err error
	// TODO: copy it to a temp dir for save test
	b, err = NewFileBackend("../test_data/file1.txt", id)
	assert.Nil(t, err)

	// Test some backend_file internals
	r, sz, err := b.readRune()
	assert.Nil(t, err)
	assert.True(t, r == '1' && sz == 1)

	b.seek(0, 9)
	r, sz, err = b.readRune()
	assert.Nil(t, err)
	assert.True(t, r == '0' && sz == 1)

	err = b.seek(200, 100)
	assert.NotNil(t, err)

	_, _, err = b.readRune() // still at EOF
	assert.NotNil(t, err)

	// Test backend interface compliance
	b.bufferSize = 17 // use a smallish buffer to make things more interesting
	assert.Eq(t, b.BufferLoc(), BufferFile(id))
	bs.testBackend(t, b, id)

	err = b.Close()
	assert.Nil(t, err)

}

func (bs *BackendSuite) TestMemBackend(t *C) {
	b2, err := NewMemBackend("../test_data/file1.txt", id2)
	bs.testBackend(t, b2, id2)
	err = b2.Close()
	assert.Nil(t, err)
}

// test Backend API methods
func (bs *BackendSuite) testBackend(t *C, b core.Backend, id int64) {
	assert.Eq(t, b.LineCount(), 12)
	assert.Eq(t, b.SrcLoc(), "../test_data/file1.txt")

	s1 := b.Slice(0, 0, 0, 9)
	s := core.RunesToString(*s1.Text())
	assert.Eq(t, s, "1234567890")
	s = core.RunesToString(*b.Slice(3, 4, 3, 4).Text())
	assert.Eq(t, s, "E")
	s = core.RunesToString(*b.Slice(2, 1, 2, 3).Text())
	assert.Eq(t, s, "bcd")
	s = core.RunesToString(*b.Slice(6, 1, 6, 5).Text())
	assert.Eq(t, s, "βξδεφ")
	// Should be an "absolute" move.
	s = core.RunesToString(*b.Slice(0, 0, 0, 9).Text())
	assert.Eq(t, s, "1234567890")
	// actual rectangle slice
	expected := "567890\n\nefghijkl\nEFGHIJKL\n\n\nεφγηιςκλ\nΕΦΓΗΙςΚΛ"
	s6 := b.Slice(0, 4, 7, 11)
	s = core.RunesToString(*s6.Text())
	assert.Eq(t, s6.R1, 0)
	assert.Eq(t, s6.C1, 4)
	assert.Eq(t, s6.R2, 7)
	assert.Eq(t, s6.C2, 11)
	assert.Eq(t, s, expected)
	s = core.RunesToString(*b.Slice(9, 2, 9, 3).Text())
	assert.Eq(t, s, "ab")
	// "backward" and mostly out of bounds slice
	s = core.RunesToString(*b.Slice(11, 20, 11, 9).Text())
	assert.Eq(t, s, `"wide" runes`)

	bs.insertionTests(t, b)
	// TODO: test save etc ....
	// Test file MD5
}

func (bs *BackendSuite) insertionTests(t *C, b core.Backend) {
	whole := core.RunesToString(*b.Slice(0, 0, -1, -1).Text())

	// Some inserts
	bs.testInsertRm(t, b, "$", 0, 2, 2, "ab$cdefghijklmnopqrstuvwxyz")
	bs.testInsertRm(t, b, "\n", 1, 2, 2, "ab\ncdefghijklmnopqrstuvwxyz")
	bs.testInsertRm(t, b, "@\n", 1, 2, 3, "ab@\ncdefghijklmnopqrstuvwxyz")
	bs.testInsertRm(t, b, "\n*", 1, 3, 0, "ab\n*cdefghijklmnopqrstuvwxyz")
	bs.testInsertRm(t, b, "!\n!\n", 2, 3, 1, "ab!\n!\ncdefghijklmnopqrstuvwxyz")
	bs.testInsertRm(t, b, "\n-\n-", 2, 4, 0, "ab\n-\n-cdefghijklmnopqrstuvwxyz")
	bs.testInsertRm(t, b, "^\n^", 1, 3, 0, "ab^\n^cdefghijklmnopqrstuvwxyz")
	bs.testInsertRm(t, b, "#\n##\n\n#", 3, 5, 0, "ab#\n##\n\n#cdefghijklmnopqrstuvwxyz")
	whole2 := core.RunesToString(*b.Slice(0, 0, -1, -1).Text())
	assert.Eq(t, whole2, whole)
}

const testLine3 = "abcdefghijklmnopqrstuvwxyz"

func (bs *BackendSuite) testInsertRm(t *C, b core.Backend, add string, lns, rl, rc int, expected string) {
	lines := b.LineCount()
	err := b.Insert(2, 2, add)
	assert.Nil(t, err)
	s := core.RunesToString(*b.Slice(2, 0, 2+lns, 30).Text())
	assert.Eq(t, s, expected)
	assert.Eq(t, b.LineCount(), lines+lns)
	err = b.Remove(2, 2, rl, rc)
	assert.Nil(t, err)
	s = core.RunesToString(*b.Slice(2, 0, 2, 30).Text())
	assert.Eq(t, s, testLine3)
	assert.Eq(t, b.LineCount(), lines)
}
