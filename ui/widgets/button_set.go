package widgets

import "github.com/tcolar/goed/core"

type ButtonSet struct {
	BaseWidget
	Buttons []*Button
	Default int
}

func NewButtonSet() *ButtonSet {
	bw := BaseWidget{}
	bw.Bg = core.NewStyle(0)
	bw.Fg = core.NewStyle(0x0F)
	bs := &ButtonSet{
		BaseWidget: bw,
	}
	return bs
}

func (w *ButtonSet) Render() {
	for _, b := range w.Buttons {
		b.Render()
	}
}

func (w *ButtonSet) AddButton(b *Button, isDefault bool) {
	b.SetParent(w)
	w.Buttons = append(w.Buttons, b)
	if isDefault {
		w.Default = len(w.Buttons) - 1
		b.Active = true
	}
	_, _, _, width := b.Bounds()
	_, _, _, x := w.Bounds()
	b.SetBounds(0, x+1, 1, x+1+width)
	w.SetBounds(0, 0, 1, x+width+2)
}
