// Goed is a terminal based editor.
// https://github.com/tcolar/goed
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/api"
	"github.com/tcolar/goed/api/client"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	apiCall = kingpin.Flag("api", "API call").Default("false").Bool()
	gui     = kingpin.Flag("g", "Start in GUI mode..").Default("false").Bool()
	test    = kingpin.Flag("testterm", "Prints colors to the terminal to check them.").Bool()
	colors  = kingpin.Flag("c", "Number of colors(0,2,16,256). 0 means Detect.").Default("0").Int()
	config  = kingpin.Flag("config", "Config file.").Default("config.toml").String()
	cpuprof = kingpin.Flag("cpuprof", "Cpu profile").Default("false").Bool()
	memprof = kingpin.Flag("memprof", "Mem profile").Default("false").Bool()

	locs = kingpin.Arg("location", "location to open").Strings()
)

func main() {

	kingpin.Version(core.Version)

	kingpin.Parse()
	if *apiCall {
		client.HandleArgs(os.Args[2:])
		return
	}
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
	if *cpuprof == true {
		f, err := os.Create("prof.cprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *memprof == true {
		f, err := os.Create("prof.mprof")
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			time.Sleep(1 * time.Minute)
			pprof.WriteHeapProfile(f)
			f.Close()
			os.Exit(0)
		}()
	}

	id := time.Now().UnixNano()

	core.Colors = *colors
	core.Bus = actions.NewActionBus()
	core.InitHome(id)
	core.ConfFile = *config
	core.Ed = ui.NewEditor(*gui)
	apiServer := api.Api{}

	startupChecks()

	defer func() {
		if fail := recover(); fail != nil {
			// writing panic to file because shell may be garbled
			fmt.Printf("Panicked with %v\n", fail)
			fmt.Printf("Writing panic to log %s \n", core.LogFile.Name())
			data := debug.Stack()
			log.Fatal(string(data))
		}
		core.Cleanup()

		core.Bus.Shutdown()

		// attempts to reset the terminal in case we left it in a bad state
		exec.Command("reset").Run()
	}()

	actions.RegisterActions()
	apiServer.Start()

	core.Ed.Start(*locs)
}

func startupChecks() {
	out, err := exec.Command("goed", "--api", "version").CombinedOutput()
	if err != nil {
		fmt.Printf("Could not find/run goed --api : %s", out)
		os.Exit(1)
	}
	v := strings.Trim(string(out), "\n\t\r ")
	if v != core.Version {
		fmt.Printf("goed --api is not at the expected version. (got %s, want %s)",
			v, core.Version)
		os.Exit(1)
	}
}
