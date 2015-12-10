package ui

import (
	"fmt"
	"os"
	"time"

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

func execTerm(args []string) int64 {
	vid := exec(args, true)
	v := core.Ed.ViewById(vid).(*View)
	b := v.backend.(*backend.BackendCmd)
	time.Sleep(500 * time.Millisecond)
	ext := ".sh"
	if os.Getenv("SHELL") == "rc" {
		ext = ".rc"
	}
	cmd := ". $HOME/.goed/default/actions/goed" +
		fmt.Sprintf("%s %d %d\n", ext, core.InstanceId, v.Id())
	b.SendBytes([]byte(cmd))
	return vid
}
