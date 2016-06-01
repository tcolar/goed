package event

import (
	"fmt"
	"log"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

var queue chan *Event = make(chan *Event)

func Queue(e *Event) {
	if e.Type == Evt_None {
		e.parseType()
	}
	queue <- e
}

func Shutdown() {
	close(queue)
}

func Listen() {
	es := &eventState{}
	for e := range queue {
		if done := handleEvent(e, es); done {
			return
		}
	}
}

func handleEvent(e *Event, es *eventState) bool {
	et := e.Type
	curView := actions.Ar.EdCurView()
	actions.Ar.ViewAutoScroll(curView, 0, 0)

	ln, col := actions.Ar.ViewCursorPos(curView)
	x, y := 0, 0 // relative mouse

	if e.hasMouse() {
		curView, y, x = actions.Ar.EdViewAt(e.MouseY+1, e.MouseX+1)
		ln, col = actions.Ar.ViewTextPos(curView, y, x)
		if e.inDrag && e.dragLn == -1 {
			e.dragLn, e.dragCol = ln, col
		}
	}

	if curView < 0 {
		return false
	}

	if et != Evt_None {
		fmt.Printf("%s %s y:%d x:%d ln:%d col:%d my:%d mx:%d - %v\n",
			et, e.String(), y, x, ln, col, e.MouseY, e.MouseX, e.inDrag)
	}

	vt := actions.Ar.ViewType(curView)
	if !e.hasMouse() && vt == core.ViewTypeShell {
		fmt.Println("te")
		handleTermEvent(curView, e)
		return false
	}

	dirty := false

	// TODO : common/termonly//cmdbar/view only
	// TODO: couldn't cmdbar be a view ?

	// TODO : dbl click
	// TODO : term Enter + VT100
	// TODO : cmdbar
	// TODO : shift selections
	// TODO : mouse select / scroll / drag / drag + scroll
	// TODO : down/pg_down selection seems buggy too (tabs ?)
	// TODO : window resize
	// TODO : allow other acme like events such as drag selection / click on selection

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
			actions.Ar.EdQuit()
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
		stretchSelection(curView, core.CursorMvmtDown)
		cs = false
	case EvtSelectEnd:
		stretchSelection(curView, core.CursorMvmtEnd)
		cs = false
	case EvtSelectHome:
		actions.Ar.ViewCursorMvmt(curView, core.CursorMvmtHome)
		stretchSelection(curView, core.CursorMvmtHome)
		cs = false
	case EvtSelectLeft:
		stretchSelection(curView, core.CursorMvmtLeft)
		cs = false
	case EvtSelectMouse:
		actions.Ar.ViewSetCursorPos(curView, ln, col)
		actions.Ar.ViewClearSelections(curView)
		fmt.Printf("%d %d %d %d\n", ln, col, e.dragLn, e.dragCol)
		actions.Ar.ViewAddSelection(curView, ln, col, e.dragLn, e.dragCol)
		cs = false
	case EvtSelectPageDown:
		stretchSelection(curView, core.CursorMvmtPgDown)
		cs = false
	case EvtSelectPageUp:
		stretchSelection(curView, core.CursorMvmtPgUp)
		cs = false
	case EvtSelectRight:
		stretchSelection(curView, core.CursorMvmtRight)
		cs = false
	case EvtSelectUp:
		stretchSelection(curView, core.CursorMvmtUp)
		cs = false
	case EvtSetCursor:
		dblClick := es.lastClickX == e.MouseX && es.lastClickY == e.MouseY &&
			time.Now().Unix()-es.lastClick <= 1
		y1, _, _, x2 := actions.Ar.ViewBounds(curView)
		// close button
		if e.MouseX+1 == x2-1 && e.MouseY+1 == y1 {
			actions.Ar.EdDelView(curView, true)
			break
		}
		// view "handle" (top left corner)
		if x == 1 && y == 1 {
			if dblClick {
				// view swap
				es.movingView = false
				cv := actions.Ar.EdCurView()
				actions.Ar.EdSwapViews(cv, curView)
				actions.Ar.EdActivateView(curView)
				break
			} // else, view move start
			es.movingView = true
			es.lastClickX = e.MouseX
			es.lastClickY = e.MouseY
			es.lastClick = time.Now().Unix()
			actions.Ar.EdSetStatusErr("Starting move, click new position or dbl click to swap")
			break
		}
		// Moving view to new position
		if es.movingView && (x == 1 || y == 1) {
			es.movingView = false
			actions.Ar.EdViewMove(es.lastClickY+1, es.lastClickX+1, e.MouseY+1, e.MouseX+1)
			break
		}
		// Set cursor position
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
		if len(e.Glyph) > 0 {
			// "plain" text
			actions.Ar.ViewInsertCur(curView, e.Glyph)
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

// Events for terminal/command views
func handleTermEvent(vid int64, e *Event) {
	es := false
	//ln, col := actions.Ar.ViewCursorCoords(vid)

	// Handle termbox special keys to VT100
	switch {
	/*	case termbox.KeyCtrlC:
			if len(actions.Ar.ViewSelections(vid)) > 0 {
				actions.Ar.ViewCopy(vid)
			} else {
				actions.Ar.TermSendBytes([]byte{byte(ev.Key)})
			}
		case termbox.KeyCtrlV:
			actions.Ar.ViewPaste(vid)*/
	// "special"/navigation keys
	case e.hasKey(KeyReturn):
		actions.Ar.TermSendBytes(vid, []byte{13})
	case e.hasKey(KeyDelete):
		actions.Ar.TermSendBytes(vid, []byte{127}) // delete (~ backspace)
	case e.hasKey(KeyUpArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'A'})
	case e.hasKey(KeyDownArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'B'})
	case e.hasKey(KeyRightArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'C'})
	case e.hasKey(KeyLeftArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'D'})
	case e.hasKey(KeyBackspace):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'C'}) // right
		actions.Ar.TermSendBytes(vid, []byte{127})          //delete
		// TODO: PgUp / pgDown not working right
	case e.hasKey(KeyNext):
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgDown)
		es = true
	case e.hasKey(KeyPrior):
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgUp)
		es = true
	case e.hasKey(KeyEnd):
		actions.Ar.TermSendBytes(vid, []byte{byte(0x05)}) // CTRL+E
		es = true
	case e.hasKey(KeyHome):
		actions.Ar.TermSendBytes(vid, []byte{byte(0x01)}) // CTRL+A
		es = true
		// function keys
	case e.hasKey(KeyF1):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'P'})
	case e.hasKey(KeyF2):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'Q'})
	case e.hasKey(KeyF3):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'R'})
	case e.hasKey(KeyF4):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'S'})
	case e.hasKey(KeyF5):
		actions.Ar.TermSendBytes(vid, []byte{27, '[', '1', '5', '~'})
	case e.hasKey(KeyF6):
		actions.Ar.TermSendBytes(vid, []byte{27, '[', '1', '7', '~'})
	case e.hasKey(KeyF7):
		actions.Ar.TermSendBytes(vid, []byte{27, '[', '1', '8', '~'})
	case e.hasKey(KeyF8):
		actions.Ar.TermSendBytes(vid, []byte{27, '[', '1', '9', '~'})
	case e.hasKey(KeyF9):
	case e.hasKey(KeyF10):
		actions.Ar.TermSendBytes(vid, []byte{27, '[', '2', '1', '~'})
	case e.hasKey(KeyF11):
		actions.Ar.TermSendBytes(vid, []byte{27, '[', '2', '3', '~'})
	case e.hasKey(KeyF12):
		actions.Ar.TermSendBytes(vid, []byte{27, '[', '2', '4', '~'})

	default:
		if len(e.Glyph) > 0 {
			actions.Ar.ViewInsertCur(vid, e.Glyph)
		} else {
			// 	TODO
			fmt.Printf("%#v\n", e)
		}
	}

	// extend keyboard selection
	if es { //&& ev.Meta == termbox.Shift {
		// TODO
		//		actions.Ar.ViewStretchSelection(vid, ln, col)
	} else {
		actions.Ar.ViewClearSelections(vid)
	}
}

func stretchSelection(vid int64, mvmt core.CursorMvmt) {
	l, c := actions.Ar.ViewCursorPos(vid)
	actions.Ar.ViewCursorMvmt(vid, mvmt)
	l2, c2 := actions.Ar.ViewCursorPos(vid)
	ss := actions.Ar.ViewSelections(vid)
	if len(ss) > 0 {
		if ss[0].LineTo == l && ss[0].ColTo == c {
			l = ss[0].LineFrom
			c = ss[0].ColFrom
		} else if ss[0].LineFrom == l && ss[0].ColFrom == c {
			l = ss[0].LineTo
			c = ss[0].ColTo
		}
	}
	actions.Ar.ViewClearSelections(vid)
	actions.Ar.ViewAddSelection(vid, l, c, l2, c2)
}
