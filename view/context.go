package view

// ShowContextHelpMessage prints a short context help dialog
func (view *View) ShowContextHelpMessage() {
	view.PrintLn("context -- shows variables information")
	view.PrintLn("")
	view.PrintLn("context local - show variables from the local scope")
	view.PrintLn("context global - show variables from the global scope")
	view.PrintLn("context constant - show user defined constants")

}
