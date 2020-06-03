package conf

import (
	"fmt"
	"strings"
)

type Command interface {
	GetLine() int
	GetName() string
	GetArgs() []string
	String() string
}

func NewCommand(line int, name string, args ...string) *DefaultCommand {
	return &DefaultCommand{
		Line: line,
		Name: name,
		Args: args,
	}
}

type DefaultCommand struct {
	Line int
	Name string
	Args []string
}

func (t *DefaultCommand) GetLine() int {
	return t.Line
}

func (t *DefaultCommand) GetName() string {
	return t.Name
}

func (t *DefaultCommand) GetArgs() []string {
	return t.Args
}

func (t *DefaultCommand) String() string {
	args := "()"
	if t.Args != nil && len(t.Args) > 0 {
		args = "('" + strings.Join(t.Args, "', '") + "')"
	}

	return fmt.Sprintf("%s%s (line %d)", t.Name, args, t.Line)
}
