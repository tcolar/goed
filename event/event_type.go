package event

type EventType string

const (
	Evt_None          EventType = "_"
	EvtBackspace                = "backspace"
	EvtBottom                   = "bottom"
	EvtCloseWindow              = "close_window"
	EvtCut                      = "cut"
	EvtCopy                     = "copy"
	EvtDelete                   = "delete"
	EvtDeleteHome               = "delete_home"
	EvtEnd                      = "end"
	EvtHome                     = "home"
	EvtEnter                    = "enter"
	EvtMoveDown                 = "move_down"
	EvtMoveLeft                 = "move_left"
	EvtMoveRight                = "move_right"
	EvtMoveUp                   = "move_up"
	EvtNavDown                  = "nav_down"
	EvtNavLeft                  = "nav_left"
	EvtNavRight                 = "nav_right"
	EvtNavUp                    = "nav_up"
	EvtOpenInNewView            = "open_in_new_view"
	EvtOpenInSameView           = "open_in_same_view"
	EvtOpenTerm                 = "open_term"
	EvtPaste                    = "paste"
	EvtPageDown                 = "page_down"
	EvtPageUp                   = "page_up"
	EvtQuit                     = "quit"
	EvtRedo                     = "redo"
	EvtReload                   = "reload"
	EvtSave                     = "save"
	EvtScrollDown               = "scroll_down"
	EvtScrollUp                 = "scroll_up"
	EvtSelectMouse              = "select_mouse"
	EvtSelectAll                = "select_all"
	EvtSelectDown               = "select_down"
	EvtSelectEnd                = "select_end"
	EvtSelectHome               = "select_home"
	EvtSelectLeft               = "select_left"
	EvtSelectPageDown           = "select_page_down"
	EvtSelectPageUp             = "select_page_up"
	EvtSelectRight              = "select_right"
	EvtSelectUp                 = "select_up"
	EvtSelectWord               = "select_word"
	EvtSetCursor                = "set_cursor"
	EvtTab                      = "tab"
	EvtToggleCmdbar             = "toggle_cmd_bar"
	EvtTop                      = "top"
	EvtUndo                     = "undo"
	EvtWinResize                = "win_resize"
)

// Default bindings, if bindings.toml not found
// mirrors res/default/bindings.toml
var defaultBindings = map[string]EventType{
	// mouse
	"MC1":  "set_cursor",       // Mouse click left
	"MC4":  "open_in_new_view", // Mouse Click right
	"MC8":  "scroll_up",        // Mouse wheel up
	"MC16": "scroll_down",      // Mouse wheel down
	"MD1":  "select_mouse",     // Mouse Drag
	"MDC1": "select_word",      // Mouse Double Click left

	// special keys
	"escape":    "toggle_cmd_bar",
	"backspace": "backspace",
	"enter":     "enter",
	"return":    "enter",
	"tab":       "tab",
	"delete":    "delete",

	// control sequences
	"ctrl+a": "home",       // as in Acme
	"ctrl+b": "select_all", // made up since ctrl+a is used
	"ctrl+c": "copy",
	//"ctrl+d": "TODO", // Delete word before cursor (ctrl+w in Acme)
	"ctrl+e": "end",
	"ctrl+h": "move_left",  // vi like mvmt
	"ctrl+j": "move_right", // vi like mvmt
	"ctrl+k": "move_up",    // vi like mvmt
	"ctrl+l": "move_down",
	"ctrl+o": "open_in_same_view",
	"ctrl+n": "open_in_new_view", // or right click
	"ctrl+q": "quit",
	"ctrl+r": "reload",
	"ctrl+s": "save",
	"ctrl+t": "open_term",
	"ctrl+u": "delete_home",
	"ctrl+v": "paste",
	"ctrl+w": "close_window",
	"ctrl+x": "cut",
	"ctrl+y": "redo",
	"ctrl+z": "undo",

	// movement
	"right_arrow":       "move_right",
	"left_arrow":        "move_left",
	"up_arrow":          "move_up",
	"down_arrow":        "move_down",
	"prior":             "page_up",   // also works with Fn + up (OSX)
	"next":              "page_down", // also works with Fn + down (OSX)
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

	// navigation
	"alt+right_arrow":   "nav_right",
	"alt+left_arrow":    "nav_left",
	"alt+down_arrow":    "nav_down",
	"alt+up_arrow":      "nav_up",
	"super+right_arrow": "nav_right", // on OSX alt+arrow comes as super+arrow
	"super+left_arrow":  "nav_left",
	"super+down_arrow":  "nav_down",
	"super+up_arrow":    "nav_up",
}
