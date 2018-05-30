package view

// ShowStepHelpMessage prints a short step help dialog
func (view *View) ShowStepHelpMessage() {
	view.PrintLn("step -- steps to the next statement," +
		" if there is a function call involved it will break" +
		" on the first statement in that function")
}
