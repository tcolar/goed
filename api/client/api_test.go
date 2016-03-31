package client

import (
	"fmt"
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

func TestApi(t *testing.T) {
	defer core.Cleanup()
	ed := core.Ed

	inst := core.InstanceId

	res, err := Action(inst, []string{"fuzz"})
	assert.NotNil(t, err)

	res, err = Action(inst, []string{"cmdbar_enable", "true"})
	assert.Nil(t, err)
	assert.True(t, ed.CmdOn())

	fmt.Println(res)
	/*
		body, err = get("/v1/cur_view")
		assert.Nil(t, err)
		assert.Equal(t, body, "1")

		body, err = get("/v1/view/1/title")
		assert.Nil(t, err)
		assert.Equal(t, body, "file1.txt")

		body, err = get("/v1/view/1/workdir")
		assert.Nil(t, err)
		d, _ := filepath.Abs("../test_data")
		assert.Equal(t, body, d)

		body, err = get("/v1/view/1/src_loc")
		p, _ := filepath.Abs("../test_data/file1.txt")
		assert.Nil(t, err)
		assert.Equal(t, body, p)

		body, err = get("/v1/view/1/dirty")
		assert.Nil(t, err)
		assert.Equal(t, body, "0")

		body, err = get("/v1/view/1/selections")
		assert.Nil(t, err)
		assert.Equal(t, body, "")

		s := core.Ed.CurView().Selections()
		sel := core.NewSelection(0, 0, 1, 9)
		sel2 := core.NewSelection(2, 2, 4, 5)
		*s = append(*s, *sel, *sel2)
		body, err = get("/v1/view/1/selections")
		assert.Nil(t, err)
		assert.Equal(t, body, "0 0 1 9\n2 2 4 5\n")

		body, err = get("/v1/view/1/line_count")
		assert.Nil(t, err)
		assert.Equal(t, body, "12")
	*/
}
