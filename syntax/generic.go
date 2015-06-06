package syntax

var SyntaxGeneric = syntax{
	Extensions: []string{"_"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Symbols1: []string{
		"=", ":",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		".", ",", ";",
	},
}
