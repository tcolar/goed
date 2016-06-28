package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

func HandleArgs(args []string) {
	if len(args) <= 0 || args[0] == "help" {
		actions.RegisterActions()
		fmt.Println("help : usage info")
		fmt.Println("edit <instid> <dir> <file>: Open a file and wait until closed.")
		fmt.Println("instances : get all goed instances Ids")
		fmt.Println("instance : get most recent goed instance Id")
		fmt.Println("open <instid> <dir> <file>: Open a file.")
		fmt.Println("version : get goed_api version")
		fmt.Println()
		fmt.Println("Goed Api methods: https://godoc.org/github.com/tcolar/goed/api")
		fmt.Println(actions.Usage())
		os.Exit(1)
	}
	switch args[0] {
	case "version":
		fmt.Println(core.Version)
	case "instance":
		handleInstances(true)
	case "instances":
		handleInstances(false)
	case "edit":
		handleEdit(args[1:])
	case "open":
		handleOpen(args[1:])
	default:
		// Everything else is passed to a goed instance
		handleAction(args)
	}
}

func handleAction(args []string) {
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
	results, err := Action(instance, append(args[0:1], args[2:]...))
	if err != nil {
		fmt.Printf("RPC call failed: %s\n", err.Error())
		os.Exit(1)
	}
	if len(results) > 0 {
		fmt.Println(strings.Join(results, " "))
	}
}

func handleInstances(lastOnly bool) {
	ids := core.Instances()
	for _, id := range ids {
		fmt.Println(id)
		if lastOnly {
			break
		}
	}
}

func handleOpen(args []string) {
	if len(args) < 3 {
		fmt.Printf("Action open needs instance, path, file arguments\n")
		os.Exit(1)
	}
	instance, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Printf("InstanceId must be a number: %s\n", err.Error())
		os.Exit(1)
	}

	err = Open(instance, args[1], args[2])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func handleEdit(args []string) {
	if len(args) < 3 {
		fmt.Printf("Action edit needs instance, path, file arguments\n")
		os.Exit(1)
	}
	instance, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Printf("InstanceId must be a number: %s\n", err.Error())
		os.Exit(1)
	}

	err = Edit(instance, args[1], args[2])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
