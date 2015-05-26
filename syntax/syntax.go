package syntax

import "sort"

var Syntaxes map[string]Syntax

func init() {
	Syntaxes = map[string]Syntax{}
	initSyntax(&syntaxGo)
	initSyntax(&syntaxMarkdown)
}

const (
	StyleNone StyleId = iota
	StyleComment
	StyleString
	StyleNumber
	StyleKw1
	StyleKw2
	StyleKw3
	StyleSymb1
	StyleSymb2
	StyleSymb3
	StyleSep1
	StyleSep2
	StyleSep3
)

type StyleId byte

type SyntaxItem struct {
	Text string
	Id   StyleId
}

func NewSyntaxItem(text string, id StyleId) SyntaxItem {
	return SyntaxItem{
		Text: text,
		Id:   id,
	}
}

// sorted SyntaxItems
type SyntaxItems []SyntaxItem

// Longest first, if equal length then alphabetically
// This is to optimize lexing
func (s SyntaxItems) Less(i, j int) bool {
	a, b := len(s[i].Text), len(s[j].Text)
	if a == b {
		return s[i].Text < s[j].Text
	}
	return a > b
}
func (s SyntaxItems) Len() int      { return len(s) }
func (s SyntaxItems) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func initSyntax(s *syntax) {
	kws := SyntaxItems{}
	for _, kw := range s.Keywords1 {
		kws = append(kws, NewSyntaxItem(kw, StyleKw1))
	}
	for _, kw := range s.Keywords2 {
		kws = append(kws, NewSyntaxItem(kw, StyleKw2))
	}
	for _, kw := range s.Keywords3 {
		kws = append(kws, NewSyntaxItem(kw, StyleKw3))
	}
	sort.Sort(kws)
	symbs := SyntaxItems{}
	for _, symb := range s.Symbols1 {
		symbs = append(symbs, NewSyntaxItem(symb, StyleSymb1))
	}
	for _, symb := range s.Symbols2 {
		symbs = append(symbs, NewSyntaxItem(symb, StyleSymb2))
	}
	for _, symb := range s.Symbols3 {
		symbs = append(symbs, NewSyntaxItem(symb, StyleSymb3))
	}
	for _, symb := range s.Separators1 {
		symbs = append(symbs, NewSyntaxItem(symb, StyleSep1))
	}
	for _, symb := range s.Separators2 {
		symbs = append(symbs, NewSyntaxItem(symb, StyleSep2))
	}
	for _, symb := range s.Separators3 {
		symbs = append(symbs, NewSyntaxItem(symb, StyleSep3))
	}
	sort.Sort(symbs)
	syntax := Syntax{
		Patterns: s.Patterns,
		Keywords: kws,
		Symbols:  symbs,
	}
	for _, ext := range s.Extensions {
		Syntaxes[ext] = syntax
	}
}

type Syntax struct {
	Patterns []SyntaxPattern
	Symbols  SyntaxItems
	Keywords SyntaxItems
}

type syntax struct {
	Extensions                            []string
	Patterns                              []SyntaxPattern
	Keywords1, Keywords2, Keywords3       []string
	Symbols1, Symbols2, Symbols3          []string
	Separators1, Separators2, Separators3 []string
	// TODO: Number patterns ??
}

type SyntaxPattern struct {
	Start         string
	End           string // if empty -> EOL
	Escape        string // if empty -> none
	MultiLine     bool   // Whether may span multiple lines
	MustStartLine bool   // Whether "start" must be the first thing on the line.
	StyleId       StyleId
}

func NewSyntaxPattern(start, end, escape string, ml bool, style StyleId) SyntaxPattern {
	if ml && len(end) == 0 {
		panic("Invalid syntax pattern")
	}
	return SyntaxPattern{
		Start:     start,
		End:       end,
		Escape:    escape,
		MultiLine: ml,
		StyleId:   style,
	}
}

func (p SyntaxPattern) WithMSL() SyntaxPattern {
	p.MustStartLine = true
	return p
}
