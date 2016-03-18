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
	bestScore := 0
	t := Evt_None
	for chord, et := range standard {
		score := e.scoreMatch(chord)
		if score > bestScore {
			t = et
			bestScore = score
		}
	}
	e.Type = t
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

func (e *EventState) scoreMatch(s string) (score int) {
	for _, k := range strings.Split(s, "+") {
		switch k {
		case "ml":
			if !e.LMouse {
				return score
			}
		case "mr":
			if !e.RMouse {
				return score
			}
		case "mm":
			if !e.MMouse {
				return score
			}

		case KeyFunction:
			if !e.Combo.Func {
				return score
			}

		case "ctrl":
			if !e.Combo.RCtrl && !e.Combo.LCtrl {
				return score
			}
		case KeyLeftControl:
			if !e.Combo.LCtrl {
				return score
			}
		case KeyRightControl:
			if !e.Combo.RCtrl {
				return score
			}

		case "alt":
			if !e.Combo.RAlt && !e.Combo.LAlt {
				return score
			}
		case "lalt":
			if !e.Combo.LAlt {
				return score
			}
		case "ralt":
			if !e.Combo.RAlt {
				return score
			}

		case "super":
			if !e.Combo.RSuper && !e.Combo.LSuper {
				return score
			}
		case "lsuper":
			if !e.Combo.LSuper {
				return score
			}
		case "rsuper":
			if !e.Combo.RSuper {
				return score
			}

		case "shift":
			if !e.Combo.RShift && !e.Combo.LShift {
				return score
			}
		case "lshift":
			if !e.Combo.LShift {
				return score
			}
		case "rshift":
			if !e.Combo.RShift {
				return score
			}

		default:
			if !e.hasKey(k) {
				return score
			}
		}
		score++
	}
	return score
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
