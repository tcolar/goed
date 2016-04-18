package client

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

var ftext []string

func init() {
	b, _ := ioutil.ReadFile("../../test_data/file1.txt")
	for _, s := range bytes.Split(b, []byte{'\n'}) {
		ftext = append(ftext, string(s))
	}
	if len(ftext[len(ftext)-1]) == 0 {
		ftext = ftext[:len(ftext)-1]
	}
}

func file1(t *testing.T) int64 {
	vid := actions.Ar.EdViewByLoc("../../test_data/file1.txt")
	assert.NotEqual(t, vid, int64(-1))
	return vid
}

func TestViewAddSelection(t *testing.T) {
	vid := file1(t)
	res, err := Action(id, []string{"view_add_selection", vidStr(vid), "1", "2", "3", "4"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	s := actions.Ar.ViewSelections(vid)
	assert.Equal(t, len(s), 1)
	assert.Equal(t, s[0], *core.NewSelection(1, 2, 3, 4))
	res, err = Action(id, []string{"view_add_selection", vidStr(vid), "6", "7", "4", "5"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	s = actions.Ar.ViewSelections(vid)
	assert.Equal(t, len(s), 2)
	assert.Equal(t, s[0], *core.NewSelection(1, 2, 3, 4))
	assert.Equal(t, s[1], *core.NewSelection(4, 5, 6, 7)) // Normalized
}

func TestViewClearSelection(t *testing.T) {
	vid := file1(t)
	actions.Ar.ViewAddSelection(vid, 1, 2, 3, 4)
	actions.Ar.ViewAddSelection(vid, 4, 5, 6, 7)
	assert.NotEqual(t, len(actions.Ar.ViewSelections(vid)), 0)
	res, err := Action(id, []string{"view_clear_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, len(actions.Ar.ViewSelections(vid)), 0)
}

func TestViewSelectAll(t *testing.T) {
	vid := file1(t)
	res, err := Action(id, []string{"view_select_all", vidStr(vid)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	s := actions.Ar.ViewSelections(vid)
	assert.Equal(t, len(s), 1)
	assert.Equal(t, s[0], *core.NewSelection(1, 1, 12, 36))
}

func TestViewSelections(t *testing.T) {
	vid := file1(t)
	actions.Ar.ViewClearSelections(vid)
	res, err := Action(id, []string{"view_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	actions.Ar.ViewAddSelection(vid, 1, 2, 3, 4)
	res, err = Action(id, []string{"view_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, "1 2 3 4", res[0])
	actions.Ar.ViewAddSelection(vid, 5, 6, 7, 8)
	res, err = Action(id, []string{"view_selections", vidStr(vid)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, "1 2 3 4", res[0])
	assert.Equal(t, "5 6 7 8", res[1]) // Normalized
}

func TestViewText(t *testing.T) {
	vid := file1(t)
	// "out of bounds" shoud return no text and not panic
	res, err := Action(id, []string{"view_text", vidStr(vid), "0", "0", "0", "0"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	res, err = Action(id, []string{"view_text", vidStr(vid), "100", "100", "200", "200"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	// "all" text
	res, err = Action(id, []string{"view_text", vidStr(vid), "1", "1", "-1", "-1"})
	assert.Nil(t, err)
	assert.Equal(t, res, ftext)
	// single char
	res, err = Action(id, []string{"view_text", vidStr(vid), "1", "1", "1", "1"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0], "1")
	// with tabs involved
	res, err = Action(id, []string{"view_text", vidStr(vid), "10", "3", "10", "4"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0], "ab")
	res, err = Action(id, []string{"view_text", vidStr(vid), "10", "3", "10", "-1"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0], "abc")
	// multiline selection
	res, err = Action(id, []string{"view_text", vidStr(vid), "10", "5", "11", "2"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], "c")
	assert.Equal(t, res[1], "aa")
	res, err = Action(id, []string{"view_text", vidStr(vid), "7", "3", "10", "4"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 4)
	assert.Equal(t, res[0], "ξδεφγηιςκλμνοπθρστυωωχψζ")
	assert.Equal(t, res[1], "ΑΒΞΔΕΦΓΗΙςΚΛΜΝΟΠΘΡΣΤΥΩΩΧΨΖ")
	assert.Equal(t, res[2], "")
	assert.Equal(t, res[3], "		ab")
	// "backward" selection
	res, err = Action(id, []string{"view_text", vidStr(vid), "1", "6", "1", "2"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0], "23456")
	res, err = Action(id, []string{"view_text", vidStr(vid), "4", "2", "3", "25"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], "yz")
	assert.Equal(t, res[1], "AB")
}

/*
view_auto_scroll(int64, int, int, bool)
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
