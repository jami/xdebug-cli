package view

import (
	"fmt"
	"path/filepath"

	"github.com/jami/xdebug-cli/dbgp"
)

// ShowInfoHelpMessage prints a short info help dialog
func (view *View) ShowInfoHelpMessage() {
	view.PrintLn("info -- shows different information")
	view.PrintLn("")
	view.PrintLn("info breakpoints - list all breakpoints")
}

// ShowInfoBreakpoints prints a short info help dialog
func (view *View) ShowInfoBreakpoints(bpl []dbgp.ProtocolBreakpoint) {
	enabledShort := func(e string) string {
		if e == "enabled" {
			return "yes"
		}
		return "no"
	}

	fileShort := func(e string, l int) string {
		return fmt.Sprintf("%s:%d", filepath.Base(e), l)
	}

	view.PrintLn(fmt.Sprintf("%-4s%-10s%-8s%-10s", "Num", "Type", "Enabled", "What"))
	for idx, bp := range bpl {
		view.PrintLn(
			fmt.Sprintf(
				"#%-3d%-10s%-8s%-10s",
				idx,
				bp.Type,
				enabledShort(bp.State),
				fileShort(bp.FileName, bp.Line)))
	}
}
