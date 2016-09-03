package event

type MouseButton int

const (
	MouseLeft       = 1 << iota // 1
	MouseMiddle                 // 2
	MouseRight                  // 4
	MouseWheelUp                // 8
	MouseWheelDown              // 16
	MouseWheelLeft              // 32
	MouseWheelRight             // 64
)
