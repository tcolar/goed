package syntax

// c#
var SyntaxCSharp = syntax{
	Extensions: []string{".cs"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`/*`, `*/`, ``, true, StyleComment),
		NewSyntaxPattern(`//`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, true, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		"abstract",
		"as",
		"base",
		"bool",
		"break",
		"byte",
		"case",
		"catch",
		"char",
		"checked",
		"const",
		"continue",
		"decimal",
		"default",
		"delegate",
		"do",
		"double",
		"else",
		"event",
		"explicit",
		"extern",
		"false",
		"finally",
		"fixed",
		"float",
		"for",
		"foreach",
		"goto",
		"if",
		"implicit",
		"in",
		"int",
		"internal",
		"is",
		"lock",
		"long",
		"namespace",
		"new",
		"null",
		"object",
		"operator",
		"out",
		"override",
		"params",
		"private",
		"protected",
		"public",
		"readonly",
		"ref",
		"return",
		"sbyte",
		"sealed",
		"short",
		"sizeof",
		"stackalloc",
		"static",
		"string",
		"struct",
		"switch",
		"this",
		"throw",
		"true",
		"try",
		"typeof",
		"uint",
		"ulong",
		"unchecked",
		"unsafe",
		"ushort",
		"using",
		"virtual",
		"volatile",
		"void",
		"while",
	},
	Keywords2: []string{
		"class", "enum", "interface", "struct", "using",
	},
	Keywords3: []string{
		"#if", "#else", "#elif", "#endif", "#define", "#undef", "#warning",
		"#error", "#line", "#region", "#endregion", "#pragma",
		"#pragma warning", "#pragma checksum",
	},
	Symbols1: []string{ // ~ assignment
		">>=", "<<=", "++", "+=", "-=", "*=", "/=", "%=",
		"|=", "&=", "^=", "--", "=",
	},
	Symbols2: []string{ // ~ comparators
		"&&", "||", ">=", "<=", "!=", "==", ">", "<", "!", "?", "?:",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "|", "&", "^", "<<", ">>",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}", "<", ">",
	},
	Separators2: []string{
		",", ".", ";", ":", "->", "->*", ".*", "::", "=>",
	},
}