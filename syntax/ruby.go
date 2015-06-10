package syntax

var SyntaxRuby = syntax{
	Extensions: []string{".rb"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`=begin`, `=end`, ``, true, StyleComment), // ML comment
		NewSyntaxPattern(`#!`, ``, ``, false, StyleKw3),
		NewSyntaxPattern(`#`, ``, ``, false, StyleComment),  // comment
		NewSyntaxPattern("'", "'", `\`, false, StyleString), // string
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString), // string
	},
	Keywords3: []string{
		"BEGIN", "END",
	},
	Keywords1: []string{
		"class", "def", "module", "begin", "end",
	},
	Keywords2: []string{
		"alias",
		"and",
		"break",
		"case",
		"do",
		"else",
		"elsif",
		"ensure",
		"false",
		"for",
		"if",
		"in",
		"next",
		"nil",
		"not",
		"or",
		"redo",
		"rescue",
		"retry",
		"return",
		"self",
		"super",
		"then",
		"true",
		"undef",
		"unless",
		"until",
		"when",
		"while",
		"yield",
		"__FILE__",
		"__LINE__",
	},
	Symbols1: []string{ // ~ assignment
		"=", "**=", "+=", "-=", "*=", "/=", "%=",
	},
	Symbols2: []string{ // ~ comparators
		"<=>", ">=", "<=", "!=", "==", ">", "<", "===", ".eql?", ".equal?",
		"||", "&&", "!", "defined?",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>", "~", "**",
		"..", "...",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":", "::",
	},
}
