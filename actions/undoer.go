package actions

import (
	"fmt"
	"sync"

	"github.com/tcolar/goed/core"
)

// TODO  ? var maxUndos = 500

// viewId keyed map of undo actions
var undos map[int64][]actionTuple = map[int64][]actionTuple{}

// viewId keyed map of redo actions
var redos map[int64][]actionTuple = map[int64][]actionTuple{}

var lock sync.Mutex

// a do/undo combo
type actionTuple struct {
	do   core.Action
	undo core.Action
}

// TODO : group together quick succesive undos (insert a, insert b, insert c) + Flushing
// or group by alphanum sequence ??
func Undo(viewId int64) error {
	lock.Lock()
	defer lock.Unlock()
	tuples, found := undos[viewId]
	if !found || len(tuples) == 0 {
		return fmt.Errorf("Nothing to undo.")
	}
	tuple := tuples[len(tuples)-1]
	undos[viewId] = undos[viewId][:len(tuples)-1]
	redos[viewId] = append(redos[viewId], tuple)
	d(tuple.undo)
	return nil
}

func Redo(viewId int64) error {
	lock.Lock()
	defer lock.Unlock()
	tuples, found := redos[viewId]
	if !found || len(tuples) == 0 {
		return fmt.Errorf("Nothing to redo.")
	}
	tuple := tuples[len(tuples)-1]
	redos[viewId] = redos[viewId][:len(tuples)-1]
	undos[viewId] = append(undos[viewId], tuple)
	d(tuple.do)
	return nil
}

func UndoAdd(viewId int64, do, undo core.Action) {
	lock.Lock()
	defer lock.Unlock()
	delete(redos, viewId)
	undos[viewId] = append(undos[viewId], actionTuple{do, undo})
}

func UndoClear(viewId int64) {
	lock.Lock()
	defer lock.Unlock()
	delete(undos, viewId)
	delete(redos, viewId)
}
