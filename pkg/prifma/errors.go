package prifma

import (
	"fmt"
	"strings"
)

type ErrWrongCall struct {
	Name string
	Args []string
}

func NewErrWrongCall(name string, args []string) *ErrWrongCall {
	return &ErrWrongCall{
		Name: name,
		Args: args,
	}
}

func (t *ErrWrongCall) Error() string {
	args := ""
	if t.Args != nil && len(t.Args) > 0 {
		args = "'" + strings.Join(t.Args, "', '") + "'"
	}

	return fmt.Sprintf("wrong call: %s(%s)", t.Name, args)
}
