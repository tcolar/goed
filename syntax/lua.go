package syntax

var SyntaxLua = syntax{
	Extensions: []string{".lua"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern("[[", "]]'", ``, true, StyleString), // ML string
		NewSyntaxPattern(`#!`, ``, ``, false, StyleKw3),
		NewSyntaxPattern(`--`, ``, ``, false, StyleComment), // comment
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString), // string
		NewSyntaxPattern("'", "'", `\`, false, StyleString), // string
	},
	Keywords1: []string{
		"function", "local", "::",
	},
	Keywords2: []string{
		"or", "and", "not",
	},
	Keywords3: []string{
		"break", "goto", "do", "while", "end", "repeat", "until", "if", "then", "elseif",
		"else", "for", "in", "return"},
	Symbols1: []string{ // ~ assignment
		"=", "+=", "-=", "*=", "/=", "%=", "//=", "**=", "&=", "|=",
		"^=", ">>=", "<<=",
	},
	Symbols2: []string{ // ~ comparators
		"~=", "==", ">", "<", ">=", "<=",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "//", "|", "&", "^", "<<", ">>", "~", "#", "..", "...",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":", "->",
	},
}
