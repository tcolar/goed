package syntax

var SyntaxToml = syntax{
	Extensions: []string{".toml"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`"""`, `"""`, ``, true, StyleString),
		NewSyntaxPattern(`'''`, `'''`, ``, true, StyleString),
		NewSyntaxPattern(`#`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, ``, false, StyleString),
	},
	Keywords1: []string{
		"true", "false",
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
	Separators3: []string{
		"[[", "]]",
	},
}
