package view

// ShowBreakpointHelpMessage prints a short breakpoint help dialog
func (view *View) ShowBreakpointHelpMessage() {
	view.PrintLn("break -- set breakpoints")
	view.PrintLn("")
	view.PrintLn("break :n - set breakpoint to current file at line n")
	view.PrintLn("break <file>:n - set breakpoint to <file> at line n")
	view.PrintLn("break call <func> - set breakpoint to function call <func>")
}
