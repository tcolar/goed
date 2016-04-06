package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
)

func TestCmdBarEnable(t *testing.T) {
	res, err := Action(id, []string{"foobar"})
	assert.NotNil(t, err)

	res, err = Action(id, []string{"cmdbar_enable", "true"})
	assert.Nil(t, err)
	assert.True(t, actions.Ar.CmdbarEnabled())

	res, err = Action(id, []string{"cmdbar_enable", "false"})
	assert.Nil(t, err)
	assert.False(t, actions.Ar.CmdbarEnabled())

	assert.Equal(t, len(res), 0)
}

func TestCmdBarToggle(t *testing.T) {
	res, err := Action(id, []string{"cmdbar_toggle"})
	assert.Nil(t, err)
	assert.True(t, actions.Ar.CmdbarEnabled())

	res, err = Action(id, []string{"cmdbar_toggle"})
	assert.Nil(t, err)
	assert.False(t, actions.Ar.CmdbarEnabled())

	assert.Equal(t, len(res), 0)
}
