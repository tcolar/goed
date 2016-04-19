package client

import (
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestViewAddSelection(c *C) {
	vid := s.openFile1(c)
	res, err := Action(s.id, []string{"view_add_selection", vidStr(vid), "1", "2", "3", "4"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	ss := actions.Ar.ViewSelections(vid)
	c.Assert(len(ss), Equals, 1)
	c.Assert(ss[0], Equals, *core.NewSelection(1, 2, 3, 4))
	res, err = Action(s.id, []string{"view_add_selection", vidStr(vid), "6", "7", "4", "5"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	ss = actions.Ar.ViewSelections(vid)
	c.Assert(len(ss), Equals, 2)
	c.Assert(ss[0], Equals, *core.NewSelection(1, 2, 3, 4))
	c.Assert(ss[1], Equals, *core.NewSelection(4, 5, 6, 7)) // Normalized
}

func (s *ApiSuite) TestViewClearSelection(c *C) {
	vid := s.openFile1(c)
	actions.Ar.ViewAddSelection(vid, 1, 2, 3, 4)
	actions.Ar.ViewAddSelection(vid, 4, 5, 6, 7)
	c.Assert(len(actions.Ar.ViewSelections(vid)), Not(Equals), 0)
	res, err := Action(s.id, []string{"view_clear_selections", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(len(actions.Ar.ViewSelections(vid)), Equals, 0)
}

func (s *ApiSuite) TestViewSelectAll(c *C) {
	vid := s.openFile1(c)
	res, err := Action(s.id, []string{"view_select_all", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	ss := actions.Ar.ViewSelections(vid)
	c.Assert(len(ss), Equals, 1)
	c.Assert(ss[0], Equals, *core.NewSelection(1, 1, 12, 36))
}

func (s *ApiSuite) TestViewSelections(c *C) {
	vid := s.openFile1(c)
	actions.Ar.ViewClearSelections(vid)
	res, err := Action(s.id, []string{"view_selections", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	actions.Ar.ViewAddSelection(vid, 1, 2, 3, 4)
	res, err = Action(s.id, []string{"view_selections", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert("1 2 3 4", Equals, res[0])
	actions.Ar.ViewAddSelection(vid, 5, 6, 7, 8)
	res, err = Action(s.id, []string{"view_selections", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 2)
	c.Assert("1 2 3 4", Equals, res[0])
	c.Assert("5 6 7 8", Equals, res[1]) // Normalized
}

func (s *ApiSuite) TestViewText(c *C) {
	vid := s.openFile1(c)
	// "out of bounds" shoud return no text and not panic
	res, err := Action(s.id, []string{"view_text", vidStr(vid), "0", "0", "0", "0"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "100", "100", "200", "200"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	// "all" text
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "1", "1", "-1", "-1"})
	c.Assert(err, IsNil)
	c.Assert(res, DeepEquals, s.ftext)
	// single char
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "1", "1", "1", "1"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert(res[0], Equals, "1")
	// with tabs involved
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "10", "3", "10", "4"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert(res[0], Equals, "ab")
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "10", "3", "10", "-1"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert(res[0], Equals, "abc")
	// multiline selection
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "10", "5", "11", "2"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 2)
	c.Assert(res[0], Equals, "c")
	c.Assert(res[1], Equals, "aa")
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "7", "3", "10", "4"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 4)
	c.Assert(res[0], Equals, "ξδεφγηιςκλμνοπθρστυωωχψζ")
	c.Assert(res[1], Equals, "ΑΒΞΔΕΦΓΗΙςΚΛΜΝΟΠΘΡΣΤΥΩΩΧΨΖ")
	c.Assert(res[2], Equals, "")
	c.Assert(res[3], Equals, "		ab")
	// "backward" selection
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "1", "6", "1", "2"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert(res[0], Equals, "23456")
	res, err = Action(s.id, []string{"view_text", vidStr(vid), "4", "2", "3", "25"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 2)
	c.Assert(res[0], Equals, "yz")
	c.Assert(res[1], Equals, "AB")
}

func (s *ApiSuite) TestViewBackspace(c *C) {
	vid := s.openFile1(c)
	actions.Ar.ViewSetCursorPos(vid, 1, 3)
	res, err := Action(s.id, []string{"view_backspace", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(actions.Ar.ViewText(vid, 1, 1, 1, -1)[0], Equals, "134567890")
	res, err = Action(s.id, []string{"view_backspace", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(actions.Ar.ViewText(vid, 1, 1, 1, -1)[0], Equals, "34567890")
	res, err = Action(s.id, []string{"view_backspace", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(actions.Ar.ViewText(vid, 1, 1, 1, -1)[0], Equals, "34567890")
	// nothing left to backspace (@ 1,1)
	res, err = Action(s.id, []string{"view_backspace", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	// backspace with line wrap
	actions.Ar.ViewSetCursorPos(vid, 4, 1)
	res, err = Action(s.id, []string{"view_backspace", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(actions.Ar.ViewText(vid, 3, 1, 3, -1)[0], Equals, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// backspace selection
	actions.Ar.ViewAddSelection(vid, 7, 3, 9, 1)
	res, err = Action(s.id, []string{"view_backspace", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(actions.Ar.ViewText(vid, 7, 1, 7, -1)[0], Equals, "ΑΒ	abc")

}

func (s *ApiSuite) TestViewCursorCoords(c *C) {
	vid := s.openFile1(c)
	res, err := Action(s.id, []string{"view_cursor_coords", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 2)
	c.Assert(res[0], Equals, "1")
	c.Assert(res[1], Equals, "1")
	actions.Ar.ViewMoveCursorRoll(vid, 3, 2)
	res, err = Action(s.id, []string{"view_cursor_coords", vidStr(vid)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 2)
	c.Assert(res[0], Equals, "4")
	c.Assert(res[1], Equals, "3")
}

/*
view_backspace(int64)
view_bounds(int64) int, int, int, int
view_cmd_stop(int64)
view_cols(int64) int
view_copy(int64)
view_cursor_coords(int64) int, int
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

// view_lock
// view_unlock ?? (to protect while editing) ?
