// goed_api provides a command line wrapper around the Goed client API
// github.com/tcolar/goed/api/client
package main

import (
	"fmt"
	"os"

	"github.com/tcolar/goed/core"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	// TODO: command line version of the api
	app = kingpin.New("goed_api", "API service for goed Editor")

	apiVersion    = app.Command("api_version", "Returns the API version.")
	curView       = app.Command("cur_view", "Returns the index of the current view.")
	highlight     = app.Command("highlight", "Highlight a given file.")
	highlightFile = highlight.Arg("file", "file to highlight").Required().String()

	//loc = test.Arg("", "").Default(".").Required().String()
)

func main() {
	kingpin.Version(core.Version)
	action := kingpin.MustParse(app.Parse(os.Args[1:]))
	Dispatch(action)
}

func Dispatch(action string) {
	switch action {
	case apiVersion.FullCommand():
		ApiVersion()
	case curView.FullCommand():
		CurView()
	case highlight.FullCommand():
		Highlight(*highlightFile)
	}
}

func ApiVersion() {
	fmt.Printf(core.ApiVersion)
}

func CurView() {
	/*if core.Ed.CurView() == nil {
		fmt.Println("No active view !")
		os.Exit(1)
	}
	fmt.Printf("%d", core.Ed.CurView().Id())*/
}

func Highlight(file string) {
}
