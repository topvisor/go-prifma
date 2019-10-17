package conf

import "fmt"

type ErrParse struct {
	Line int
	Data string
}

func NewErrParse(line int, data string) error {
	return &ErrParse{
		Line: line,
		Data: data,
	}
}

func (t *ErrParse) Error() string {
	return fmt.Sprintf("parse error(line %d): %s", t.Line, t.Data)
}

type ErrCommand struct {
	Command Command
}

func NewErrCommand(command Command) error {
	return &ErrCommand{
		Command: command,
	}
}

func (t *ErrCommand) Error() string {
	return fmt.Sprintf("wrong directive: %s", t.Command.String())
}
