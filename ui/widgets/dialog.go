package widgets

type Dialog struct {
	BaseWidget
	Label     *Label
	ButtonSet *ButtonSet
}

func NewDialog(text string, defaultButton int, buttons ...*Button) *Dialog {
	d := &Dialog{}
	d.Label = NewLabel(text)
	d.ButtonSet = NewButtonSet()
	d.Label.SetParent(d)
	d.ButtonSet.SetParent(d)
	for i, b := range buttons {
		d.ButtonSet.AddButton(b, i == defaultButton)
	}
	_, _, by2, _ := d.Label.Bounds()
	d.Label.Move(0, 1)
	d.ButtonSet.Move(by2+1, 0)

	// center the question ??
	// set overall bounds
	// draw frame ?
	// d.SetBounds(0, 0)

	return d
}

func (w *Dialog) Render() {
	w.Label.Render()
	w.ButtonSet.Render()
}
