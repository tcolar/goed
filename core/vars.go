// package core contains the core data structures and functionality
// leveraged y the other other Goed packages.
package core

import "os"

const ApiVersion = "v1"

var Trace = false
var ShowEvents = false

// Ed is thew editor singleton
var Ed Editable

// Colors is the number of colors to use in the terminal
var Colors int

// Home represent the goed "home" folder.
var Home string

// testing : whether we are in "unit test" mode.
var Testing bool

// ConfigFile holds the path to the config file currently in use.
var ConfFile string

// LogFile holds the path of the log file currently in use.
var LogFile *os.File

// terminal as defined by $SHELL
var Terminal string

var Bus ActionDispatcher

var ApiPort int

var Socket string // instance RPC socket

var InstanceId int64 // instance ID

type CursorMvmt byte

const (
	CursorMvmtRight      CursorMvmt = 0
	CursorMvmtLeft                  = 1
	CursorMvmtUp                    = 2
	CursorMvmtDown                  = 3
	CursorMvmtPgDown                = 4
	CursorMvmtPgUp                  = 5
	CursorMvmtHome                  = 6
	CursorMvmtEnd                   = 7
	CursorMvmtTop                   = 8
	CursorMvmtBottom                = 9
	CursorMvmtScrollDown            = 10
	CursorMvmtScrollUp              = 11
)

type ViewType int

const (
	ViewTypeStandard  ViewType = 0 // editable file
	ViewTypeShell              = 1 // interactive shell
	ViewTypeCmdOutput          = 2 // static command output
)
