package client

import (
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/assert"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestCmdBarEnable(c *C) {
	res, err := Action(s.id, []string{"cmdbar_enable", "true"})
	assert.Nil(c, err)
	assert.True(c, actions.Ar.CmdbarEnabled())

	res, err = Action(s.id, []string{"cmdbar_enable", "false"})
	assert.Nil(c, err)
	assert.False(c, actions.Ar.CmdbarEnabled())

	assert.Eq(c, len(res), 0)
}

func (s *ApiSuite) TestCmdBarToggle(c *C) {
	res, err := Action(s.id, []string{"cmdbar_toggle"})
	assert.Nil(c, err)
	assert.True(c, actions.Ar.CmdbarEnabled())

	res, err = Action(s.id, []string{"cmdbar_toggle"})
	assert.Nil(c, err)
	assert.False(c, actions.Ar.CmdbarEnabled())

	assert.Eq(c, len(res), 0)
}
