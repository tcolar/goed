package syntax

var SyntaxMarkup = syntax{
	Extensions: []string{".html", ".htm", ".xml", ".xhtml"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern(`<!--`, `-->`, ``, true, StyleComment),
		NewSyntaxPattern(`<!`, `>`, ``, false, StyleKw3),    //directive
		NewSyntaxPattern(`<?`, `?>`, ``, false, StyleSymb2), // xml
		NewSyntaxPattern(`"`, `"`, `\`, false, StyleString),
		//NewSyntaxPattern(`'`, `'`, `\`, false, StyleString),
	},
	Keywords1: []string{
		"charset", "name", "value", "content", "id",
		"class", "hidden", "disabled", "meta", "style", "title",
	},
	Symbols1: []string{
		"=", ":",
	},
	Separators3: []string{
		"<", "</", "/>", ">",
	},
	Separators1: []string{
		".", ",",
	},
}
