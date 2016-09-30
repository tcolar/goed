// Goed is a terminal based editor.
// https://github.com/tcolar/goed
package goed

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/api"
	"github.com/tcolar/goed/api/client"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	apiCall    = kingpin.Flag("api", "API call").Default("false").Bool()
	termColors = kingpin.Flag("term-colors", "Prints colors to the terminal to check them.").Bool()
	termEvents = kingpin.Flag("term-events", "Display received events in a view.").Bool()
	colors     = kingpin.Flag("c", "Number of colors(0,2,16,256). 0 means Detect.").Default("0").Int()
	config     = kingpin.Flag("config", "Config file.").Default("config.toml").String()
	cpuprof    = kingpin.Flag("cpuprof", "Cpu profile").Default("false").Bool()
	memprof    = kingpin.Flag("memprof", "Mem profile").Default("false").Bool()

	locs = kingpin.Arg("location", "location to open").Strings()
)

func Initialize() *core.Config {

	kingpin.Version(core.Version)

	kingpin.Parse()
	if *apiCall {
		client.HandleArgs(os.Args[2:])
		return nil
	}
	if *termColors {
		core.TermColors()
		return nil
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
	core.ShowEvents = *termEvents
	core.Bus = actions.NewActionBus()
	core.InitHome(id)
	core.ConfFile = *config

	startupChecks()

	// capture the error stream (hopefully empty) so it does not get dumped onto our terminal UI.
	logFile, _ := os.Create(path.Join(core.Home, "logs", "stderr.txt"))
	syscall.Dup2(int(logFile.Fd()), 2)

	return core.LoadConfig(core.ConfFile)
}

func Terminate(term core.Term) {
	if fail := recover(); fail != nil {
		// writing panic to file because shell may be garbled
		fmt.Printf("Panicked with %v\n", fail)
		fmt.Printf("Writing panic to log %s \n", core.LogFile.Name())
		data := debug.Stack()
		log.Fatal(string(data))
	}
	core.Cleanup()

	core.Bus.Shutdown()

	if term != nil {
		term.Close()
	}
}

func Start(term core.Term, config *core.Config) {
	defer Terminate(term)
	core.Ed = ui.NewEditor(term, config)
	actions.RegisterActions()
	apiServer := api.Api{}
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
