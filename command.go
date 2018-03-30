package xdebugcli

type CommandProcessor struct {
}

func (cp *CommandProcessor) Process(msg string) {

}

func NewCommandProcessor() *CommandProcessor {
	cp := &CommandProcessor{}
	return cp
}
