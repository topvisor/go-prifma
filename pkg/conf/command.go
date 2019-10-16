package conf

import "strings"

type Command interface {
	GetName() string
	GetArgs() []string
	String() string
}

func NewCommand(name string, args ...string) Command {
	return &DefaultCommand{
		Name: name,
		Args: args,
	}
}

type DefaultCommand struct {
	Name string
	Args []string
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

	return t.Name + args
}
