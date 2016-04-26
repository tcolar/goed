package event

import (
	"fmt"
	"log"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

var queue chan EventState = make(chan EventState)

func Queue(es EventState) {
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
	actions.Ar.ViewAutoScroll(curView, 0, 0)

	y, x := actions.Ar.ViewCursorCoords(curView)

	if es.hasMouse() {
		curView, y, x = actions.Ar.EdViewAt(es.MouseY+1, es.MouseX+1)
	}

	if curView < 0 {
		return false
	}

	//	col := actions.Ar.ViewLineRunesTo(curView, y, x)
	ln, col := actions.Ar.ViewTextPos(curView, y, x)

	dirty := false

	//	if et != Evt_None {
	fmt.Printf("%s %s y:%d x:%d ln:%d col:%d my:%d mx:%d - %v\n",
		et, es.String(), y, x, ln, col, es.MouseY, es.MouseX, es.inDrag)
	//	}

	// TODO : common/termonly//cmdbar/view only
	// TODO: couldn't cmdbar be a view ?

	// TODO : right click select/open still broken
	// TODO : dbl click
	// TODO : swap view
	// TODO : move view
	// TODO : click to close view
	// TODO : term Enter + VT100
	// TODO : cmdbar
	// TODO : shift selections
	// TODO : mouse select / scroll / drag / drag + scroll
	// TODO : down/pg_down slection seems buggy too (tabs ?)
	// TODO : window resize
	// TODO : allow other acme like events such as drag selection / click on selection
	// TODO : events & actions tests

	cs := true // clear selections

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
		actions.Ar.ViewSetCursorPos(curView, ln, col)
		actions.Ar.ViewOpenSelection(curView, true)
	case EvtOpenInSameView:
		actions.Ar.ViewSetCursorPos(curView, ln, col)
		actions.Ar.ViewOpenSelection(curView, false)
	case EvtOpenTerm:
		v := actions.Ar.EdOpenTerm([]string{core.Terminal})
		actions.Ar.EdActivateView(v)
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
		cs = false
	case EvtSelectDown:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtDown)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSelectEnd:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtEnd)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSelectHome:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtHome)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSelectLeft:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtLeft)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSelectMouse:
		actions.Ar.ViewSetCursorPos(curView, ln, col)
		actions.Ar.ViewStretchSelection(curView, es.dragLn, es.dragCol)
		cs = false
	case EvtSelectPageDown:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtPgDown)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSelectPageUp:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtPgUp)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSelectRight:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtRight)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSelectUp:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtUp)
		actions.Ar.ViewStretchSelection(curView, ln, col)
		cs = false
	case EvtSetCursor:
		if x == 1 && y == 1 {
			break
		}
		actions.Ar.ViewSetCursorPos(curView, ln, col)
		actions.Ar.EdActivateView(curView)
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
		if len(es.Glyph) > 0 {
			actions.Ar.ViewInsertCur(curView, es.Glyph)
			dirty = true
		}
	default:
		log.Println("Unhandled action : " + string(et))
		actions.Ar.EdSetStatusErr("Unhandled action : " + string(et))
	}

	if cs {
		actions.Ar.ViewClearSelections(curView)
	}

	if dirty {
		actions.Ar.ViewSetDirty(curView, true)
	}

	actions.Ar.EdRender()

	return false
}
