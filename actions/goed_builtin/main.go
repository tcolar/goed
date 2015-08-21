package main

import (
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/tcolar/goed/core"
)

var (
	app = kingpin.New("goed_builtin", "API service for goed Editor")

	_goto            = app.Command(":", "Got to / Open a location.")
	_gotoLoc         = _goto.Arg("location", "[path]:[line]:[col]").Required().String()
	search           = app.Command("s", "Search text.")
	searchPattern    = search.Arg("pattern", "Grep pattern").Required().String()
	searchPath       = search.Arg("path", "Path to search, current file if ommited.").Default(".").String()
	searchIgnoreCase = search.Flag("i", "IgnoreCase").Default("false").Bool()
	searchReplace    = search.Flag("r", "Replace found text").Default("").String()
	viDelete         = app.Command("d", "Delete line(s), VI style")
	viDeleteCount    = viDelete.Arg("number of lines", "").Default("1").Int()
	viPaste          = app.Command("p", "Paste line(s), VI style")
	viYank           = app.Command("p", "Yank line(s), VI style")
	viYankCount      = viYank.Arg("number of lines", "").Default("1").Int()
)

func main() {
	kingpin.Version(core.Version)
	action := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch action {
	case "s":
		Search(*searchPattern, *searchPath, *searchIgnoreCase)
	default:
		os.Exit(9)
	}
}
