package syntax

var syntaxGo = syntax{
	Extensions: []string{".go"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern("`", "`", ``, true, StyleString),
		NewSyntaxPattern(`/*`, `*/`, ``, true, StyleComment),
		NewSyntaxPattern(`//`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		"const", "go", "import", "interface", "package", "struct", "type", "var",
	},
	Keywords2: []string{
		"break", "case", "chan", "continue", "default", "else",
		"fallthrough", "for", "goto", "if", "nmap", "range",
		"return", "select", "switch", "defer", "func",
	},
	Symbols1: []string{ // ~ assignment
		">>=", "<<=", "&^=", "++", "+=", "-=", "*=", "/=", "%=",
		"|=", "&=", "^=", "--", ":=", "=",
	},
	Symbols2: []string{ // ~ comparators
		">=", "<=", "&&", "||", ">=", "<=", "!=", "==", ">", "<", "!",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>", "&^",
		"...", "<-", "->",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":",
	},
}
