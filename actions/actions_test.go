package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/core"
)

func init() {
	core.Testing = true
	core.Bus = NewActionBus()
	go core.Bus.Start()
}

// Note: Most actions are tested/exercised via the api tests in api/client

func TestUndo(t *testing.T) {
	b := core.Bus
	i := 7
	j := 1
	v1 := int64(1)
	v2 := int64(2)
	Undo(v1)
	Redo(v1)
	assert.Equal(t, i, 7)
	add(v1, &i, 3)
	assert.Equal(t, i, 10)
	Redo(v1)
	b.Flush()
	assert.Equal(t, i, 10, "test1")
	Undo(v1)
	b.Flush()
	assert.Equal(t, i, 7, "test1")
	Undo(v1)
	add(v1, &i, 9)
	add(v1, &i, 11)
	add(v2, &j, 17)
	assert.Equal(t, i, 27) // 7 +9 + 11
	assert.Equal(t, j, 18) // 1 + 17

	Undo(v2)
	b.Flush()
	assert.Equal(t, j, 1)
	assert.Equal(t, i, 27)
	Undo(v2)
	b.Flush()
	assert.Equal(t, j, 1)
	assert.Equal(t, i, 27)
	Redo(v2)
	b.Flush()
	assert.Equal(t, j, 18) // 1 + 17
	assert.Equal(t, i, 27)

	Undo(v1)
	b.Flush()
	assert.Equal(t, i, 16) // 7 + 9
	Undo(v1)
	b.Flush()
	assert.Equal(t, i, 7)

	add(v1, &i, 3)
	add(v1, &i, 5)
	assert.Equal(t, i, 15) // 7 + 3 +5
	Redo(v1)
	b.Flush()

	Undo(v1) // 7 + 3
	b.Flush()
	assert.Equal(t, i, 10)

	Undo(v1) // 7
	b.Flush()
	assert.Equal(t, i, 7)

	Redo(v1)
	b.Flush()
	assert.Equal(t, i, 10) // 7 + 3
	Redo(v1)
	b.Flush()
	assert.Equal(t, i, 15) // 7 + 3 + 5
	Redo(v1)

	UndoClear(v1)

	Undo(v1)
	Redo(v1)
}

func TestUndoLimit(t *testing.T) {
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
	assert.Equal(t, len(undos[v]), 3)
	assert.Equal(t, i, 39)
	Undo(v)
	assert.Equal(t, len(undos[v]), 2)
	assert.Equal(t, i, 26)
	Undo(v)
	assert.Equal(t, len(undos[v]), 1)
	assert.Equal(t, i, 15)
	Undo(v)
	assert.Equal(t, len(undos[v]), 0)
	assert.Equal(t, i, 8)
	Undo(v)
}

func add(v int64, i *int, inc int) {
	d(addAction{i, inc})
	UndoAdd(v, addAction{i, inc}, addAction{i, -inc})
	core.Bus.Flush()
}

type addAction struct {
	val *int
	inc int
}

func (a addAction) Run() {
	*a.val += a.inc
}
