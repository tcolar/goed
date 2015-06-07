package syntax

var SyntaxJs = syntax{
	Extensions: []string{".js"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`/*`, `*/`, ``, true, StyleComment), // ML comment
		NewSyntaxPattern(`//`, ``, ``, false, StyleComment),  // comment
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),  // string
		NewSyntaxPattern("'", "'", `\`, false, StyleString),  // char
	},
	Keywords1: []string{
		"const", "import", "class", "interface", "extends", "implements",
		"package", "enum", "var",
	},
	Keywords2: []string{
		"abstract",
		"boolean",
		"break",
		"byte",
		"case",
		"catch",
		"char",
		"continue",
		"debugger",
		"default",
		"delete",
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
		"in",
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
		"super",
		"switch",
		"synchronized",
		"this",
		"throw",
		"throws",
		"transient",
		"true",
		"try",
		"typeof",
		"undefined",
		"void",
		"volatile",
		"while",
		"with",
	},
	Symbols1: []string{ // ~ assignment
		"++", "+=", "-=", "*=", "/=", "%=",
		"--", "=", "&=", "^=", "|=", "<<=", ">>=", ">>>=",
	},
	Symbols2: []string{ // ~ comparators
		"&&", "||", ">=", "<=", "!=", "==", ">", "<", "===", "!==", "?:",
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
