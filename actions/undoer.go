package actions

import (
	"fmt"
	"sync"

	"github.com/tcolar/goed/core"
)

// TODO  ? var maxUndos = 500
// TODO : group together quick succesive undos (insert a, insert b, insert c) + Flushing

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

// or group by alphanum sequence ??
func Undo(viewId int64) error {
	action, err := func() (core.Action, error) {
		lock.Lock()
		defer lock.Unlock()
		tuples, found := undos[viewId]
		if !found || len(tuples) == 0 {
			return nil, fmt.Errorf("Nothing to undo.")
		}
		tuple := tuples[len(tuples)-1]
		undos[viewId] = undos[viewId][:len(tuples)-1]
		redos[viewId] = append(redos[viewId], tuple)
		return tuple.undo, nil
	}()
	if err != nil {
		return err
	}
	return action.Run()
}

func Redo(viewId int64) error {
	action, err := func() (core.Action, error) {
		lock.Lock()
		defer lock.Unlock()
		tuples, found := redos[viewId]
		if !found || len(tuples) == 0 {
			return nil, fmt.Errorf("Nothing to redo.")
		}
		tuple := tuples[len(tuples)-1]
		redos[viewId] = redos[viewId][:len(tuples)-1]
		undos[viewId] = append(undos[viewId], tuple)
		return tuple.do, nil
	}()
	if err != nil {
		return err
	}
	return action.Run()
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

// Dump prints out the undo/redo stack of a view, for debugging
func Dump(viewId int64) {
	fmt.Printf("Undos:\n")
	for _, u := range undos[viewId] {
		fmt.Printf("\t %#v\n", u)
	}
	fmt.Printf("Redos:\n")
	for _, r := range redos[viewId] {
		fmt.Printf("\t %#v\n", r)
	}
}
