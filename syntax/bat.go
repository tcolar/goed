package syntax

// DOS batch files
var SyntaxBat = syntax{
	Extensions: []string{".bat", ".BAT"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`REM `, ``, ``, false, StyleComment),
		NewSyntaxPattern(`rem `, ``, ``, false, StyleComment),
	},
	Keywords1: []string{
		"NOT", "NUL", "null", "ECHO", "VAR", "IN", "DO", "GOTO", "PAUSE", "CHOICE",
		"EXIST", "CALL", "COMMAND", "SET", "SHIFT", "SIGN", "ERRORLEVEL",
		"CON", "PRN", "@ECHO", "IF", "ELSE", "END",
	},
	Keywords2: []string{
		"EQU", "NEQ", "LSS", "LEQ", "GTR", "GEQ",
	},
	Symbols1: []string{ // ~ assignment
		"=", "+=", "-=", "*=", "/=", "&=", "|=", "^=", ">>=", "<<=", "%%=",
	},
	Symbols2: []string{ // ~ comparators
		"==",
	},
	Symbols3: []string{ // others
		"+", "-", "*", "/", "%", "++", "--", "&", "|", "!", "~", "^", ">>", "<<",
	},
	Separators1: []string{
		"(", ")", "{", "}",
	},
	Separators2: []string{
		",", ".",
	},
}
