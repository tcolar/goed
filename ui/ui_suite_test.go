package ui

import (
	"math/rand"
	"testing"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type UiSuite struct {
}

var _ = Suite(&UiSuite{})

func (s *UiSuite) SetUpSuite(c *C) {
	rand.Seed(time.Now().UTC().UnixNano())
	core.Testing = true
	core.InitHome(time.Now().Unix())
	core.Ed = NewMockEditor()
	core.Bus = actions.NewActionBus()
	// Note: not starting the action bus in UI tests so not to have to worry
	// about potential races internally.
	//go core.Bus.Start()
	core.Ed.Start([]string{})
}
