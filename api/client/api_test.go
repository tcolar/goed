package client

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/api"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
)

var id int64

func init() {
	id = time.Now().Unix()
	core.Testing = true
	core.InitHome(id)
	core.Ed = ui.NewMockEditor()
	core.Bus = actions.NewActionBus()
	actions.RegisterActions()
	apiServer := api.Api{}
	apiServer.Start()
	go core.Bus.Start()
	core.Ed.Start([]string{"../test_data/file1.txt"})
}

func vidStr(vid int64) string {
	return fmt.Sprintf("%d", vid)
}

func TestNoSuchAction(t *testing.T) {
	res, err := Action(id, []string{"foobar"})
	assert.NotNil(t, err)
	assert.Equal(t, len(res), 0)
}

func TestEdit(t *testing.T) {
	done := false
	completed := make(chan struct{})
	go func() {
		err := Edit(id, "test_data", "fooedit")
		done = true
		assert.Nil(t, err)
		close(completed)
	}()
	vid := int64(-1)
	// view should open up and stay open until the view is closed
	// at which time the open action should be completed
	loc, _ := filepath.Abs("./test_data/fooedit")
	for vid == -1 {
		vid = actions.Ar.EdViewByLoc(loc)
		time.Sleep(100 * time.Millisecond)
	}
	assert.False(t, done)
	actions.Ar.EdDelView(vid, true)
	select {
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout waiting for edit to complete.")
	case <-completed: // good
	}
}

func TestOpen(t *testing.T) {
	err := Open(id, "test_data", "empty.txt")
	assert.Nil(t, err)
	loc, _ := filepath.Abs("./test_data/empty.txt")
	vid := actions.Ar.EdViewByLoc(loc)
	assert.NotEqual(t, vid, "-1")
	actions.Ar.EdDelView(vid, true)
}
