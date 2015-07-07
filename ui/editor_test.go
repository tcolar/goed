package ui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/core"
)

func init() {
	core.Testing = true
	core.InitHome(time.Now().Unix())
	core.Ed = NewMockEditor()
	core.Ed.Start([]string{})
}

func TestQuitCheck(t *testing.T) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView()
	v2 := Ed.NewView()
	col := Ed.NewCol(1.0, []*View{v, v2})
	Ed.Cols = []*Col{col}
	then := time.Now()
	assert.True(t, Ed.QuitCheck(), "quitcheck1")
	v2.SetDirty(true)
	assert.False(t, Ed.QuitCheck(), "quitcheck2")
	assert.True(t, v2.lastCloseTs.After(then), "quitcheck ts")
	assert.True(t, Ed.QuitCheck(), "quitcheck3")
}

func TestStartMany(t *testing.T) {
	ed := NewMockEditor()
	ed.Start([]string{"./test_data", "./test_data/empty.txt", "./test_data/file1.txt"})
	assert.Equal(t, len(ed.Cols), 2)
	assert.Equal(t, len(ed.Cols[0].Views), 1)
	assert.Equal(t, len(ed.Cols[1].Views), 2)
}
