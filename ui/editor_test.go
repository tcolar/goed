package ui

import (
	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func (us *UiSuite) TestRingBuffer(t *C) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	title := "foo"
	b, _ := backend.NewMemBackendCmd([]string{"echo", "1"}, ".", v.Id(), &title, false)
	b.MaxRows = 5
	fg := core.NewStyle(1)
	bg := core.NewStyle(2)

	l, r := b.Overwrite(0, 0, "1", fg, bg)
	assert.Eq(t, b.Head(), 0)
	assert.Eq(t, l, 0)
	assert.Eq(t, r, 1)
	s1 := b.Slice(0, 0, 0, -1)
	s := core.RunesToString(*s1.Text())
	assert.Eq(t, s, "1")
	c1, c2 := b.ColorAt(0, 0)
	assert.Eq(t, c1, fg)
	assert.Eq(t, c2, bg)

	fg = core.NewStyle(98)
	bg = core.NewStyle(99)

	l, r = b.Overwrite(1, 0, "22\n333\n4444\n55555\n666666", fg, bg)
	assert.Eq(t, l, 4)
	assert.Eq(t, r, 6)
	assert.Eq(t, b.Head(), 1)
	s1 = b.Slice(0, 0, 0, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "22")
	s1 = b.Slice(4, 0, 4, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "666666")
	c1, c2 = b.ColorAt(0, 0)
	assert.Eq(t, c1, fg)
	assert.Eq(t, c2, bg)

	l, r = b.Overwrite(5, 0, "7", fg, bg)
	assert.Eq(t, l, 4)
	assert.Eq(t, r, 1)
	assert.Eq(t, b.Head(), 2)
	s1 = b.Slice(0, 0, 0, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "333")
	s1 = b.Slice(4, 0, 4, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "7")
	c1, c2 = b.ColorAt(0, 0)
	assert.Eq(t, c1, fg)
	assert.Eq(t, c2, bg)

	s1 = b.Slice(0, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "333\n4444\n55555\n666666\n7")

	s1 = b.Slice(0, 0, 100, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "333\n4444\n55555\n666666\n7")

	s1 = b.Slice(1, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "4444\n55555\n666666\n7")

	s1 = b.Slice(2, 0, 100, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "55555\n666666\n7")

	l, r = b.Overwrite(4, 0, "77\n88", fg, bg)
	s1 = b.Slice(0, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "4444\n55555\n666666\n77\n88")

	l, r = b.Overwrite(4, 0, "888", fg, bg)
	s1 = b.Slice(0, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Eq(t, s, "4444\n55555\n666666\n77\n888")
}
