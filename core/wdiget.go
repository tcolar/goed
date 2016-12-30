package core

type Widget interface {
	// Get the widget bounds (within parent)
	Bounds() (y1, x1, y2, x2 int)
	// Get parent widget or nil if none
	GetParent() Widget
	// Render forces re-rendering the view UI.
	Render()
	// Set the widget bounds (within parent)
	SetBounds(y1, x1, y2, x2 int)
	// Set the parent, typically internal use only
	SetParent(w Widget)
}
