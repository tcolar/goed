package client

import (
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
var ed core.Editable

func init() {
	id = time.Now().Unix()
	core.Testing = true
	core.InitHome(id)
	core.Ed = ui.NewMockEditor()
	ed = core.Ed
	core.Bus = actions.NewActionBus()
	actions.RegisterActions()
	apiServer := api.Api{}
	apiServer.Start()
	go core.Bus.Start()
	core.Ed.Start([]string{"../test_data/file1.txt"})
}

func TestCmdBarEnable(t *testing.T) {

	res, err := Action(id, []string{"foobar"})
	assert.NotNil(t, err)

	res, err = Action(id, []string{"cmdbar_enable", "true"})
	assert.Nil(t, err)
	assert.True(t, ed.CmdOn())

	res, err = Action(id, []string{"cmdbar_enable", "false"})
	assert.Nil(t, err)
	assert.False(t, ed.CmdOn())

	assert.Equal(t, len(res), 0)
}

func TestCmdBarToggle(t *testing.T) {
	res, err := Action(id, []string{"cmdbar_toggle"})
	assert.Nil(t, err)
	assert.True(t, ed.CmdOn())

	res, err = Action(id, []string{"cmdbar_toggle"})
	assert.Nil(t, err)
	assert.False(t, ed.CmdOn())

	assert.Equal(t, len(res), 0)
}

func TestOpen(t *testing.T) {
	err := Open(id, "test_data", "empty.txt")
	assert.Nil(t, err)
	loc, _ := filepath.Abs("./test_data/empty.txt")
	vid := ed.ViewByLoc(loc)
	assert.NotEqual(t, vid, "-1")
	ed.DelView(vid, true)
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
		vid = ed.ViewByLoc(loc)
		time.Sleep(100 * time.Millisecond)
	}
	assert.False(t, done)
	ed.DelView(vid, true)
	select {
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout waiting for edit to complete.")
	case <-completed: // good
	}
}
