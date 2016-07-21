package core

// interafce for the command bar
type Commander interface {
	Backspace()
	Clear()
	CursorMvmt(mvmt CursorMvmt)
	Delete()
	Insert(text string)
	NewLine()
}
