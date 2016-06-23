package event

import (
	"fmt"
	"log"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

var queue = make(chan *Event, 200)

// Queue - Note: Queue a copy of the event
func Queue(e Event) {
	queue <- &e
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
	if e.Type == Evt_None {
		e.parseType()
	}
	et := e.Type
	curView := actions.Ar.EdCurView()
	actions.Ar.ViewAutoScroll(curView, 0, 0)

	ln, col := actions.Ar.ViewCursorPos(curView)
	x, y := 0, 0 // relative mouse

	if e.hasMouse() {
		curView, y, x = actions.Ar.EdViewAt(e.MouseY+1, e.MouseX+1)
		ln, col = actions.Ar.ViewTextPos(curView, y, x)
		if e.inDrag && e.dragLn == -1 {
			// start from prev click position as in case of terminal we are only going
			// to get a drag event *after* we have started drag motion.
			_, py, px := actions.Ar.EdViewAt(es.lastClickY+1, es.lastClickX+1)
			dl, dc := actions.Ar.ViewTextPos(curView, py, px)
			e.dragLn, e.dragCol = dl, dc
		}
	}

	if curView < 0 {
		return false
	}

	vt := actions.Ar.ViewType(curView)
	if !e.hasMouse() && vt == core.ViewTypeShell {
		handleTermEvent(curView, e)
		return false
	}

	dirty := false

	// TODO : cmdbar support -> couldn't cmdbar be a view ? -> redo ?

	// parity

	// TODO : sometimes moouse scroll puts gargabe character in terminal, seems to be termbox bug on OsX
	// TODO : down/pg_down selection not working -> seem to be termbox / Os MX issue
	// TODO : dbl click
	// TODO : allow other acme like events such as drag selection / click on selection

	cs := true // clear selections

	log.Printf("%#v", e)

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
		actions.Ar.ViewAddSelection(curView, ln, col, e.dragLn, e.dragCol)
		// Deal with selection autoscroll
		cols, rows := actions.Ar.ViewCols(curView), actions.Ar.ViewRows(curView)
		actions.Ar.EdSetStatus(fmt.Sprintf("%d %d %d %d", y, x, rows, cols))

		if y < 3 { // scroll up
			actions.Ar.ViewAutoScroll(curView, -5, 0)
		} else if y > rows+3 { // scroll down
			actions.Ar.ViewAutoScroll(curView, 5, 0)
		} else if x < 3 { //scroll left
			actions.Ar.ViewAutoScroll(curView, 0, -5)
		} else if x >= cols+2 { // scroll right
			actions.Ar.ViewAutoScroll(curView, 0, 5)
		}
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

			// Moving view to new position
		if es.movingView && (x == 1 || y == 1) {
			es.movingView = false
			actions.Ar.EdViewMove(es.lastClickY+1, es.lastClickX+1, e.MouseY+1, e.MouseX+1)
			break
		}

		y1, _, _, x2 := actions.Ar.ViewBounds(curView)
		es.lastClickX = e.MouseX
		es.lastClickY = e.MouseY
		es.lastClick = time.Now().Unix()

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
			actions.Ar.EdSetStatusErr("Starting move, click new position or dbl click to swap")
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
	case EvtWinResize:
		actions.Ar.ViewRender(curView)
	case Evt_None:
		if len(e.Glyph) > 0 {
			// "plain" text
			actions.Ar.ViewInsertCur(curView, e.Glyph)
			dirty = true
		} else {
			log.Println("Unhandled action : " + string(et))
			cs = false
		}
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
	cs := true
	ln, col := actions.Ar.ViewCursorCoords(vid)

	// Handle termbox special keys to VT100
	switch {
	case e.Type == EvtSelectMouse:
		actions.Ar.ViewSetCursorPos(vid, ln, col)
		actions.Ar.ViewClearSelections(vid)
		actions.Ar.ViewAddSelection(vid, ln, col, e.dragLn, e.dragCol)
		cs = false
	case e.Type == EvtCopy && len(actions.Ar.ViewSelections(vid)) > 0:
		// copy if copy event and there is a selection
		// if no selection, it may be Ctrl+C which is also used to terminate a command
		// (next case)
		actions.Ar.ViewCopy(vid)
		break
	case (e.Combo.LCtrl || e.Combo.RCtrl) && e.hasKey(KeyC): // CTRL+C
		actions.Ar.TermSendBytes(vid, []byte{byte(0x03)})
	case e.Type == EvtPaste:
		actions.Ar.ViewPaste(vid)
	// "special"/navigation keys
	case e.hasKey(KeyReturn):
		actions.Ar.TermSendBytes(vid, []byte{13})
	case e.hasKey(KeyTab):
		actions.Ar.TermSendBytes(vid, []byte{9})
	case e.hasKey(KeyDelete):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'C'})
		actions.Ar.TermSendBytes(vid, []byte{127}) // delete
	case e.hasKey(KeyUpArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'A'})
	case e.hasKey(KeyDownArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'B'})
	case e.hasKey(KeyRightArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'C'})
	case e.hasKey(KeyLeftArrow):
		actions.Ar.TermSendBytes(vid, []byte{27, 'O', 'D'})
	case e.hasKey(KeyBackspace):
		actions.Ar.TermSendBytes(vid, []byte{127})
		// TODO: PgUp / pgDown not working right
	case e.hasKey(KeyNext):
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgDown)
		cs = false
	case e.hasKey(KeyPrior):
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgUp)
		cs = false
	case e.hasKey(KeyEnd):
		actions.Ar.TermSendBytes(vid, []byte{byte(0x05)}) // CTRL+E
		cs = false
	case e.hasKey(KeyHome):
		actions.Ar.TermSendBytes(vid, []byte{byte(0x01)}) // CTRL+A
		cs = false
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
		if (e.Combo.LCtrl || e.Combo.RCtrl) && len(e.Keys) == 1 {
			// CTRL+? terminal combos
			val := byte([]rune(e.Keys[0])[0] - 'a' + 1)
			actions.Ar.TermSendBytes(vid, []byte{val})
		} else if len(e.Glyph) > 0 {
			actions.Ar.ViewInsertCur(vid, e.Glyph)
		} else {
			// TODO ??
			cs = false
		}
	}

	if cs { // clear selections
		actions.Ar.ViewClearSelections(vid)
	}
	actions.Ar.EdRender()

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
