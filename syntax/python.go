package syntax

var SyntaxPython = syntax{
	Extensions: []string{".py", ".module"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern("'''", "'''", `\`, true, StyleString), // ML string
		NewSyntaxPattern(`"""`, `"""`, `\`, true, StyleString), // ML string
		NewSyntaxPattern(`#!`, ``, ``, false, StyleKw3),
		NewSyntaxPattern(`#`, ``, ``, false, StyleComment),  // ML comment
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString), // string
		NewSyntaxPattern("'", "'", `\`, false, StyleString), // string
	},
	Keywords1: []string{
		"None", "True", "False",
	},
	Keywords3: []string{
		"class", "def", "from", "global", "import",
	},
	Keywords2: []string{
		"and",
		"as",
		"assert",
		"break",
		"continue",
		"del",
		"elif",
		"else",
		"except",
		"exec",
		"finally",
		"for",
		"if",
		"in",
		"is",
		"lambda",
		"not",
		"or",
		"pass",
		"print",
		"raise",
		"return",
		"try",
		"while",
		"with",
		"yield",
	},
	Symbols1: []string{ // ~ assignment
		"=", "+=", "-=", "*=", "/=", "%=", "//=", "**=", "&=", "|=",
		"^=", ">>=", "<<=",
	},
	Symbols2: []string{ // ~ comparators
		"!=", "==", ">", "<", ">=", "<=",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "//", "**", "|", "&", "^", "<<", ">>", "~",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":", "->", "=>",
	},
}
