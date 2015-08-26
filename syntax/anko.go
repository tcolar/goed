package syntax

var SyntaxAnko = syntax{
	Extensions: []string{".ank"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern("`", "`", ``, true, StyleString),
		NewSyntaxPattern(`#!`, ``, ``, false, StyleKw3),
		NewSyntaxPattern(`#`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		"module", "import", "var",
	},
	Keywords2: []string{
		"break", "continue", "throw", "if", "else", "swicth", "case",
		"try", "catch", "finaly", "default", "for", "range", "in", "new",
	},
	Keywords3: []string{
		"keys", "len", "println", "printf", "print", "true", "false", "nil",
		"range",
	},
	Symbols1: []string{ // ~ assignment
		"++", "+=", "-=", "*=", "/=", "|=", "--", "=",
	},
	Symbols2: []string{ // ~ comparators
		"&&", "||", ">=", "<=", "!=", "==", ">", "<", "!",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "|", "&", "^", "**", "...",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":",
	},
}
