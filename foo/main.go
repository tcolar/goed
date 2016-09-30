package main

import (
	"github.com/andlabs/ui"
	"github.com/kr/pretty"
)

func main() {
	err := ui.Main(func() {
		area := ui.NewArea(&AreaHandler{})
		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetChild(area)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

var _ ui.AreaHandler = (*AreaHandler)(nil)

type AreaHandler struct {
}

// Draw is sent when a part of the Area needs to be drawn.
// dp will contain a drawing context to draw on, the rectangle
// that needs to be drawn in, and (for a non-scrolling area) the
// size of the area. The rectangle that needs to be drawn will
// have been cleared by the system prior to drawing, so you are
// always working on a clean slate.
//
// If you call Save on the drawing context, you must call Release
// before returning from Draw, and the number of calls to Save
// and Release must match. Failure to do so results in undefined
// behavior.
func (h *AreaHandler) Draw(a *ui.Area, dp *ui.AreaDrawParams) {
	pretty.Println("draw")
	pretty.Println(a)
	pretty.Println(dp)
	fd := &ui.FontDescriptor{
		Family:  "Arial",
		Size:    14.0, // as a text size, for instance 12 for a 12-point font
		Weight:  ui.TextWeightNormal,
		Italic:  ui.TextItalicNormal,
		Stretch: ui.TextStretchNormal,
	}
	font := ui.LoadClosestFont(fd)
	tl := ui.NewTextLayout("foo", font, 200.0)
	dp.Context.Text(20.0, 10.0, tl)
}

// MouseEvent is called when the mouse moves over the Area
// or when a mouse button is pressed or released. See
// AreaMouseEvent for more details.
//
// If a mouse button is being held, MouseEvents will continue to
// be generated, even if the mouse is not within the area. On
// some systems, the system can interrupt this behavior;
// see DragBroken.
func (h *AreaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	//	pretty.Println("me")
	//	pretty.Println(a)
	//	pretty.Println(me)
}

// MouseCrossed is called when the mouse either enters or
// leaves the Area. It is called even if the mouse buttons are being
// held (see MouseEvent above). If the mouse has entered the
// Area, left is false; if it has left the Area, left is true.
//
// If, when the Area is first shown, the mouse is already inside
// the Area, MouseCrossed will be called with left=false.
// TODO what about future shows?
func (h *AreaHandler) MouseCrossed(a *ui.Area, left bool) {
	//	pretty.Println("mc")
	//	pretty.Println(a)
}

// DragBroken is called if a mouse drag is interrupted by the
// system. As noted above, when a mouse button is held,
// MouseEvent will continue to be called, even if the mouse is
// outside the Area. On some systems, this behavior can be
// stopped by the system itself for a variety of reasons. This
// method is provided to allow your program to cope with the
// loss of the mouse in this case. You should cope by cancelling
// whatever drag-related operation you were doing.
//
// Note that this is only generated on some systems under
// specific conditions. Do not implement behavior that only
// takes effect when DragBroken is called.
func (h *AreaHandler) DragBroken(a *ui.Area) {
	//	pretty.Println("db")
	//	pretty.Println(a)
}

// KeyEvent is called when a key is pressed while the Area has
// keyboard focus (if the Area has been tabbed into or if the
// mouse has been clicked on it). See AreaKeyEvent for specifics.
//
// Because some keyboard events are handled by the system
// (for instance, menu accelerators and global hotkeys), you
// must return whether you handled the key event; return true
// if you did or false if you did not. If you wish to ignore the
// keyboard outright, the correct implementation of KeyEvent is
// 	func (h *MyHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
// 		return false
// 	}
// DO NOT RETURN TRUE UNCONDITIONALLY FROM THIS
// METHOD. BAD THINGS WILL HAPPEN IF YOU DO.
func (h *AreaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	pretty.Println("ke")
	pretty.Println(a)
	pretty.Println(ke)
	return false
}
