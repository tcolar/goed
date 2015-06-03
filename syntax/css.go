package syntax

var SyntaxCss = syntax{
	Extensions: []string{".css"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`/*`, `*/`, ``, true, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		"true", "false", "block", "none", "auto",
	},
	Keywords2: []string{
		"@keyframes", "@media",
	},
	Symbols1: []string{ // ~ assignment
		":",
	},
	Symbols2: []string{
		"%", "px", "em",
	},
	Symbols3: []string{ // others
		"*", "#", "!",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ">",
	},
}
