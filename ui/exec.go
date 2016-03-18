package ui

import (
	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
)

func exec(args []string, interactive bool) int64 {
	workDir := "."
	ed := core.Ed.(*Editor)
	if ed.CurView() != nil {
		workDir = ed.CurView().WorkDir()
	}
	v := ed.AddViewSmart(nil)
	v.highlighter = &TermHighlighter{}
	if interactive {
		v.SetViewType(core.ViewTypeInteractive)
	}
	b, err := backend.NewMemBackendCmd(args, workDir, v.Id(), nil, false)
	b.MaxRows = core.Ed.Config().MaxCmdBufferLines
	if err != nil {
		ed.SetStatusErr(err.Error())
	}
	v.backend = b
	return v.Id()
}
