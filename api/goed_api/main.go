// goed_api provides a command line wrapper around the Goed client API
// github.com/tcolar/goed/api/client
package main

import (
	"fmt"
	"os"

	"github.com/tcolar/goed/api/client"
	"github.com/tcolar/goed/core"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	app = kingpin.New("goed_api", "API service for goed Editor")

	instances   = app.Command("instances", "Returns the known goed instance ID'S, space separated, latest first.")
	instances1  = instances.Flag("1", "Returns only the most recent instance ID.").Default("false").Bool()
	apiVersion  = app.Command("api_version", "Returns the API version.")
	apiVersionI = apiVersion.Arg("ApiInstanceId", "InstanceId").Required().Int64()
	viewReload  = app.Command("view_reload", "Reload a view's buffer from the source.")
	viewReloadI = viewReload.Arg("InstanceId", "InstanceId").Required().Int64()
	viewReloadV = viewReload.Arg("ViewId", "ViewId").Required().Int64()
	viewSave    = app.Command("view_save", "Save the view buffer to the source.")
	viewSaveI   = viewSave.Arg("InstanceId", "InstanceId").Required().Int64()
	viewSaveV   = viewSave.Arg("ViewId", "ViewId").Required().Int64()
	viewSrcLoc  = app.Command("view_src_loc", "Get the view's source document path.")
	viewSrcLocI = viewSrcLoc.Arg("InstanceId", "InstanceId").Required().Int64()
	viewSrcLocV = viewSrcLoc.Arg("ViewId", "ViewId").Required().Int64()
	viewCwd     = app.Command("view_cwd", "Change view working directory.")
	viewCwdI    = viewCwd.Arg("InstanceId", "InstanceId").Required().Int64()
	viewCwdV    = viewCwd.Arg("ViewId", "ViewId").Required().Int64()
	viewCwdLoc  = viewCwd.Arg("dir", "dir").Required().String()
	open        = app.Command("open", "Open a file/directory in Goed (New view).")
	openI       = open.Arg("InstanceId", "InstanceId").Required().Int64()
	openCwd     = open.Arg("cwd", "cwd").Required().String()
	openLoc     = open.Arg("loc", "loc").Required().String()
	edit        = app.Command("edit", "Edit a file in Goed and wait until saved.")
	editI       = edit.Arg("InstanceId", "InstanceId").Required().Int64()
	editCwd     = edit.Arg("cwd", "cwd").Required().String()
	editLoc     = edit.Arg("loc", "loc").Required().String()
)

func main() {
	kingpin.Version(core.Version)
	action := kingpin.MustParse(app.Parse(os.Args[1:]))
	Dispatch(action)
}

func Dispatch(action string) {
	switch action {
	case instances.FullCommand():
		Instances()
	case apiVersion.FullCommand():
		ApiVersion()
	case viewReload.FullCommand():
		ViewReload()
	case viewSave.FullCommand():
		ViewSave()
	case viewSrcLoc.FullCommand():
		ViewSrcLoc()
	case viewCwd.FullCommand():
		ViewCwd()
	case open.FullCommand():
		Open()
	case edit.FullCommand():
		Edit()
	default:
		kingpin.Usage()
	}
}

func Instances() {
	ids := core.Instances()
	for _, id := range ids {
		fmt.Println(id)
		if *instances1 {
			break
		}
	}
}

func ApiVersion() {
	version, err := client.ApiVersion(*apiVersionI)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(version)
}

func ViewReload() {
	err := client.ViewReload(*viewReloadI, *viewReloadV)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func ViewSave() {
	err := client.ViewSave(*viewSaveI, *viewSaveV)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func ViewSrcLoc() {
	loc, err := client.ViewSrcLoc(*viewSrcLocI, *viewSrcLocV)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(loc)
}

func Open() {
	err := client.Open(*openI, *openCwd, *openLoc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func ViewCwd() {
	err := client.ViewCwd(*viewCwdI, *viewCwdV, *viewCwdLoc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func Edit() {
	err := client.Edit(*editI, *editCwd, *editLoc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
