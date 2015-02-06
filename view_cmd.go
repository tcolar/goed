package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TODO: this is lame for now, provided incremental output,
// and/or interactive mode
func (v *View) Exec(workDir string) {
	workDir, _ = filepath.Abs(workDir)
	v.WorkDir = workDir
	if v.Cmd == nil {
		Ed.SetStatusErr("Command missing !")
		return
	}
	file, err := os.Create(Ed.BufferFile(v.Id))
	if err != nil {
		Ed.SetStatusErr(err.Error())
		return
	}
	defer file.Close()
	v.Cmd.Stdout = file
	v.Cmd.Stderr = file
	v.Cmd.Dir = workDir
	title := strings.Join(v.Cmd.Args, " ")
	v.title = fmt.Sprintf("[RUNNING] %s", title)
	Ed.Render()
	err = v.Cmd.Run()
	Ed.Open(file.Name(), v, "")
	v.WorkDir = workDir // open() would have modified this
	if err != nil {
		v.title = fmt.Sprintf("[FAILED] %s", title)
		Ed.SetStatusErr(err.Error())
	} else {
		v.title = fmt.Sprintf("[DONE] %s", title)
		Ed.SetStatus(workDir)
	}
	Ed.Render()
}
