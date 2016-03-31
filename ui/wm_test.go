package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/core"
)

func TestWm(t *testing.T) {
	var err error
	var v, v2 core.Viewable
	Ed := core.Ed.(*Editor)
	// Cleanup leftovers of previous tests
	for _, c := range Ed.Cols {
		for _, v := range c.Views {
			Ed.DelView(v, true)
		}
	}
	v = Ed.ViewById(Ed.Cols[0].Views[0])
	Ed.views = map[int64]*View{}
	Ed.views[v.Id()] = viewCast(v)

	assert.Equal(t, len(Ed.Cols), 1)
	c := Ed.Cols[0]
	assert.Equal(t, len(c.Views), 1)
	v = Ed.ViewById(c.Views[0])
	vid, err := Ed.Open("file1.txt", v.Id(), "../test_data", false)
	v = Ed.ViewById(vid)
	assertBounds(t, v, 1, 0, 23, 49)
	assert.Nil(t, err)
	assert.Equal(t, len(c.Views), 1)

	v2 = Ed.NewFileView("../test_data/no_eol.txt")
	Ed.InsertView(viewCast(v2), viewCast(v), 0.5)
	assert.Equal(t, len(c.Views), 2)
	assertBounds(t, v, 1, 0, 11, 49)
	assertBounds(t, v2, 12, 0, 23, 49)
	assert.IsType(t, Ed.WidgetAt(0, 0), Ed.Cmdbar)
	assert.Equal(t, Ed.WidgetAt(1, 0), v)
	assert.Equal(t, Ed.WidgetAt(1, 49), v)
	assert.Equal(t, Ed.WidgetAt(1, 0), v)
	assert.Equal(t, Ed.WidgetAt(20, 0), v2)
	assert.Equal(t, Ed.WidgetAt(20, 49), v2)
	assert.IsType(t, Ed.WidgetAt(24, 0), Ed.Statusbar)
	assert.Equal(t, Ed.ViewIndex(c, v.Id()), 0)
	assert.Equal(t, Ed.ViewIndex(c, v2.Id()), 1)
	assert.Equal(t, Ed.ViewColumn(v2.Id()), c)
	assert.Nil(t, Ed.ViewById(0))
	assert.Equal(t, Ed.ViewById(v.Id()), v)
	assert.Equal(t, Ed.ViewByLoc(""), int64(-1))
	assert.Equal(t, Ed.ViewByLoc(v.Backend().SrcLoc()), v.Id())
	assert.Equal(t, Ed.CurView(), v)
	Ed.ViewActivate(v2.Id())
	assert.Equal(t, Ed.CurView(), v2)

	Ed.DelView(v2.Id(), true)
	assert.Equal(t, len(c.Views), 1)
	assert.Equal(t, Ed.WidgetAt(20, 0), v)

	c2 := Ed.AddCol(c, 0.5)
	assert.Equal(t, len(Ed.Cols), 2)
	_, err = Ed.Open("no_eol.txt", c2.Views[0], "../test_data", false)
	v2 = Ed.views[c2.Views[0]]
	assert.Equal(t, len(c.Views), 1)
	assert.Equal(t, len(c2.Views), 1)
	assertBounds(t, v, 1, 0, 23, 24)
	assertBounds(t, v2, 1, 25, 23, 49)
	v3 := Ed.AddView(viewCast(v2), 0.5)
	assert.Equal(t, len(Ed.Cols), 2)
	assert.Equal(t, len(c2.Views), 2)
	assert.Equal(t, Ed.WidgetAt(2, 30), v2)
	assert.Equal(t, Ed.WidgetAt(20, 30), v3)
	Ed.ViewMove(viewCast(v2).y1, viewCast(v2).x1, v3.y1+5, viewCast(v2).x1)
	assert.Equal(t, Ed.WidgetAt(2, 30), v3)
	assert.Equal(t, Ed.WidgetAt(20, 30), v2)

	v3.SetDirty(true)
	Ed.DelViewCheck(v3.Id(), true)
	assert.Equal(t, len(c2.Views), 2) // dirty disallow it
	Ed.DelViewCheck(v3.Id(), true)
	assert.Equal(t, len(c2.Views), 1) // allowed second time
	assert.Equal(t, len(Ed.Cols), 2)

	Ed.DelCol(c2, true)
	assert.Equal(t, len(Ed.Cols), 1)

	Ed.DelCol(c, true)
	assert.Equal(t, len(Ed.Cols), 1) // can't remove last view/col
}

func assertBounds(t *testing.T, v core.Viewable, y1, x1, y2, x2 int) {
	b1, b2, b3, b4 := v.Bounds()
	assert.Equal(t, b1, y1)
	assert.Equal(t, b2, x1)
	assert.Equal(t, b3, y2)
	assert.Equal(t, b4, x2)
}
