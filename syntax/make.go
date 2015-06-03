package syntax

// Makefile
var SyntaxMake = syntax{
	FileNames:  []string{"makefile"},
	Extensions: []string{".make", ".mk"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`@#`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`#`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		"include", "export", "define", "endef", "undefine", "ifdef", "ifndef",
		"ifeq", "ifneq", "undefine", "override", "unexport", "private", "vpath",
		"endif",
	},
	Keywords2: []string{
		"MAKEFILES", "MAKE", "VPATH", "SHELL", "MAKESHELL", "MAKE_VERSION",
		"MAKE_HOST", "MAKELEVEL", "MAKEFLAGS", "GNUMAKEFLAGS", "MAKECMDGOALS",
		"CURDIR", "SUFFIXES", ".LIBPATTERNS",
		".PHONY", ".SUFFIXES", ".DEFAULT", ".PRECIOUS", ".INTERMEDIATE",
		".SECONDARY", ".SECONDEXPANSION", ".DELETE_ON_ERROR", ".IGNORE",
		".LOW_RESOLUTION_TIME", ".SILENT", ".EXPORT_ALL_VARIABLES",
		".NOTPARALLEL", ".ONESHELL", ".POSIX",
	},
	Keywords3: []string{
		"subst", "patsubst", "strip", "findstring", "filter", "filter-out",
		"sort", "word", "words", "wordlist", "firstword", "lastword", "dir",
		"notdir", "suffix", "basename", "addsuffix", "addprefix", "join",
		"wildcard", "realpath", "abspath", "error", "warning", "shell",
		"origin", "flavor", "foreach", "if", "or", "and", "call", "eval",
		"file", "value",
	},
	Symbols1: []string{
		"=", "?=", ":=", "::=", "+=", "!=",
	},
	Symbols2: []string{
		"$", "$$", "$@", "$$@", "$?", "$%", "$<", "$^", "$+", "$*",
		"$(@D)", "$(@F)", "$(%D)", "$(%F)", "$(<D)", "$(<D)",
		"$(^D)", "$(^D)", "$(+D)", "$(+D)", "$(?D)", "$(?D)",
	},
	Symbols3: []string{ // others
		"+", "-", "<", "^", "*",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", ";", `\`,
	},
	Separators3: []string{
		":", "::", "%:", "%.", "@", "%",
	},
}
