package view

// ShowHelpMessage prints a short overall help dialog
func (view *View) ShowHelpMessage() {
	view.PrintLn("List of commands:")
	view.PrintLn("")
	view.PrintLn("info|i - get debug context information")
	view.PrintLn("break|b - set debug breakpoints")
	view.PrintLn("list|l - list the current code lines")
	view.PrintLn("run|r - run the code")
	view.PrintLn("step|s - step into the next instruction")
	view.PrintLn("next|n - step over instruction")
	view.PrintLn("finish|f - closes the debug session")
	view.PrintLn("print|p - prints out a variable")
	view.PrintLn("context|c - context variable dump")
}
