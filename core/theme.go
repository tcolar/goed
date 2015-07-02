package core

import (
	"path"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
)

const (
	Plain uint16 = iota + (1 << 8)
	Bold
	Underlined
)

// Theme represents a goed theme data.
type Theme struct {
	Bg       Style // default to term bg
	Fg       Style // default to term fg
	BgSelect Style // default to term bg
	FgSelect Style // default to term fg
	BgCursor Style
	FgCursor Style

	Comment                            Style
	String                             Style
	Number                             Style
	Keyword1, Keyword2, Keyword3       Style
	Symbol1, Symbol2, Symbol3          Style
	Separator1, Separator2, Separator3 Style

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

func ReadDefaultTheme() (*Theme, error) {
	return ReadTheme(path.Join(Home, "standard", "themes", "default.toml"))
}

func ReadTheme(loc string) (*Theme, error) {
	var theme Theme
	if _, err := toml.DecodeFile(loc, &theme); err != nil {
		return ReadDefaultTheme()
	}
	return &theme, nil
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
	if err != nil {
		return err
	}
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
	return nil
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
