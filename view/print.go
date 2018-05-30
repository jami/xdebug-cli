package view

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jami/xdebug-cli/dbgp"
)

// ShowPrintHelpMessage prints a short step help dialog
func (v *View) ShowPrintHelpMessage() {
	v.PrintLn("print -- prints out a variable")
}

func printProperty(v *View, p *dbgp.ProtocolProperty, indent int) {
	size := p.NumChildren
	ot := ""
	ct := ""
	// scalar
	if p.Children == 0 {
		size = p.Size
	}

	switch p.Type {
	case "array":
		ot = "["
		ct = "]"
	case "object":
		ot = "{"
		ct = "}"
	}

	fmtTypeDesc := strings.Repeat(" ", indent) + "%s %s(%d)"
	buffer := fmt.Sprintf(fmtTypeDesc, p.Name, p.Type, size)

	if len(ot) > 0 {
		buffer = buffer + " " + ot
		if size == 0 {
			buffer = buffer + ct
			ct = ""
		}
	}

	if p.Children == 0 && p.Encoding == "base64" && size > 0 {
		data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(p.Content))

		if err != nil {
			v.PrintErrorLn(err.Error())
		}

		value := string(data)
		if len(value) > 40 {
			value = value[:36] + " ..."
		}

		buffer = buffer + " = " + value
	}

	v.PrintLn(buffer)

	if p.Property != nil && p.Children == 1 {
		for _, c := range *p.Property {
			printProperty(v, &c, indent+4)
		}
	}

	if len(ct) > 0 {
		v.PrintLn(strings.Repeat(" ", indent) + ct)
	}
}

// PrintPropertyListWithDetails print
func (v *View) PrintPropertyListWithDetails(scopeName string, propertyList []dbgp.ProtocolProperty) {
	v.PrintLn("scope: " + scopeName)

	for _, p := range propertyList {
		printProperty(v, &p, 4)
	}
}
