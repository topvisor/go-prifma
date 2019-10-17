package conf

import "fmt"

type ParseError struct {
	Line int
	Data string
}

func NewParseError(line int, data string) error {
	return &ParseError{
		Line: line,
		Data: data,
	}
}

func (t *ParseError) Error() string {
	return fmt.Sprintf("parse error(line %d): %s", t.Line, t.Data)
}

type CommandError struct {
	Command Command
}

func NewCommandError(command Command) error {
	return &CommandError{
		Command: command,
	}
}

func (t *CommandError) Error() string {
	return fmt.Sprintf("wrong directive: %s", t.Command.String())
}
