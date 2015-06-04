package syntax

var SyntaxShell = syntax{
	Extensions: []string{".sh", ".zsh", ".bash", ".ksh", ".rc"}, // TODO: rc should be separate
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`#!`, ``, ``, false, StyleKw3),
		NewSyntaxPattern(`#`, ``, ``, false, StyleComment),
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		// sh
		".", "break", ":", "cd", "continue", "eval", "exec", "exit", "export",
		"getopts", "hash", "pwd", "readonly", "return", "shift", "test", "times",
		"trap", "umask", "unset",
		// bash
		"alias", "bind", "builtin", "caller", "command", "declare", "echo",
		"enable", "help", "let", "local", "logout", "mapfile", "printf",
		"read", "readarray", "source", "type", "typeset", "ulimit", "unalias",
		// zsh
		"autoload", "bg", "bindkey", "bye", "cap", "chdir", "clone",
		"comparguments", "compcall", "compctl", "compdescribe", "compfiles",
		"compgroups", "compquote", "comptags", "comptry", "compvalues",
		"dirs", "disable", "disown", "echotc", "echoti", "emulate", "false",
		"fc", "fg", "float", "getcap", "getln", "history", "integer", "jos",
		"kill", "limit", "log", "noglob", "popd", "print", "pushd", "pushln",
		"r", "rehash", "sched", "set", "setcap", "setopt", "stat", "suspend",
		"true", "ttyctl-fu", "unfunction", "unhash", "unlimit", "unset",
		"unsetopt", "vared", "wait", "whence", "where", "which", "zcompile",
		"zformat", "zftp", "zle", "zmodload", "zparseopts", "zprof", "zpty",
		"zregexparse", "zsocket", "zstyle", "ztcp",
	},
	Keywords2: []string{
		// sh
		"CDPATH", "HOME", "IFS", "MAIL", "MAILPATH", "OPTARG", "OPTIND", "PATH",
		"PS1", "PS2",
		// bash
		"BASH", "BASHOPTS", "BASHPID", "BASH_ALIASES", "BASH_ARGC", "BASH_ARGV",
		"BASH_CMDS", "BASH_COMMAND", "BASH_COMPAT", "BASH_ENV",
		"BASH_EXECUTION_STRING", "BASH_LINENO", "BASH_REMATCH", "BASH_SOURCE",
		"BASH_SUBSHELL", "BASH_VERSINFO", "BASH_VERSION", "BASH_XTRACEFD",
		"CHILD_MAX", "COLUMNS", "COMP_CWORD", "COMP_LINE", "COMP_POINT", "COMP_TYPE",
		"COMP_KEY", "COMP_WORDBREAKS", "COMP_WORDS", "COMPREPLY", "COPROC",
		"DIRSTACK", "EMACS", "ENV", "EUID", "FCEDIT", "FIGNORE", "FUNCNAME",
		"FUNCNEST", "GLOBIGNORE", "GROUPS", "histchars", "HISTCMD", "HISTCONTROL",
		"HITSFILE", "HISTFILESIZE", "HISTIGNORE", "HISTSIZE", "HISTTIMEFORMAT",
		"HOSTFILE", "HOSTNAME", "HOSTTYPE", "IGNOREEOF", "INPUTRC", "LANG",
		"LC_ALL", "LC_COLLATE", "LC_CTYPE", "LC_MESSAGES", "LC_NUMERIC",
		"LINENO", "LINES", "MATCHTYPE", "MAILCHECK", "MAPFILE", "OLDPWD",
		"OPTERR", "OSTYPE", "PIPESTATUS", "POSIXLY_CORRECT", "PPID",
		"PROMPT_COMMAND", "PROMPT_DIRTRIM", "PS3", "PS4", "PWD", "RANDOM",
		"READLINE_LINE", "READLINE_POINT", "REPLY", "SECONDS", "SHELL",
		"SHELLOPTS", "SHLVL", "TIMEFORMAT", "TMOUT", "TMPDIR", "UID",
	},
	Symbols1: []string{ // ~ assignment
		"--", "++", "=", "*=", "/=", "+=", "-=", "<<=", ">>=", "&=", "^=", "|=",
	},
	Symbols2: []string{ // ~ comparators
		"==", "!=", "<", ">", " -eq ", " -ne ", " -lt ", " -gt ", " -ge ",
		" -a ", " -b ", " -c ", " -d ", " -e ", " -f ", " -g ", " -h ", " -k ",
		" -p ", " -r ", " -s ", " -t ", " -u ", " -w ", " -x ", " -G ", " -L ",
		" -N ", " -O ", " -S ", " -ef ", " -nt ", " -ot ", " -o ", " - v ",
		" -R ", " -z ", " -n ", "<=", ">=", "&&", "||",
	},
	Symbols3: []string{ // others
		"+", "-", "!", "~", "**", "*", "/", "%", ">>", "<<", "&", "^", "|",
		"$", "$#", "$@", "$?", "$$", "$-", "$_", "$!", "$*", "$0", "$1", "$2",
		"$3", "$4", "$5", "$6", "$", "$8", "$9",
	},
	Separators1: []string{
		"(", ")", "[", "]", "{", "}",
	},
	Separators2: []string{
		",", "?",
	},
	Separators3: []string{
		"[[", "]]", "((", "))",
	},
}
