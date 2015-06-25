package core

import "os"

const Version = "0.0.2"
const ApiVersion = "v1"

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

var ApiPort int
