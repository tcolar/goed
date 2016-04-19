package client

import (
	"github.com/tcolar/goed/actions"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestCmdBarEnable(c *C) {
	res, err := Action(s.id, []string{"cmdbar_enable", "true"})
	c.Assert(err, IsNil)
	c.Assert(actions.Ar.CmdbarEnabled(), Equals, true)

	res, err = Action(s.id, []string{"cmdbar_enable", "false"})
	c.Assert(err, IsNil)
	c.Assert(actions.Ar.CmdbarEnabled(), Equals, false)

	c.Assert(len(res), Equals, 0)
}

func (s *ApiSuite) TestCmdBarToggle(c *C) {
	res, err := Action(s.id, []string{"cmdbar_toggle"})
	c.Assert(err, IsNil)
	c.Assert(actions.Ar.CmdbarEnabled(), Equals, true)

	res, err = Action(s.id, []string{"cmdbar_toggle"})
	c.Assert(err, IsNil)
	c.Assert(actions.Ar.CmdbarEnabled(), Equals, false)

	c.Assert(len(res), Equals, 0)
}
