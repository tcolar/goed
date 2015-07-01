// goed_api provides a command line wrapper around the Goed client API
// github.com/tcolar/goed/api/client
package main

import (
	"fmt"
	"os"

	"github.com/tcolar/goed/api/client"
	"github.com/tcolar/goed/core"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	// TODO: command line version of the api
	app = kingpin.New("goed_api", "API service for goed Editor")

	instances  = app.Command("instances", "Returns the known goed instance ID'S, space separated, latest first.")
	instances1 = instances.Flag("1", "Returns only the most recent instance ID.").Default("false").Bool()
	apiVersion = app.Command("api_version", "Returns the API version.")
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
	id := 0 // TODO get insatnceid as arg
	version, err := client.ApiVersion(id)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(version)
}
