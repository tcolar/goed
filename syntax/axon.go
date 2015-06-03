package syntax

var SyntaxAxon = syntax{
	Extensions: []string{".axon"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`/*`, `*/`, ``, true, StyleComment),
		NewSyntaxPattern(`//`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`**`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		"and",
		"catch",
		"do",
		"else",
		"end",
		"false",
		"if",
		"not",
		"null",
		"or",
		"return",
		"throw",
		"true",
		"try",
	},
	Symbols1: []string{ // ~ assignment
		"=",
	},
	Symbols2: []string{ // ~ comparators
		">=", "<=", "!=", "==", ">", "<", "<=>",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ":", "->", "=>",
	},
}
