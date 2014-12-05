package main

import (
	"fmt"
	"os"
	"strings"
)

// TODO: this is lame for now, provided incremental output,
// and/or interactive mode
func (v *View) Exec() {
	if v.Cmd == nil {
		Ed.SetStatusErr("Command missing !")
		return
	}
	file, err := os.Create(Ed.BufferFile(v))
	if err != nil {
		Ed.SetStatusErr(err.Error())
		return
	}
	defer file.Close()
	v.Cmd.Stdout = file
	v.Cmd.Stderr = file
	title := strings.Join(v.Cmd.Args, " ")
	v.title = fmt.Sprintf("[RUNNING] %s", title)
	Ed.Render()
	err = v.Cmd.Run()
	if err != nil {
		v.title = fmt.Sprintf("[FAILED] %s", title)
		Ed.SetStatusErr(err.Error())
		return
	}
	Ed.Open(file.Name(), v, "")
	v.title = fmt.Sprintf("[DONE] %s", title)
	Ed.Render()
}
