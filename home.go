package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"
	"strconv"
)

func (e *Editor) initHome() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Error : %s \n", err.Error())
		e.Home = "goed"
	} else if runtime.GOOS == "windows" { // meh
		e.Home = path.Join(usr.HomeDir, "goed")
	} else {
		e.Home = path.Join(usr.HomeDir, ".goed")
	}
	os.MkdirAll(e.Home, 0777)
	os.MkdirAll(path.Join(e.Home, "buffers"), 0777)
}

func (e *Editor) BufferFile(v *View) string {
	return path.Join(e.Home, "buffers", strconv.Itoa(v.Id))
}
