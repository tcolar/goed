package syntax

var SyntaxMarkdown = syntax{
	Extensions: []string{".md"},
	Patterns: []SyntaxPattern{
		NewSyntaxPattern("```", "```", "", true, StyleKw2).WithMSL(), //Code
		NewSyntaxPattern("`", "`", "", false, StyleKw2),              //Code
		NewSyntaxPattern("    ", "", "", false, StyleKw2).WithMSL(),  //Code
		NewSyntaxPattern("#", "", "", false, StyleString).WithMSL(),  //Header
		NewSyntaxPattern("=", "", "", false, StyleString).WithMSL(),  //Header
		NewSyntaxPattern("---", "", "", false, StyleKw1).WithMSL(),   //HR
		NewSyntaxPattern("***", "", "", false, StyleKw1).WithMSL(),   //HR
		NewSyntaxPattern("___", "", "", false, StyleKw1).WithMSL(),   //HR
		NewSyntaxPattern(">", "", "", false, StyleComment).WithMSL(), //BlockQuote
		NewSyntaxPattern("**", "**", "", true, StyleSymb2),           //Bold
		NewSyntaxPattern("__", "__", "", true, StyleSymb2),           //Bold
	},
}
