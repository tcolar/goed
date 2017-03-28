package style

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

var Colors = 256

const (
	Plain uint16 = 1 << (8 + iota)
	Bold
	Underlined
)

// The format of a style as stored in a file is 4 bytes, HexaDecimal as follows:
// color, attr
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

func (s Style) Color() byte {
	return byte(s.uint16 & 0xFF)
}

func (s Style) IsBold() bool {
	return s.uint16&0xF00 == Bold
}

func (s Style) IsUnderlined() bool {
	return s.uint16&0xF00 == Underlined
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
