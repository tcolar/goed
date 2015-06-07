package syntax

var SyntaxJava = syntax{
	Extensions: []string{".java"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`/*`, `*/`, ``, true, StyleComment), // ML comment
		NewSyntaxPattern(`//`, ``, ``, false, StyleComment),  // comment
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),  // string
		NewSyntaxPattern("'", "'", `\`, false, StyleString),  // char
	},
	Keywords1: []string{
		"const", "import", "class", "interface", "extends", "implements",
		"package", "emum",
	},
	Keywords2: []string{
		"abstract",
		"assert",
		"boolean",
		"break",
		"byte",
		"case",
		"catch",
		"char",
		"continue",
		"default",
		"do",
		"double",
		"else",
		"false",
		"final",
		"finally",
		"float",
		"for",
		"goto",
		"if",
		"instanceof",
		"int",
		"long",
		"native",
		"new",
		"null",
		"private",
		"protected",
		"public",
		"return",
		"short",
		"static",
		"strictfp",
		"super",
		"switch",
		"synchronized",
		"this",
		"throw",
		"throws",
		"transient",
		"true",
		"try",
		"void",
		"volatile",
		"while",
	},
	Symbols1: []string{ // ~ assignment
		"++", "+=", "-=", "*=", "/=", "%=",
		"--", "=", "&=", "^=", "|=", "<<=", ">>=", ">>>=",
	},
	Symbols2: []string{ // ~ comparators
		"&&", "||", ">=", "<=", "!=", "==", ">", "<",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>", ">>>", "~", "!",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ".", ";", ":",
	},
}
