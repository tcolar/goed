package event

import (
	"fmt"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

var queue chan EventState = make(chan EventState)

func Queue(es EventState) {
	fmt.Println(es.Type)
	if es.Type == Evt_None {
		es.parseType()
	}
	queue <- es
}

func Shutdown() {
	close(queue)
}

func Listen() {
	for es := range queue {
		if done := handleEvent(&es); done {
			return
		}
	}
}

func handleEvent(es *EventState) bool {
	et := es.Type

	curView := actions.Ar.EdCurView()
	if curView < 0 {
		return false
	}
	// TODO: only do those as needed ?
	ln, col := actions.Ar.ViewCurPos(curView)
	l1, c1, _, _ := actions.Ar.ViewBounds(curView)
	offy, offx := actions.Ar.ViewScrollPos(curView)

	mouseView := actions.Ar.EdViewAt(es.MouseY, es.MouseX)
	mouseCol := es.MouseX - c1 + offx - 2
	mouseLn := es.MouseY - l1 + offy - 2

	dirty := false
	es.inDrag = false

	actions.Ar.ViewAutoScroll(curView, 0, 0, false)

	// TODO: cmdbar, term(ctrl+c)
	// TODO : common/termonly//cmdbar/view only
	// TODO: couldn't cmdbar e a view ?
	//fmt.Printf("%s %s mv:%d [%d:%d] (%d,%d)\n", et, es.String(), mouseView, mouseLn, mouseCol, es.MouseY, es.MouseX)

	switch et {
	case EvtBackspace:
		actions.Ar.ViewBackspace(curView)
		dirty = true
	case EvtBottom:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtBottom)
	case EvtCloseWindow:
		actions.Ar.EdDelView(curView, true)
	case EvtCut:
		actions.Ar.ViewCut(curView)
		dirty = true
	case EvtCopy:
		actions.Ar.ViewCopy(curView)
		dirty = true
	case EvtDelete:
		actions.Ar.ViewDeleteCur(curView)
		dirty = true
	case EvtEnd:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtEnd)
	case EvtEnter:
		actions.Ar.ViewInsertNewLine(curView)
		dirty = true
	case EvtHome:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtHome)
	case EvtMoveDown:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtDown)
	case EvtMoveLeft:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtLeft)
	case EvtMoveRight:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtRight)
	case EvtMoveUp:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtUp)
	case EvtNavDown:
		actions.Ar.EdViewNavigate(core.CursorMvmtDown)
	case EvtNavLeft:
		actions.Ar.EdViewNavigate(core.CursorMvmtLeft)
	case EvtNavRight:
		actions.Ar.EdViewNavigate(core.CursorMvmtRight)
	case EvtNavUp:
		actions.Ar.EdViewNavigate(core.CursorMvmtUp)
	case EvtOpenInNewView:
		actions.Ar.ViewClearSelections(mouseView)
		actions.Ar.ViewMoveCursor(mouseView, es.MouseY-l1-2-ln, es.MouseX-c1-2-col)
		actions.Ar.ViewOpenSelection(mouseView, true)
	case EvtOpenInSameView:
		actions.Ar.ViewOpenSelection(curView, false)
	case EvtOpenTerm:
		v := actions.Ar.EdOpenTerm([]string{core.Terminal})
		actions.Ar.EdActivateView(v, 0, 0)
	case EvtPaste:
		actions.Ar.ViewPaste(curView)
		dirty = true
	case EvtPageDown:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtPgDown)
	case EvtPageUp:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtPgUp)
	case EvtQuit:
		if actions.Ar.EdQuitCheck() {
			return true
		}
	case EvtRedo:
		actions.Ar.ViewRedo(curView)
		dirty = true
	case EvtReload:
		actions.Ar.ViewReload(curView)
	case EvtSave:
		actions.Ar.ViewSave(curView)
	case EvtSelectAll:
		actions.Ar.ViewSelectAll(curView)
	// TODO other selects
	case EvtSetCursor:
		//	fmt.Printf("set cursor %d %d\n", mouseLn, mouseCol)
		actions.Ar.ViewClearSelections(mouseView)
		actions.Ar.EdActivateView(mouseView, mouseLn, mouseCol)
	case EvtTab:
		actions.Ar.ViewInsertCur(curView, "\t")
		dirty = true
	case EvtToggleCmdbar:
		actions.Ar.CmdbarToggle()
	case EvtTop:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtTop)
	case EvtUndo:
		actions.Ar.ViewUndo(curView)
		dirty = true
	//case EvtWinResize:
	//	actions.Ar.EdResize(ev.Height, ev.Width)
	case Evt_None:
		actions.Ar.ViewInsertCur(curView, es.Glyph)
		dirty = true
	default:
		actions.Ar.EdSetStatusErr("Unhandled action : " + string(et))
	}

	if dirty {
		actions.Ar.ViewSetDirty(curView, true)
	}

	actions.Ar.EdRender()

	return false
}
