// goed_api provides a command line wrapper around the Goed client API
// github.com/tcolar/goed/api/client
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/api/client"
	"github.com/tcolar/goed/core"
)

/*
var (
	app = kingpin.New("goed_api", "API service for goed Editor")

	version = app.Command("version", "goed_api version.")

	instances   = app.Command("instances", "Returns the known goed instance ID'S, space separated, latest first.")
	instances1  = instances.Flag("1", "Returns only the most recent instance ID.").Default("false").Bool()
	apiVersion  = app.Command("api_version", "Returns the API version.")
	apiVersionI = apiVersion.Arg("ApiInstanceId", "InstanceId").Required().Int64()
	viewReload  = app.Command("view_reload", "Reload a view's buffer from the source.")
	viewReloadI = viewReload.Arg("InstanceId", "InstanceId").Required().Int64()
	viewReloadV = viewReload.Arg("ViewId", "ViewId").Required().Int64()
	viewCols    = app.Command("view_cols", "Size of the view in columns.")
	viewColsI   = viewCols.Arg("InstanceId", "InstanceId").Required().Int64()
	viewColsV   = viewCols.Arg("ViewId", "ViewId").Required().Int64()
	viewRows    = app.Command("view_rows", "Size of the views in rows.")
	viewRowsI   = viewRows.Arg("InstanceId", "InstanceId").Required().Int64()
	viewRowsV   = viewRows.Arg("ViewId", "ViewId").Required().Int64()
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
	viewVtCols  = app.Command("view_vt_cols", "Set vt100 wrap column.")
	viewVtColsI = viewVtCols.Arg("InstanceId", "InstanceId").Required().Int64()
	viewVtColsV = viewVtCols.Arg("ViewId", "ViewId").Required().Int64()
	viewVtColsC = viewVtCols.Arg("Cols", "Cols").Required().Int()
	open        = app.Command("open", "Open a file/directory in Goed (New view).")
	openI       = open.Arg("InstanceId", "InstanceId").Required().Int64()
	openCwd     = open.Arg("cwd", "cwd").Required().String()
	openLoc     = open.Arg("loc", "loc").Required().String()
	edit        = app.Command("edit", "Edit a file in Goed and wait until saved.")
	editI       = edit.Arg("InstanceId", "InstanceId").Required().Int64()
	editCwd     = edit.Arg("cwd", "cwd").Required().String()
	editLoc     = edit.Arg("loc", "loc").Required().String()
	generic     = app.Command("generic", "testing")
)*/

func main() {
	if len(os.Args) <= 1 || os.Args[1] == "--help" {
		//todo : usage
		fmt.Println("--help : usage info")
		fmt.Println("version : get goed_api version")
		fmt.Println("instances : get all goed insatnces Ids")
		fmt.Println("instance : get most recent goed instance Id")
		fmt.Println()
		fmt.Println("Goed Api methods: (More details at http://github.com/tcolar/goed/api/)")
		fmt.Println(actions.Usage())
		os.Exit(1)
	}
	switch os.Args[1] {
	case "version":
		fmt.Println(core.Version)
	case "instance":
		Instances(true)
	case "instances":
		Instances(false)
		/*	case apiVersion.FullCommand():
				ApiVersion()
			case viewReload.FullCommand():
				ViewReload()
			case viewRows.FullCommand():
				ViewRows()
			case viewCols.FullCommand():
				ViewCols()
			case viewSave.FullCommand():
				ViewSave()
			case viewSrcLoc.FullCommand():
				ViewSrcLoc()
			case viewCwd.FullCommand():
				ViewCwd()
			case viewVtCols.FullCommand():
				ViewVtCols()
			case open.FullCommand():
				Open()
			case edit.FullCommand():
				Edit()
			case generic.FullCommand():
				Generic()*/
	default:
		// Everything else is passed to a goed instance
		Action(os.Args[1:])
	}
}

func Action(args []string) {
	action := args[0]
	if len(args) < 2 {
		fmt.Printf("Action '%s' needs instanceId as first argument\n", action)
		os.Exit(1)
	}
	instance, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Printf("InstanceId must be a number: %s\n", err.Error())
		os.Exit(1)
	}
	results, err := client.Action(instance, append(args[0:1], args[2:]...))
	if err != nil {
		fmt.Printf("RPC call failed: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(strings.Join(results, " "))
}

func Instances(lastOnly bool) {
	ids := core.Instances()
	for _, id := range ids {
		fmt.Println(id)
		if lastOnly {
			break
		}
	}
}

/*
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

func ViewRows() {
	rows, err := client.ViewRows(*viewRowsI, *viewRowsV)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(rows)
}

func ViewCols() {
	cols, err := client.ViewCols(*viewColsI, *viewColsV)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(cols)
}

func ViewVtCols() {
	err := client.ViewVtCols(*viewVtColsI, *viewVtColsV, *viewVtColsC)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
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
}*/
