package syntax

var SyntaxFantom = syntax{
	Extensions: []string{".fan", ".fwt", ".fog"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`"""`, `"""`, ``, true, StyleString), // triple quoted
		NewSyntaxPattern(`<|`, `|>`, ``, true, StyleString),   // DSL
		NewSyntaxPattern(`/*`, `*/`, ``, true, StyleComment),  // ML comment
		NewSyntaxPattern(`//`, ``, ``, false, StyleComment),   // comment
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),   // string
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),   // url
		NewSyntaxPattern("`", "`", `\`, false, StyleString),   // char
	},
	Keywords1: []string{
		"const", "using", "class", "interface", "mixin", "enum",
	},
	Keywords2: []string{
		"abstract",
		"as",
		"assert",
		"break",
		"case",
		"catch",
		"continue",
		"default",
		"do",
		"else",
		"false",
		"final",
		"finally",
		"for",
		"foreach",
		"if",
		"internal",
		"is",
		"isnot",
		"it",
		"native",
		"new",
		"null",
		"once",
		"override",
		"private",
		"protected",
		"public",
		"readonly",
		"return",
		"static",
		"super",
		"switch",
		"this",
		"throw",
		"true",
		"try",
		"virtual",
		"volatile",
		"void",
		"while",
	},
	Symbols1: []string{ // ~ assignment
		"++", "+=", "-=", "*=", "/=", "%=",
		"--", ":=", "=", "?:",
	},
	Symbols2: []string{ // ~ comparators
		"&&", "||", ">=", "<=", "!=", "==", ">", "<", "!", "===", "!==",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":", "->", "?.", "?->", "..", "..<",
	},
}
