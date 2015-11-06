package ui

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	core.Testing = true
	core.InitHome(time.Now().Unix())
	core.Ed = NewMockEditor()
	core.Bus = actions.NewActionBus()
	go core.Bus.Start()
	core.Ed.Start([]string{})
}

func TestQuitCheck(t *testing.T) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v2 := Ed.NewView("")
	col := Ed.NewCol(1.0, []int64{v.Id(), v2.Id()})
	Ed.Cols = []*Col{col}
	then := time.Now()
	assert.True(t, Ed.QuitCheck(), "quitcheck1")
	v2.SetDirty(true)
	assert.False(t, Ed.QuitCheck(), "quitcheck2")
	assert.True(t, v2.lastCloseTs.After(then), "quitcheck ts")
	assert.True(t, Ed.QuitCheck(), "quitcheck3")
}

func TestRingBuffer(t *testing.T) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	title := "foo"
	b, _ := backend.NewMemBackendCmd([]string{"echo", "1"}, ".", v.Id(), &title, false)
	b.MaxRows = 5
	fg := core.NewStyle(1)
	bg := core.NewStyle(2)

	l, r := b.Overwrite(0, 0, "1", fg, bg)
	assert.Equal(t, b.Head(), 0)
	assert.Equal(t, l, 0)
	assert.Equal(t, r, 1)
	s1 := b.Slice(0, 0, 0, -1)
	s := core.RunesToString(*s1.Text())
	assert.Equal(t, s, "1")
	c1, c2 := b.ColorAt(0, 0)
	assert.Equal(t, c1, fg)
	assert.Equal(t, c2, bg)

	fg = core.NewStyle(98)
	bg = core.NewStyle(99)

	l, r = b.Overwrite(1, 0, "22\n333\n4444\n55555\n666666", fg, bg)
	assert.Equal(t, l, 4)
	assert.Equal(t, r, 6)
	assert.Equal(t, b.Head(), 1)
	s1 = b.Slice(0, 0, 0, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "22")
	s1 = b.Slice(4, 0, 4, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "666666")
	c1, c2 = b.ColorAt(0, 0)
	assert.Equal(t, c1, fg)
	assert.Equal(t, c2, bg)

	l, r = b.Overwrite(5, 0, "7", fg, bg)
	assert.Equal(t, l, 4)
	assert.Equal(t, r, 1)
	assert.Equal(t, b.Head(), 2)
	s1 = b.Slice(0, 0, 0, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "333")
	s1 = b.Slice(4, 0, 4, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "7")
	c1, c2 = b.ColorAt(0, 0)
	assert.Equal(t, c1, fg)
	assert.Equal(t, c2, bg)

	s1 = b.Slice(0, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "333\n4444\n55555\n666666\n7")

	s1 = b.Slice(0, 0, 100, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "333\n4444\n55555\n666666\n7")

	s1 = b.Slice(1, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "4444\n55555\n666666\n7")

	s1 = b.Slice(2, 0, 100, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "55555\n666666\n7")

	l, r = b.Overwrite(4, 0, "77\n88", fg, bg)
	s1 = b.Slice(0, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "4444\n55555\n666666\n77\n88")

	l, r = b.Overwrite(4, 0, "888", fg, bg)
	s1 = b.Slice(0, 0, -1, -1)
	s = core.RunesToString(*s1.Text())
	assert.Equal(t, s, "4444\n55555\n666666\n77\n888")
}
