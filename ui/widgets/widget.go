package widgets

type Widget interface {
	// Get the widget bounds
	Bounds() (y1, x1, y2, x2 int)
	// Render forec re=-rendering the view UI.
	Render()
	// Set the widget bounds
	SetBounds(y1, x1, y2, x2 int)
}
