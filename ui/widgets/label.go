package widgets

import "strings"

type Label struct {
	BaseWidget
	Lines []string
}

func NewLabel(text string) *Label {
	bw := BaseWidget{}
	lines := strings.Split(text, "\n")
	longest := 0
	for _, l := range lines {
		if len(l) > longest {
			longest = len(l)
		}
	}
	bw.SetBounds(0, 0, len(lines), longest)
	return &Label{
		BaseWidget: bw,
		Lines:      lines,
	}
}

func (w *Label) Render() {
	t := GetTermWidget(w)
	y1, x1, _, _ := w.Bounds()
	for i, line := range w.Lines {
		for j, r := range line {
			t.Char(y1+i, x1+j, r, w.Fg, w.Bg)
		}
	}
}
