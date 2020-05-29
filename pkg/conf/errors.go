package conf

import (
	"fmt"
)

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
	Message string
}

func NewErrCommand(command Command, message string) error {
	return &ErrCommand{
		Command: command,
		Message: message,
	}
}

func NewErrCommandArgsNumber(command Command) error {
	return NewErrCommand(command, "wrong arguments number")
}

func NewErrCommandArg(command Command, arg string) error {
	return NewErrCommand(command, "wrong argument - "+arg)
}

func NewErrCommandName(command Command) error {
	return NewErrCommand(command, "wrong directive name")
}

func NewErrCommandMustHaveBlock(command Command) error {
	return NewErrCommand(command, "directive must have block ({})")
}

func NewErrCommandMustHaveNoBlock(command Command) error {
	return NewErrCommand(command, "directive must have no block ({})")
}

func (t *ErrCommand) Error() string {
	return t.Message + ": " + t.Command.String()
}
