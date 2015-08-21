package main

import (
	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
)

func Search(pattern, path string, ignorecase bool) {
	ed := core.Ed.(*Editor)
	workDir := "."
	if ed.curView != nil {
		workDir = ed.CurView().WorkDir()
	}
	v := ed.AddViewSmart()
	b, err := backend.NewMemBackendCmd(args, workDir, v.Id(), nil)
	if err != nil {
		ed.SetStatusErr(err.Error())
	}
	v.backend = b
}
