package widgets

import (
	"strings"

	"github.com/tcolar/goed/core"
)

// A sigle button, to be added to a ButtonSet
type Button struct {
	BaseWidget
	Text          string
	ShortcutIndex int
	OnSelect      func() error
	ShortcutBg    core.Style
	ShortcutFg    core.Style
	Active        bool
}

func NewButton(text string, shortcut rune, onSelect func() error) *Button {
	bw := BaseWidget{}
	bw.SetBounds(0, 0, 2, len(text)+4)
	bw.Bg = core.NewStyle(0)
	bw.Fg = core.NewStyle(0x0F)
	shortcutIndex := -1
	if i := strings.Index(strings.ToLower(text), strings.ToLower(string(shortcut))); i >= 0 {
		shortcutIndex = i
	}
	return &Button{
		BaseWidget:    bw,
		Text:          text,
		ShortcutIndex: shortcutIndex,
		OnSelect:      onSelect,
		ShortcutBg:    core.NewStyle(0x0F),
		ShortcutFg:    core.NewStyle(0),
		Active:        false,
	}
}

func (w *Button) Render() {
	t := GetTermWidget(w)
	py1, px1, _, _ := w.GetParent().Bounds()
	y1, x1, _, _ := w.Bounds()
	y1 += py1
	x1 += px1
	t.Char(y1, x1, 0x250C, w.Fg, w.Bg)
	for i := 0; w.Active && i <= len(w.Text)+1; i++ {
		t.Char(y1, x1+i+1, 0x2500, w.Fg, w.Bg)
	}
	t.Char(y1, x1+len(w.Text)+3, 0x2510, w.Fg, w.Bg)
	if w.Active {
		t.Char(y1+1, x1, 0x2502, w.Fg, w.Bg)
	}
	for i := 0; i != len(w.Text); i++ {
		if i == w.ShortcutIndex {
			t.Char(y1+1, x1+i+2, rune(w.Text[i]), w.Fg.WithAttr(core.Bold), w.Bg)
		} else {
			t.Char(y1+1, x1+i+2, rune(w.Text[i]), w.Fg, w.Bg)
		}
	}
	if w.Active {
		t.Char(y1+1, x1+len(w.Text)+3, 0x2502, w.Fg, w.Bg)
	}
	t.Char(y1+2, x1, 0x2514, w.Fg, w.Bg)
	for i := 0; w.Active && i <= len(w.Text)+1; i++ {
		t.Char(y1+2, x1+i+1, 0x2500, w.Fg, w.Bg)
	}
	t.Char(y1+2, x1+len(w.Text)+3, 0x2518, w.Fg, w.Bg)
	t.Flush()
}
