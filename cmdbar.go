package main

import (
	"fmt"
	"strconv"
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
		Ed.OpenFile(parts[1], Ed.CurView)
		Ed.CmdOn = false
	case "h", "help":
		Ed.SetStatus("TBD help")
	case "nc", "newcol":
		// nc : newblank col
		// nc [file], new col, open file
		// nc 40 -> new col 40% width
		// nc 40 [file] -> new col 40% width, open file
		loc := ""
		pct := 50
		if len(parts) > 1 {
			p, err := strconv.Atoi(parts[1])
			if err == nil {
				pct = p
				if len(parts) > 2 {
					loc = strings.Join(parts[2:], " ")
				}
			} else {
				loc = strings.Join(parts[1:], " ")
			}
		}
		v := Ed.AddCol(float64(pct) / 100.0).Views[0]
		if len(loc) > 0 {
			Ed.OpenFile(loc, v)
		}
	case "nv", "newview":
		loc := ""
		pct := 50
		if len(parts) > 1 {
			p, err := strconv.Atoi(parts[1])
			if err == nil {
				pct = p
				if len(parts) > 2 {
					loc = strings.Join(parts[2:], " ")
				}
			} else {
				loc = strings.Join(parts[1:], " ")
			}
		}
		v := Ed.AddView(float64(pct) / 100.0)
		if len(loc) > 0 {
			Ed.OpenFile(loc, v)
		}
	case "dc", "delcol":
		Ed.DelCol()
	case "dv", "delview":
		Ed.DelView()
	case "gf", "gofmt":
		Ed.SetStatus("TBD gofmt")
	default:
		Ed.SetStatusErr("Unexpected command " + parts[0])
	}
}
