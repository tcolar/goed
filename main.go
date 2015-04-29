package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime/debug"

	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	app    = kingpin.New("goed", "A code editor")
	test   = kingpin.Flag("testterm", "Prints colors to the terminal to check them.").Bool()
	colors = kingpin.Flag("c", "Number of colors(0,2,16,256). 0 means Detect.").Default("0").Int()
	loc    = kingpin.Arg("location", "location to open").Default(".").String()
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	kingpin.Version(core.Version)

	kingpin.Parse()
	if *test {
		core.TestTerm()
		return
	}
	if *colors == 0 {
		*colors = core.DetectColors()
	}
	if *colors != 256 && *colors != 16 {
		*colors = 2
	}

	core.Colors = *colors
	core.InitHome()
	core.Ed = ui.NewEditor()

	defer func() {
		if fail := recover(); fail != nil {
			err := fail.(error)
			// writing panic to file because shell will be garbled
			fmt.Printf("Panicked with %s\n", err.Error())
			f := path.Join(core.Home, "panic.txt")
			fmt.Printf("Writing panic to %s \n", f)
			data := debug.Stack()
			ioutil.WriteFile(f, data, 0644)
		}
	}()

	core.Ed.Start(*loc)
}
