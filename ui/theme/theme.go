package theme

import (
	"github.com/BurntSushi/toml"
	"github.com/tcolar/goed/ui/style"
)

// Theme represents a goed theme data.
type Theme struct {
	Bg       style.Style // default to term bg
	Fg       style.Style // default to term fg
	BgSelect style.Style // default to term bg
	FgSelect style.Style // default to term fg
	BgCursor style.Style
	FgCursor style.Style

	Comment                            style.Style
	String                             style.Style
	Number                             style.Style
	Keyword1, Keyword2, Keyword3       style.Style
	Symbol1, Symbol2, Symbol3          style.Style
	Separator1, Separator2, Separator3 style.Style

	FileClean        style.StyledRune
	FileDirty        style.StyledRune
	Scrollbar        style.StyledRune
	ScrollTab        style.StyledRune
	Statusbar        style.StyledRune
	StatusbarText    style.Style
	StatusbarTextErr style.Style
	Cmdbar           style.StyledRune
	CmdbarText       style.Style
	CmdbarTextOn     style.Style
	Viewbar          style.StyledRune
	ViewbarText      style.Style
	MoreTextSide     style.StyledRune
	MoreTextUp       style.StyledRune
	MoreTextDown     style.StyledRune
	TabChar          style.StyledRune
	Margin           style.StyledRune
	Close            style.StyledRune
}

func ReadTheme(loc string) (*Theme, error) {
	var theme Theme
	_, err := toml.DecodeFile(loc, &theme)
	return &theme, err
}
