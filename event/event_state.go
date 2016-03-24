package event

import (
	"fmt"
	"strings"
)

type EventState struct {
	Type EventType
	// current values
	Glyph          string
	Keys           []string
	Combo          Combo
	MouseBtns      map[int]bool
	MouseY, MouseX int

	// state
	//movingView               bool
	//lastLClickX, lastLClickY int
	//lastLClick               int64 // timestamp
	dragLn, dragCol int  // drag start point
	inDrag          bool // mouse dragging
}

func NewEventState() *EventState {
	return &EventState{
		MouseBtns: map[int]bool{},
		Keys:      []string{},
	}
}

func (e *EventState) hasMouse() bool {
	for _, on := range e.MouseBtns {
		if on {
			return true
		}
	}
	return false
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
	e.inDrag = false
}

func (e *EventState) KeyUp(key string) {
	e.updKey(key, false)
}

func (e *EventState) MouseUp(button, y, x int) {
	e.MouseBtns[button] = false
	e.inDrag = false
	// TODO some sort of endDrag event ?
}

func (e *EventState) MouseDown(button, y, x int) {
	if e.MouseBtns[button] && (e.MouseX != x || e.MouseY != y) {
		e.inDrag = true
	}
	e.MouseY, e.MouseX = y, x
	e.MouseBtns[button] = true
	if !e.inDrag {
		e.dragLn, e.dragCol = y, x
	}
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
outer:
	for _, k := range strings.Split(s, "+") {
		if k[0] == 'M' { // mouse
			a := "MC"
			if e.inDrag {
				a = "MD"
			}
			for btn, b := range e.MouseBtns {
				if b && fmt.Sprintf("%s%d", a, btn) == k {
					score++
					continue outer
				}
			}
			return 0
		}
		switch k { // kb
		case KeyFunction:
			if !e.Combo.Func {
				return 0
			}
		case "ctrl":
			if !e.Combo.RCtrl && !e.Combo.LCtrl {
				return 0
			}
		case KeyLeftControl:
			if !e.Combo.LCtrl {
				return 0
			}
		case KeyRightControl:
			if !e.Combo.RCtrl {
				return 0
			}
		case "alt":
			if !e.Combo.RAlt && !e.Combo.LAlt {
				return 0
			}
		case "lalt":
			if !e.Combo.LAlt {
				return 0
			}
		case "ralt":
			if !e.Combo.RAlt {
				return 0
			}
		case "super":
			if !e.Combo.RSuper && !e.Combo.LSuper {
				return 0
			}
		case "lsuper":
			if !e.Combo.LSuper {
				return 0
			}
		case "rsuper":
			if !e.Combo.RSuper {
				return 0
			}
		case "shift":
			if !e.Combo.RShift && !e.Combo.LShift {
				return 0
			}
		case "lshift":
			if !e.Combo.LShift {
				return 0
			}
		case "rshift":
			if !e.Combo.RShift {
				return 0
			}
		default:
			if !e.hasKey(k) {
				return 0
			}
		}
		score++
	}
	return score
}

func (e *EventState) String() string {
	s := []string{}
	for btn, b := range e.MouseBtns {
		if b {
			a := "MC"
			if e.inDrag {
				a = "MD"
			}
			s = append(s, fmt.Sprintf("%s%d", a, btn))
		}
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
