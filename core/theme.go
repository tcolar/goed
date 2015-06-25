package core

import (
	"io/ioutil"
	"os"
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
	dir := path.Join(Home, "themes")
	loc := path.Join(dir, "default.toml")
	// If the theme does not exist yet(first start ?), create it
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
		err := ioutil.WriteFile(loc, []byte(defaultTheme), 0755)
		if err != nil {
			return nil, err
		}
	}
	return ReadTheme(loc)
}

func ReadTheme(loc string) (*Theme, error) {
	var theme Theme
	if _, err := toml.DecodeFile(loc, &theme); err != nil {
		return nil, err
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

var defaultTheme = `
Bg="EA000000"
Fg="DF030F01"
BgSelect="DF030F00"
FgSelect="EA000001"
BgCursor="21060F00"
FgCursor="EA000001"
Comment="F7070F00"
String="49060F00"
Keyword1="21030F01"
Keyword2="6F030F01"
Keyword3="69030F01"
Separator1="40050F00"
Separator2="95050F00"
Separator3="E2050F00"
Symbol1="CC040F01"
Symbol2="8D040F01"
Symbol3="D1040F01"
Statusbar = "❊,EB070000,EB000000"
StatusbarText = "BD660000"
StatusbarTextErr = "01C50001"
Cmdbar = "❊,EB070000,EB000000"
CmdbarText = "BD030000"
CmdbarTextOn = "BF660001"
Viewbar = "–,ED070000,ED000000"
ViewbarText = "E7020000"
FileClean = "✔,1C020000,E9000000"
FileDirty = "✗,A0010001,E9000000"
Scrollbar = "░,66060000,66000000"
ScrollTab = "▒,64060000,64000000"
MoreTextSide = "…,1F040000,EA000000"
MoreTextUp = "⇡,1F040000,EA000000"
MoreTextDown = "⇣,1F040000,EA000000"
TabChar = "⇨,EF000000,EA080000"
Margin = "|,EF000000,EA080000"
Close = "✕,33060000,EC080000"`
