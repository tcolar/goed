package actions

import (
	"testing"

	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ActionSuite struct {
}

var _ = Suite(&ActionSuite{})

func (s *ActionSuite) SetUpSuite(c *C) {
	core.Testing = true
	core.Bus = NewActionBus()
	go core.Bus.Start()
}

// Note: Most actions are tested/exercised via the api tests in api/client

func (s *ActionSuite) TestUndo(t *C) {
	b := core.Bus
	i := 7
	j := 1
	v1 := int64(1)
	v2 := int64(2)
	Undo(v1)
	Redo(v1)
	assert.Eq(t, i, 7)
	add(v1, &i, 3)
	assert.Eq(t, i, 10)
	Redo(v1)
	b.Flush()
	assert.Eq(t, i, 10)
	Undo(v1)
	b.Flush()
	assert.Eq(t, i, 7)
	Undo(v1)
	add(v1, &i, 9)
	add(v1, &i, 11)
	add(v2, &j, 17)
	assert.Eq(t, i, 27) // 7 +9 + 11
	assert.Eq(t, j, 18) // 1 + 17

	Undo(v2)
	b.Flush()
	assert.Eq(t, j, 1)
	assert.Eq(t, i, 27)
	Undo(v2)
	b.Flush()
	assert.Eq(t, j, 1)
	assert.Eq(t, i, 27)
	Redo(v2)
	b.Flush()
	assert.Eq(t, j, 18) // 1 + 17
	assert.Eq(t, i, 27)

	Undo(v1)
	b.Flush()
	assert.Eq(t, i, 16) // 7 + 9
	Undo(v1)
	b.Flush()
	assert.Eq(t, i, 7)

	add(v1, &i, 3)
	add(v1, &i, 5)
	assert.Eq(t, i, 15) // 7 + 3 +5
	Redo(v1)
	b.Flush()

	Undo(v1) // 7 + 3
	b.Flush()
	assert.Eq(t, i, 10)

	Undo(v1) // 7
	b.Flush()
	assert.Eq(t, i, 7)

	Redo(v1)
	b.Flush()
	assert.Eq(t, i, 10) // 7 + 3
	Redo(v1)
	b.Flush()
	assert.Eq(t, i, 15) // 7 + 3 + 5
	Redo(v1)

	UndoClear(v1)

	Undo(v1)
	Redo(v1)
}

func (s *ActionSuite) TestUndoLimit(t *C) {
	v := int64(3)
	i := 0
	maxUndos = 3
	defer func() { maxUndos = 1000 }()
	add(v, &i, 3)
	add(v, &i, 5)
	add(v, &i, 7)
	add(v, &i, 11)
	add(v, &i, 13)
	core.Bus.Flush()
	assert.Eq(t, len(undos[v]), 3)
	assert.Eq(t, i, 39)
	Undo(v)
	assert.Eq(t, len(undos[v]), 2)
	assert.Eq(t, i, 26)
	Undo(v)
	assert.Eq(t, len(undos[v]), 1)
	assert.Eq(t, i, 15)
	Undo(v)
	assert.Eq(t, len(undos[v]), 0)
	assert.Eq(t, i, 8)
	Undo(v)
}

func add(v int64, i *int, inc int) {
	d(addAction{i, inc})
	UndoAdd(v, []core.Action{addAction{i, inc}}, []core.Action{addAction{i, -inc}})
	core.Bus.Flush()
}

type addAction struct {
	val *int
	inc int
}

func (a addAction) Run() {
	*a.val += a.inc
}
