package dbgp

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CommandBuilder type
type CommandBuilder func(int, string) string

// CommandProcessor for term input handling
type CommandProcessor struct {
	reader   *bufio.Reader
	commands map[*regexp.Regexp]CommandBuilder
}

// GetCommand accumulates input runes
func (cp *CommandProcessor) GetCommand(transID int) string {
	buffer, _ := cp.reader.ReadString('\n')
	fmt.Println("Buffer", buffer)
	for matcher, builder := range cp.commands {
		if matcher.MatchString(buffer) {
			buffer = builder(transID, buffer)
			fmt.Println("Builder", buffer)
			break
		}
	}

	return strings.TrimSpace(buffer)
}

// NewCommandProcessor creator
func NewCommandProcessor() *CommandProcessor {
	cp := &CommandProcessor{
		reader:   bufio.NewReader(os.Stdin),
		commands: map[*regexp.Regexp]CommandBuilder{},
	}

	reg, b := createStepIntoCommand()
	cp.commands[reg] = b

	return cp
}

func createStatusIntoCommand() (*regexp.Regexp, CommandBuilder) {
	r, _ := regexp.Compile("status")

	return r, func(tid int, cmd string) string {
		return fmt.Sprintf("status -i %d", tid)
	}
}

func createStepIntoCommand() (*regexp.Regexp, CommandBuilder) {
	r, _ := regexp.Compile("step")

	return r, func(tid int, cmd string) string {
		return "huihuihui"
	}
}
