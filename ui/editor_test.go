package ui

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	core.Testing = true
	core.InitHome(time.Now().Unix())
	core.Ed = NewMockEditor()
	core.Bus = actions.NewActionBus()
	go core.Bus.Start()
	core.Ed.Start([]string{})
}

func TestQuitCheck(t *testing.T) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v2 := Ed.NewView("")
	col := Ed.NewCol(1.0, []int64{v.Id(), v2.Id()})
	Ed.Cols = []*Col{col}
	then := time.Now()
	assert.True(t, Ed.QuitCheck(), "quitcheck1")
	v2.SetDirty(true)
	assert.False(t, Ed.QuitCheck(), "quitcheck2")
	assert.True(t, v2.lastCloseTs.After(then), "quitcheck ts")
	assert.True(t, Ed.QuitCheck(), "quitcheck3")
}
