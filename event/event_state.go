package event

import (
	"fmt"
	"strings"
)

type EventState struct {
	Type EventType
	// current values
	Glyph                  string
	Keys                   []string
	Combo                  Combo
	LMouse, MMouse, RMouse bool
	MouseBtn               int
	MouseY, MouseX         int

	// state
	movingView                         bool
	lastClickX, lastClickY             int
	lastLClick, lastMClick, lastRClick int64 // timestamp
	dragLn, dragCol                    int
	inDrag                             bool
}

func (e *EventState) parseType() {
	for chord, et := range standard {
		if e.matches(chord) {
			e.Type = et
			return
		}
	}
	e.Type = Evt_None
}

func (e *EventState) KeyDown(key string) {
	e.updKey(key, true)
}

func (e *EventState) KeyUp(key string) {
	e.updKey(key, false)
}

func (e *EventState) MouseUp(button, y, x int) {
	e.MouseY, e.MouseX = y, x
	e.MouseBtn = 0
}

func (e *EventState) MouseDown(button, y, x int) {
	e.MouseY, e.MouseX = y, x
	e.MouseBtn = button
}

func (e *EventState) updKey(key string, isDown bool) {
	switch key {
	case KeyLeftSuper:
		e.Combo.LSuper = isDown
	case KeyRightSuper:
		e.Combo.RSuper = isDown
	case KeyLeftControl:
		e.Combo.LCtrl = isDown
	case KeyRightControl:
		e.Combo.RCtrl = isDown
	case KeyLeftAlt:
		e.Combo.LAlt = isDown
	case KeyRightAlt:
		e.Combo.RAlt = isDown
	case KeyLeftShift:
		e.Combo.LShift = isDown
	case KeyRightShift:
		e.Combo.RShift = isDown
	case KeyFunction:
		e.Combo.Func = isDown
	default:
		i := 0
		for _, k := range e.Keys {
			if k == key {
				break
			}
			i++
		}
		if isDown && i >= len(e.Keys) {
			e.Keys = append(e.Keys, key)
		}
		if !isDown && i < len(e.Keys) {
			e.Keys = append(e.Keys[:i], e.Keys[i+1:]...)
		}
	}
}

func (e *EventState) hasKey(key string) bool {
	for _, k := range e.Keys {
		if k == key {
			return true
		}
	}
	return false
}

// todo : validate any combo does not start with another combo

// ie: "super-rshift-alt-ctrl-function-x"
// ie: "mm"
func (e *EventState) matches(s string) bool {
	for _, k := range strings.Split(s, "+") {
		switch k {
		case "ml":
			if !e.LMouse {
				return false
			}
		case "mr":
			if !e.RMouse {
				return false
			}
		case "mm":
			if !e.MMouse {
				return false
			}

		case KeyFunction:
			if !e.Combo.Func {
				return false
			}

		case "ctrl":
			if !e.Combo.RCtrl && !e.Combo.LCtrl {
				return false
			}
		case KeyLeftControl:
			if !e.Combo.LCtrl {
				return false
			}
		case KeyRightControl:
			if !e.Combo.RCtrl {
				return false
			}

		case "alt":
			if !e.Combo.RAlt && !e.Combo.LAlt {
				return false
			}
		case "lalt":
			if !e.Combo.LAlt {
				return false
			}
		case "ralt":
			if !e.Combo.RAlt {
				return false
			}

		case "super":
			if !e.Combo.RSuper && !e.Combo.LSuper {
				return false
			}
		case "lsuper":
			if !e.Combo.LSuper {
				return false
			}
		case "rsuper":
			if !e.Combo.RSuper {
				return false
			}

		case "shift":
			if !e.Combo.RShift && !e.Combo.LShift {
				return false
			}
		case "lshift":
			if !e.Combo.LShift {
				return false
			}
		case "rshift":
			if !e.Combo.RShift {
				return false
			}

		default:
			if !e.hasKey(k) {
				return false
			}
		}
	}
	return true
}

func (e *EventState) String() string {
	s := []string{}
	if e.MouseBtn != 0 {
		s = append(s, fmt.Sprintf("M%d", e.MouseBtn))
	}
	if e.LMouse {
		s = append(s, "lm")
	}
	if e.MMouse {
		s = append(s, "mm")
	}
	if e.RMouse {
		s = append(s, "rm")
	}
	if e.Combo.LSuper {
		s = append(s, "lsuper")
	}
	if e.Combo.RSuper {
		s = append(s, "rsuper")
	}
	if e.Combo.LShift {
		s = append(s, "lshift")
	}
	if e.Combo.RShift {
		s = append(s, "rshift")
	}
	if e.Combo.LAlt {
		s = append(s, "lalt")
	}
	if e.Combo.RAlt {
		s = append(s, "ralt")
	}
	if e.Combo.LCtrl {
		s = append(s, "lctrl")
	}
	if e.Combo.RCtrl {
		s = append(s, "rctrl")
	}
	if e.Combo.Func {
		s = append(s, "func")
	}
	s = append(s, e.Keys...)
	return strings.Join(s, "+")
}

type Combo struct {
	LSuper, RSuper bool
	LShift, RShift bool
	LAlt, RAlt     bool
	LCtrl, RCtrl   bool
	Func           bool
}
