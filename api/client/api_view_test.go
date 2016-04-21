package client

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func (as *ApiSuite) TestViewAddSelection(t *C) {
	vid := as.openFile1(t)
	res, err := Action(as.id, []string{"view_add_selection", vidStr(vid), "1", "2", "3", "4"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	s := actions.Ar.ViewSelections(vid)
	assert.Eq(t, len(s), 1)
	assert.Eq(t, s[0], *core.NewSelection(1, 2, 3, 4))
	res, err = Action(as.id, []string{"view_add_selection", vidStr(vid), "6", "7", "4", "5"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	s = actions.Ar.ViewSelections(vid)
	assert.Eq(t, len(s), 2)
	assert.Eq(t, s[0], *core.NewSelection(1, 2, 3, 4))
	assert.Eq(t, s[1], *core.NewSelection(4, 5, 6, 7)) // Normalized
}

func (as *ApiSuite) TestViewAutoScroll(t *C) {
	vid := as.openFile1(t)
	actions.Ar.ViewInsert(vid, 1, 1, "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nxx", true)
	actions.Ar.ViewSetCursorPos(vid, 1, 1)
	ln, col := actions.Ar.ViewScrollPos(vid)
	assert.Eq(t, ln, 1)
	assert.Eq(t, col, 1)
	actions.Ar.ViewAddSelection(vid, 1, 1, -1, -1)
	res, err := Action(as.id, []string{"view_auto_scroll", vidStr(vid), "5", "5"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	time.Sleep(300 * time.Millisecond)
	ln, col = actions.Ar.ViewScrollPos(vid)
	assert.True(t, ln > 1) // scrolled down some
	res, err = Action(as.id, []string{"view_auto_scroll", vidStr(vid), "-10", "-10"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	time.Sleep(300 * time.Millisecond)
	ln, col = actions.Ar.ViewScrollPos(vid)
	assert.Eq(t, ln, 1) // scrolled back to top
	res, err = Action(as.id, []string{"view_auto_scroll", vidStr(vid), "0", "0"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
}

func (as *ApiSuite) TestViewBackspace(t *C) {
	vid := as.openFile1(t)
	actions.Ar.ViewSetCursorPos(vid, 1, 3)
	res, err := Action(as.id, []string{"view_backspace", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, actions.Ar.ViewText(vid, 1, 1, 1, -1)[0], "134567890")
	res, err = Action(as.id, []string{"view_backspace", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, actions.Ar.ViewText(vid, 1, 1, 1, -1)[0], "34567890")
	res, err = Action(as.id, []string{"view_backspace", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, actions.Ar.ViewText(vid, 1, 1, 1, -1)[0], "34567890")
	// nothing left to backspace (@ 1,1)
	res, err = Action(as.id, []string{"view_backspace", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	// backspace with line wrap
	actions.Ar.ViewSetCursorPos(vid, 4, 1)
	res, err = Action(as.id, []string{"view_backspace", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, actions.Ar.ViewText(vid, 3, 1, 3, -1)[0], "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// backspace selection
	actions.Ar.ViewAddSelection(vid, 7, 3, 9, 1)
	res, err = Action(as.id, []string{"view_backspace", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, actions.Ar.ViewText(vid, 7, 1, 7, -1)[0], "ΑΒ	abc")
}

func (as *ApiSuite) TestViewBounds(t *C) {
	views := actions.Ar.EdViews()
	assert.Eq(t, len(views), 1)
	res, err := Action(as.id, []string{"view_bounds", vidStr(views[0])})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 4)
	assert.Eq(t, strings.Join(res, " "), "2 1 24 50") // whole editor
	vid := as.openFile1(t)
	res, err = Action(as.id, []string{"view_bounds", vidStr(views[0])})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 4)
	assert.Eq(t, strings.Join(res, " "), "2 1 12 50") //top half
	res, err = Action(as.id, []string{"view_bounds", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 4)
	assert.Eq(t, strings.Join(res, " "), "13 1 24 50") //bottom half
}

func (as *ApiSuite) TestViewClearSelection(t *C) {
	vid := as.openFile1(t)
	actions.Ar.ViewAddSelection(vid, 1, 2, 3, 4)
	actions.Ar.ViewAddSelection(vid, 4, 5, 6, 7)
	assert.NotEq(t, len(actions.Ar.ViewSelections(vid)), 0)
	res, err := Action(as.id, []string{"view_clear_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, len(actions.Ar.ViewSelections(vid)), 0)
}

func (as *ApiSuite) TestViewCmdStop(t *C) {
	marker := "4224"
	vid := actions.Ar.EdOpenTerm([]string{"sleep", marker})
	// "sleep" command should be running a while
	out, err := exec.Command("ps", "-ax").CombinedOutput()
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(out), "sleep "+marker))
	// This should stop it
	res, err := Action(as.id, []string{"view_cmd_stop", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	// check it's gone
	out, err = exec.Command("ps", "-ax").CombinedOutput()
	assert.Nil(t, err)
	assert.False(t, strings.Contains(string(out), "sleep "+marker))
}

func (as *ApiSuite) TestViewCols(t *C) {
	views := actions.Ar.EdViews()
	assert.Eq(t, len(views), 1)
	res, err := Action(as.id, []string{"view_cols", vidStr(views[0])})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "46") //49 - 1 - 1(scrollar) -2(left/right padding)
	// add new view and move to new column
	vid := as.openFile1(t)
	l, c2, _, _ := actions.Ar.ViewBounds(vid)
	actions.Ar.EdViewMove(vid, l, c2, 2, c2+20) // force view to it's own column
	res, err = Action(as.id, []string{"view_cols", vidStr(views[0])})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "16") // 20-4
	res, err = Action(as.id, []string{"view_cols", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "26") // 30-4
	actions.Ar.EdDelView(vid, true)
	res, err = Action(as.id, []string{"view_cols", vidStr(views[0])})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "46")
}

func (as *ApiSuite) TestViewCopy(t *C) {
	vid := as.openFile1(t)
	// copy line
	res, err := Action(as.id, []string{"view_copy", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	core.Bus.Flush()
	cb, _ := core.ClipboardRead()
	assert.Eq(t, cb, "1234567890")
	// copy selection
	actions.Ar.ViewClearSelections(vid)
	actions.Ar.ViewAddSelection(vid, 10, 2, 11, 10)
	actions.Ar.ViewSetCursorPos(vid, 11, 10)
	res, err = Action(as.id, []string{"view_copy", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	core.Bus.Flush()
	cb, _ = core.ClipboardRead()
	assert.Eq(t, cb, "	abc\naaa aaa.go")
}

func (as *ApiSuite) TestViewCursorCoords(t *C) {
	vid := as.openFile1(t)
	res, err := Action(as.id, []string{"view_cursor_coords", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 2)
	assert.Eq(t, res[0], "1")
	assert.Eq(t, res[1], "1")
	actions.Ar.ViewMoveCursorRoll(vid, 3, 2)
	res, err = Action(as.id, []string{"view_cursor_coords", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 2)
	assert.Eq(t, res[0], "4")
	assert.Eq(t, res[1], "3")
}

func (as *ApiSuite) TestViewCursorMvmt(t *C) {
	vid := as.openFile1(t)
	actions.Ar.ViewSetCursorPos(vid, 3, 5)
	ln, col := actions.Ar.ViewCursorPos(vid)
	assert.Eq(t, ln, 3)
	assert.Eq(t, col, 5)
	as.checkMvmt(t, vid, core.CursorMvmtRight, 3, 6)
	as.checkMvmt(t, vid, core.CursorMvmtLeft, 3, 5)
	as.checkMvmt(t, vid, core.CursorMvmtDown, 4, 5)
	as.checkMvmt(t, vid, core.CursorMvmtUp, 3, 5)
	as.checkMvmt(t, vid, core.CursorMvmtHome, 3, 1)
	as.checkMvmt(t, vid, core.CursorMvmtEnd, 3, 27)
	actions.Ar.ViewSetCursorPos(vid, 1, 5) // view is 12 lines (= page size)
	core.Bus.Flush()
	as.checkMvmt(t, vid, core.CursorMvmtPgDown, 12, 5)
	actions.Ar.ViewSetCursorPos(vid, 13, 5)
	core.Bus.Flush()
	as.checkMvmt(t, vid, core.CursorMvmtPgUp, 1, 5)
	actions.Ar.ViewSetCursorPos(vid, 1, 7)
	core.Bus.Flush()
	as.checkMvmt(t, vid, core.CursorMvmtBottom, 12, 37)
	as.checkMvmt(t, vid, core.CursorMvmtTop, 1, 1)
}

func (as *ApiSuite) checkMvmt(t *C, vid int64, mvmt core.CursorMvmt, eln, ecol int) {
	res, err := Action(as.id, []string{"view_cursor_mvmt", vidStr(vid), fmt.Sprintf("%d", mvmt)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	ln, col := actions.Ar.ViewCursorPos(vid)
	assert.Eq(t, ln, eln)
	assert.Eq(t, col, ecol)
}

func (as *ApiSuite) TestViewSelectAll(t *C) {
	vid := as.openFile1(t)
	res, err := Action(as.id, []string{"view_select_all", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	s := actions.Ar.ViewSelections(vid)
	assert.Eq(t, len(s), 1)
	assert.Eq(t, s[0], *core.NewSelection(1, 1, 12, 36))
}

func (as *ApiSuite) TestViewSelections(t *C) {
	vid := as.openFile1(t)
	actions.Ar.ViewClearSelections(vid)
	res, err := Action(as.id, []string{"view_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	actions.Ar.ViewAddSelection(vid, 1, 2, 3, 4)
	res, err = Action(as.id, []string{"view_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, "1 2 3 4", res[0])
	actions.Ar.ViewAddSelection(vid, 5, 6, 7, 8)
	res, err = Action(as.id, []string{"view_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 2)
	assert.Eq(t, "1 2 3 4", res[0])
	assert.Eq(t, "5 6 7 8", res[1]) // Normalized
}

func (as *ApiSuite) TestViewText(t *C) {
	vid := as.openFile1(t)
	// "out of bounds" shoud return no text and not panic
	res, err := Action(as.id, []string{"view_text", vidStr(vid), "0", "0", "0", "0"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "100", "100", "200", "200"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	// "all" text
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "1", "1", "-1", "-1"})
	assert.Nil(t, err)
	assert.DeepEq(t, res, as.ftext)
	// single char
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "1", "1", "1", "1"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "1")
	// with tabs involved
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "10", "3", "10", "4"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "ab")
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "10", "3", "10", "-1"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "abc")
	// multiline selection
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "10", "5", "11", "2"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 2)
	assert.Eq(t, res[0], "c")
	assert.Eq(t, res[1], "aa")
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "7", "3", "10", "4"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 4)
	assert.Eq(t, res[0], "ξδεφγηιςκλμνοπθρστυωωχψζ")
	assert.Eq(t, res[1], "ΑΒΞΔΕΦΓΗΙςΚΛΜΝΟΠΘΡΣΤΥΩΩΧΨΖ")
	assert.Eq(t, res[2], "")
	assert.Eq(t, res[3], "		ab")
	// "backward" selection
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "1", "6", "1", "2"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, res[0], "23456")
	res, err = Action(as.id, []string{"view_text", vidStr(vid), "4", "2", "3", "25"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 2)
	assert.Eq(t, res[0], "yz")
	assert.Eq(t, res[1], "AB")
}

/*
view_cursor_mvmt(int64, core.CursorMvmt)
view_cursor_pos(int64) int, int
view_cut(int64)
view_delete(int64, int, int, int, int, bool)
view_delete_cur(int64)
view_insert(int64, int, int, string, bool)
view_insert_cur(int64, string)
view_insert_new_line(int64)
view_move_cursor(int64, int, int)
view_move_cursor_roll(int64, int, int)
view_open_selection(int64, bool)
view_paste(int64)
view_redo(int64)
view_reload(int64)
view_render(int64)
view_rows(int64) int
view_save(int64)
view_scroll_pos(int64) int, int
view_set_cursor_pos(int64, int, int)
view_set_dirty(int64, bool)
view_set_title(int64, string)
view_set_vt_cols(int64, int)
view_set_workdir(int64, string)
view_stretch_selection(int64, int, int)
view_src_loc(int64) string
view_sync_slice(int64)
view_text(int64, int, int,int, int) string
view_undo(int64)
*/

// view_lock ?? (to protect while editing) ? -> with timeout ?
// view_unlock

// TODO : review what api calls should be internal only
