package main

import (
	"fmt"
	"strings"
)

// Cmdbar widget
type Cmdbar struct {
	Widget
	Cmd     []rune
	History [][]rune // TBD
}

func (c *Cmdbar) Render() {
	Ed.FB(Ed.Theme.Cmdbar.Fg, Ed.Theme.Cmdbar.Bg)
	Ed.Fill(Ed.Theme.Cmdbar.Rune, c.x1, c.y1, c.x2, c.y2)
	if Ed.CmdOn {
		Ed.FB(Ed.Theme.CmdbarTextOn, Ed.Theme.Cmdbar.Bg)
		Ed.Str(c.x1, c.y1, fmt.Sprintf("> %s", string(c.Cmd)))
	} else {
		Ed.FB(Ed.Theme.CmdbarText, Ed.Theme.Cmdbar.Bg)
		Ed.Str(c.x1, c.y1, fmt.Sprintf("> %s", string(c.Cmd)))
	}
}

func (c *Cmdbar) RunCmd() {
	// TODO: This is temporary until I create real fs based events & actions
	s := string(c.Cmd)
	parts := strings.Fields(s)
	if len(parts) < 1 {
		return
	}
	switch parts[0] {
	case "s", "save":
		//TODO: might need a new id etc.../
		if Ed.CurView != nil {
			Ed.CurView.Save()
		}
	case "o", "open":
		if len(parts) < 2 {
			Ed.SetStatusErr("Missing file path")
			return
		}
		// if no active view, create one ??
		// if active view is dirty, create one ??
		OpenFile(parts[1], Ed.CurView)
		Ed.CmdOn = false
	case "h", "help":
		Ed.SetStatus("TBD help")
	case "nc", "newcol":
		Ed.SetStatus("TBD nc")
	case "nv", "newview":
		Ed.SetStatus("TBD nv")
	case "dc", "delcol":
		Ed.SetStatus("TBD dc")
	case "dv", "delview":
		Ed.SetStatus("TBD dv")
	case "gf", "gofmt":
		Ed.SetStatus("TBD gofmt")
	default:
		Ed.SetStatusErr("Unexpected command " + parts[0])
	}
}
