// Goed is a terminal based editor.
// https://github.com/tcolar/goed
package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime/debug"
	"time"

	"github.com/tcolar/goed/api"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	app    = kingpin.New("goed", "A code editor")
	test   = kingpin.Flag("testterm", "Prints colors to the terminal to check them.").Bool()
	colors = kingpin.Flag("c", "Number of colors(0,2,16,256). 0 means Detect.").Default("0").Int()
	config = kingpin.Flag("config", "Config file.").Default("config.toml").String()
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

	id := time.Now().UnixNano()

	core.Colors = *colors
	core.InitHome(id)
	core.ConfFile = *config
	core.Ed = ui.NewEditor()
	apiServer := api.Api{}

	defer func() {
		if fail := recover(); fail != nil {
			// writing panic to file because shell will be garbled
			fmt.Printf("Panicked with %v\n", fail)
			fmt.Printf("Writing panic to log %s \n", core.LogFile.Name())
			data := debug.Stack()
			log.Fatal(string(data))
		}
		core.Cleanup()
		// attempts to reset the terminal in case we left it in a bad state
		exec.Command("reset").Run()
	}()

	apiServer.Start()
	core.Ed.Start(*loc)
}
