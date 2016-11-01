package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui/widgets"
)

var _ core.Commander = (*Cmdbar)(nil)

// Cmdbar is the CommandBar widget
// It's sort of a temporary crutch as of now.
type Cmdbar struct {
	widgets.BaseWidget
	cmd        []rune
	history    [][]rune
	cursorX    int
	historyPos int
}

func (c *Cmdbar) Render() {
	ed := core.Ed
	t := ed.Theme()
	ed.TermFB(t.Cmdbar.Fg, t.Cmdbar.Bg)
	y1, x1, y2, x2 := c.Bounds()
	ed.TermFill(t.Cmdbar.Rune, y1, x1, y2, x2)
	fg := t.CmdbarText
	bg := t.Cmdbar.Bg
	if ed.CmdOn() {
		fg = t.CmdbarTextOn
	}
	ed.TermFB(fg, bg)
	for i, r := range "> " + string(c.cmd) + " " {
		if ed.CmdOn() && i == c.cursorX+2 {
			ed.TermFB(t.FgCursor, t.BgCursor)
		}
		ed.TermChar(y1, x1+i, r)
		if ed.CmdOn() && i == c.cursorX+2 {
			ed.TermFB(fg, bg)
		}
	}
	ed.TermFB(t.CmdbarText, t.Cmdbar.Bg)
	ed.TermStr(y1, x2-11, fmt.Sprintf("|GoEd %s", core.Version))
}

func (e *Editor) CmdbarToggle() {
	e.cmdOn = !e.cmdOn
	if e.cmdOn {
		e.Cmdbar.historyPos = 0
	}
}

func (c *Cmdbar) Backspace() {
	if c.cursorX <= 0 {
		return
	}
	c.cmd = append(c.cmd[:c.cursorX-1], c.cmd[c.cursorX:]...)
	c.cursorX--
}

func (c *Cmdbar) Clear() {
	c.cmd = []rune{}
	c.cursorX = 0
}

func (c *Cmdbar) Delete() {
	if c.cursorX >= len(c.cmd) {
		return
	}
	if c.cursorX == len(c.cmd)-1 {
		c.cmd = c.cmd[:c.cursorX+1]
	} else {
		c.cmd = append(c.cmd[:c.cursorX+1], c.cmd[c.cursorX+2:]...)
	}
	c.cursorX--
}

func (c *Cmdbar) Insert(s string) {
	c.cmd = append(c.cmd[:c.cursorX], append([]rune(s), c.cmd[c.cursorX:]...)...)
	c.cursorX += len(s)
}

func (c *Cmdbar) CursorMvmt(m core.CursorMvmt) {
	switch m {
	case core.CursorMvmtLeft:
		if c.cursorX > 0 {
			c.cursorX--
		}
	case core.CursorMvmtRight:
		if c.cursorX < len(c.cmd) {
			c.cursorX++
		}
	case core.CursorMvmtUp:
		if c.historyPos >= len(c.history)-1 {
			return
		}
		c.historyPos++
		c.cmd = c.history[len(c.history)-c.historyPos]
	case core.CursorMvmtDown:
		if c.historyPos <= 1 {
			return
		}
		c.historyPos--
		c.cmd = c.history[len(c.history)-c.historyPos]
	}
}

func (c *Cmdbar) NewLine() { // run the command
	// TODO: Temporary hard coded commands until I create real fs based events & actions
	s := string(c.cmd)
	parts := strings.Fields(s)
	if len(parts) < 1 {
		return
	}
	args := parts[1:]
	var err error
	switch parts[0] {
	case "o", "open":
		err = c.open(args)
	case ":", "line":
		c.line(args)
	case "/", "search":
		if len(c.cmd) < 2 {
			break
		}
		query := string(c.cmd[2:])
		c.Search(query)
	default:
		exec(parts, false)
	}

	c.history = append(c.history, c.cmd)

	if err != nil {
		actions.Ar.EdSetStatusErr(err.Error())
	} else {
		actions.Ar.CmdbarClear()
		actions.Ar.CmdbarEnable(false)
	}
	actions.Ar.EdRender()
}

func (c *Cmdbar) open(args []string) error {
	if len(args) < 1 {
		// try to expand a location from the current view
		return fmt.Errorf("No path provided")
	}
	ed := core.Ed.(*Editor)
	v := ed.NewView(args[0])
	ed.InsertViewSmart(v)
	cv := ed.CurView()
	_, err := ed.Open(args[0], cv.Id(), cv.WorkDir(), true)
	if err != nil {
		return err
	}
	ed.ViewActivate(cv.Id())
	return nil
}

func (c *Cmdbar) line(args []string) {
	ed := core.Ed.(*Editor)
	if len(args) < 0 {
		ed.SetStatusErr("Expected a line number argument.")
		return
	}
	l, err := strconv.Atoi(args[0])
	if err != nil {
		ed.SetStatusErr("Expected a line number argument.")
		return
	}
	v := ed.ViewById(ed.CurViewId())
	if v != nil {
		actions.Ar.ViewMoveCursor(ed.CurViewId(), l-v.CurLine()-1, 0, false)
	}
}

func (c *Cmdbar) Search(query string) {
	exec([]string{"grep", "-rni", "--color", query, "."}, false)
}
