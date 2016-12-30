package event

import (
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

type FileWatcher struct {
	sync.Mutex
	watcher  *fsnotify.Watcher
	watchMap map[string][]int64 // map of watched "path" to view(s)
	done     chan struct{}
}

func NewFileWatcher() *FileWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	return &FileWatcher{
		watcher:  w,
		done:     make(chan struct{}),
		watchMap: map[string][]int64{},
	}
}

func (w *FileWatcher) Start() {
	for {
		select {
		case <-w.done:
			return
		case event := <-w.watcher.Events:
			actions.Ar.EdFileEvent(core.FileOp(event.Op), event.Name)
		}
	}
}

func (w *FileWatcher) Stop() {
	w.Lock()
	defer w.Unlock()
	for loc, _ := range w.watchMap {
		w.watcher.Remove(loc)
	}
	close(w.done)
}

func (w *FileWatcher) Watch(vid int64, loc string) {
	if w == nil {
		return
	}
	loc, err := filepath.Abs(loc)
	if err != nil {
		return
	}
	w.Lock()
	defer w.Unlock()
	if m, found := w.watchMap[loc]; found {
		// already watching this loc
		idx := -1
		for i, id := range m {
			if id == vid {
				idx = i
			}
		}
		if idx != -1 {
			return // already watching this loc/view
		}
		// already watching but for a different view
		w.watchMap[loc] = append(m, vid)
		return
	}
	w.watchMap[loc] = []int64{vid}
	w.watcher.Add(loc)
}

func (w *FileWatcher) Unwatch(vid int64, loc string) {
	if w == nil {
		return
	}
	loc, err := filepath.Abs(loc)
	if err != nil {
		return
	}
	w.Lock()
	defer w.Unlock()
	if m, found := w.watchMap[loc]; found {
		idx := -1
		for i, id := range m {
			if id == vid {
				idx = i
			}
		}
		if idx == -1 {
			return // not watching this loc/view
		}
		// remove this view for this loc
		w.watchMap[loc] = append(m[:idx], m[idx+1:]...)
		// if no views left for this loc, no longer need to watch it.
		if len(w.watchMap[loc]) == 0 {
			delete(w.watchMap, loc)
			w.watcher.Remove(loc)
		}
	}
}
