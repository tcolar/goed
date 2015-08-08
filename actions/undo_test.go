package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/core"
)

func init() {
	core.Bus = NewActionBus()
	go core.Bus.Start()
}

func TestUndo(t *testing.T) {
	b := core.Bus
	i := 7
	j := 1
	v1 := int64(1)
	v2 := int64(2)
	assert.Error(t, Undo(v1), "Nothing to undo")
	assert.Error(t, Redo(v1), "Nothing to redo")
	assert.Equal(t, i, 7)
	add(v1, &i, 3)
	assert.Equal(t, i, 10)
	assert.Error(t, Redo(v1), "Nothing to redo")
	b.Flush()
	assert.Equal(t, i, 10, "test1")
	assert.Nil(t, Undo(v1), "Undo 1")
	b.Flush()
	assert.Equal(t, i, 7, "test1")
	assert.Error(t, Undo(v1), "Undo 2")
	add(v1, &i, 9)
	add(v1, &i, 11)
	add(v2, &j, 17)
	assert.Equal(t, i, 27) // 7 +9 + 11
	assert.Equal(t, j, 18) // 1 + 17

	assert.Nil(t, Undo(v2))
	b.Flush()
	assert.Equal(t, j, 1)
	assert.Equal(t, i, 27)
	assert.Error(t, Undo(v2))
	b.Flush()
	assert.Equal(t, j, 1)
	assert.Equal(t, i, 27)
	assert.Nil(t, Redo(v2))
	b.Flush()
	assert.Equal(t, j, 18) // 1 + 17
	assert.Equal(t, i, 27)

	assert.Nil(t, Undo(v1))
	b.Flush()
	assert.Equal(t, i, 16) // 7 + 9
	assert.Nil(t, Undo(v1))
	b.Flush()
	assert.Equal(t, i, 7)

	add(v1, &i, 3)
	add(v1, &i, 5)
	assert.Equal(t, i, 15) // 7 + 3 +5
	assert.Error(t, Redo(v1))
	b.Flush()

	assert.Nil(t, Undo(v1)) // 7 + 3
	b.Flush()
	assert.Equal(t, i, 10)

	assert.Nil(t, Undo(v1)) // 7
	b.Flush()
	assert.Equal(t, i, 7)

	assert.Nil(t, Redo(v1))
	b.Flush()
	assert.Equal(t, i, 10) // 7 + 3
	assert.Nil(t, Redo(v1))
	b.Flush()
	assert.Equal(t, i, 15) // 7 + 3 + 5
	assert.Error(t, Redo(v1))

	UndoClear(v1)

	assert.Error(t, Undo(v1))
	assert.Error(t, Redo(v1))
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
	assert.Nil(t, Undo(v))
	assert.Equal(t, len(undos[v]), 2)
	assert.Equal(t, i, 26)
	assert.Nil(t, Undo(v))
	assert.Equal(t, len(undos[v]), 1)
	assert.Equal(t, i, 15)
	assert.Nil(t, Undo(v))
	assert.Equal(t, len(undos[v]), 0)
	assert.Equal(t, i, 8)
	assert.Error(t, Undo(v))
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

func (a addAction) Run() error {
	*a.val += a.inc
	return nil
}
