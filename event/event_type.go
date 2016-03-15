package event

// TODO: temporary, read from config file
var standard = map[string]EventType{
	"ml":                "set_cursor",
	"mr":                "open_new_view",
	"esc":               "toggle_cmd_bar",
	"ctrl+q":            "quit",
	"backspace":         "backspace",
	"enter":             "enter",
	"ctrl+a":            "select_all",
	"ctrl+c":            "copy",
	"ctrl+r":            "reload",
	"ctrl+s":            "save",
	"ctrl+v":            "paste",
	"ctrl+w":            "close_window",
	"ctrl+x":            "cut",
	"ctrl+y":            "redo",
	"ctrl+z":            "undo",
	"right_arrow":       "move_right",
	"left_arrow":        "move_left",
	"up_arrow":          "move_up",
	"down_arrow":        "move_down",
	"prior":             "page_up",
	"next":              "page_down",
	"home":              "home",
	"end":               "end",
	"shift+right_arrow": "select_right",
	"shift+left_arrow":  "select_left",
	"shift+up_arrow":    "select_up",
	"shift+down_arrow":  "select_down",
	"shift+prior":       "select_page_up",
	"shift+next":        "select_page_down",
	"shift+home":        "select_home",
	"shift+end":         "select_end",
	"tab":               "tab",
	"delete":            "delete",
	"alt+o":             "open_new_view",
	"ctrl+o":            "open_same_view",
	"alt+right_arrow":   "nav_right",
	"alt+left_arrow":    "nav_left",
	"alt+down_arrow":    "nav_down",
	"alt+up_arrow":      "nav_up",
	"ctrl+t":            "open_term",
}

type EventType string

const (
	Evt_              EventType = "_"
	EvtBackspace                = "backspace"
	EvtCloseWindow              = "close_window"
	EvtCut                      = "cut"
	EvtCopy                     = "copy"
	EvtDelete                   = "delete"
	EvtEnd                      = "end"
	EvtHome                     = "home"
	EvtKeyEnter                 = "enter"
	EvtMoveDown                 = "move_down"
	EvtMoveLeft                 = "move_left"
	EvtMoveRight                = "move_right"
	EvtMoveUp                   = "move_up"
	EvtNavDown                  = "nav_down"
	EvtNavLeft                  = "nav_left"
	EvtNavRight                 = "nav_right"
	EvtNavUp                    = "nav_up"
	EvtOpenNewView              = "open_new_view"
	EvtOpenSameView             = "open_same_view"
	EvtOpenTerm                 = "open_term"
	EvtPaste                    = "paste"
	EvtPageDown                 = "page_down"
	EvtPageUp                   = "page_up"
	EvtQuit                     = "quit"
	EvtRedo                     = "redo"
	EvtReload                   = "reload"
	EvtSave                     = "save"
	EvtSelectAll                = "select_all"
	EvtSelectDown               = "select_down"
	EvtSelectEnd                = "select_end"
	EvtSelectHome               = "select_home"
	EvtSelectLeft               = "select_left"
	EvtSelectPageDown           = "select_page_down"
	EvtSelectPageUp             = "select_page_up"
	EvtSelectRight              = "select_right"
	EvtSelectUp                 = "select_up"
	EvtSetCursor                = "set_cursor"
	EvtStart                    = "start"
	EvtTab                      = "tab"
	EvtToggleCmdBar             = "toggle_cmd_bar"
	EvtUndo                     = "undo"
	EvtWinResize                = "win_resize" // no key/mouse
)
