package core

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
)

type Theme struct {
	Bg       Style // default to term bg
	Fg       Style // default to term fg
	BgSelect Style // default to term bg
	FgSelect Style // default to term fg
	BgCursor Style
	FgCursor Style
	Comment  Style
	String   Style
	Keyword  Style
	// Numbers Style
	// Custom1 Style
	// Symbols := = etc ...
	// brackets {[]}()
	FileClean        StyledRune
	FileDirty        StyledRune
	Scrollbar        StyledRune
	ScrollTab        StyledRune
	Statusbar        StyledRune
	StatusbarText    Style
	StatusbarTextErr Style
	Cmdbar           StyledRune
	CmdbarText       Style
	CmdbarTextOn     Style
	Viewbar          StyledRune
	ViewbarText      Style
	MoreTextSide     StyledRune
	MoreTextUp       StyledRune
	MoreTextDown     StyledRune
	TabChar          StyledRune
	Margin           StyledRune
	Close            StyledRune
}

func ReadTheme(path string) (*Theme, error) {
	theme := defaultTheme()
	if _, err := toml.DecodeFile(path, &theme); err != nil {
		return nil, err
	}
	return theme, nil
}

func defaultTheme() *Theme {
	fg := Style{0x0001}
	bg := Style{0x0000}
	return &Theme{
		Bg:               bg,
		Fg:               fg,
		BgSelect:         fg,
		FgSelect:         bg,
		BgCursor:         fg,
		FgCursor:         bg,
		Comment:          fg,
		String:           fg,
		Keyword:          fg,
		StatusbarText:    fg,
		StatusbarTextErr: Style{0x0101},
		CmdbarText:       fg,
		CmdbarTextOn:     Style{0x0101},
		ViewbarText:      fg,
		FileClean:        StyledRune{'✔', fg, bg},
		FileDirty:        StyledRune{'✗', fg, bg},
		Scrollbar:        StyledRune{'░', fg, bg},
		ScrollTab:        StyledRune{'▒', fg, bg},
		Statusbar:        StyledRune{'❊', fg, bg},
		Cmdbar:           StyledRune{'❊', fg, bg},
		Viewbar:          StyledRune{'–', fg, bg},
		MoreTextSide:     StyledRune{'…', fg, bg},
		MoreTextUp:       StyledRune{'⇡', fg, bg},
		MoreTextDown:     StyledRune{'⇣', fg, bg},
		TabChar:          StyledRune{'⇨', fg, bg},
		Margin:           StyledRune{'|', fg, bg},
		Close:            StyledRune{'✕', fg, bg},
	}
}

// The format of a style as stored in a file is 4 bytes, HexaDecimal as follows:
// Byte 1 : 256 color palette index
// Byte 2 : 16 color palette index
// Byte 3 : 2 color (monochrome) palette index
// Byte 4 : Attribute such as (0: plain, 1:Bold, 3 : Underlined)
type Style struct {
	uint16
}

func NewStyle(s uint16) Style {
	return Style{
		uint16: s,
	}
}

func (s Style) Uint16() uint16 {
	return s.uint16
}

func (s Style) WithAttr(attr uint16) Style {
	return Style{s.uint16 | attr}
}

func (s *Style) UnmarshalText(text []byte) error {
	parsed, err := strconv.ParseUint(string(text), 16, 32)
	var val = uint16(parsed & 0x3)
	val = val << 8
	switch Colors {
	case 256:
		val = val | uint16((parsed&0xFF000000)>>24)
	case 16:
		val = val | uint16((parsed&0x0F0000)>>16)
	default:
		val = val | uint16((parsed&0x0F00)>>8)
	}
	s.uint16 = val
	return err
}

type StyledRune struct {
	Rune rune
	Fg   Style
	Bg   Style
}

func (s *StyledRune) UnmarshalText(text []byte) error {
	str := string(text)
	parts := strings.Split(str, ",")
	s.Rune, _ = utf8.DecodeRune([]byte(parts[0]))
	st := Style{}
	st.UnmarshalText([]byte(parts[1]))
	s.Fg = st
	st.UnmarshalText([]byte(parts[2]))
	s.Bg = st
	return nil
}
