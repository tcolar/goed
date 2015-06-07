package syntax

var SyntaxPerl = syntax{
	Extensions: []string{".pl"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`=begin`, `=end`, ``, true, StyleComment), // ML comment
		NewSyntaxPattern("'", "'", `\`, true, StyleString),         // ML string
		NewSyntaxPattern(`#`, ``, ``, false, StyleComment),         // comment
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),        // string
		NewSyntaxPattern("`", "`", `\`, false, StyleString),        // string
	},
	Keywords1: []string{
		"BEGIN", "END",
	},
	Keywords2: []string{
		"continue",
		"else",
		"elsif",
		"for",
		"foreach",
		"goto",
		"if",
		"last",
		"next",
		"redo",
		"unless",
		"until",
		"while",
	},
	Keywords3: []string{
		"eq", "ne", "lt", "gt", "le", "ge", "cmp", "and", "or", "not", "xor",
	},
	Symbols1: []string{ // ~ assignment
		"**=", ".=", "+=", "-=", "*=", "/=", "%=",
	},
	Symbols2: []string{ // ~ comparators
		"<=>", "=~", "!~", "&&", "||", ">=", "<=", "!=", "==", ">", "<",
		"?:", "//", "~~",
	},
	Symbols3: []string{ // others
		"**", "..", "...",
		"+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>", "~", "!",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":", "->", "=>",
	},
}
