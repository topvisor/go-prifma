package prifma_new

import (
	"fmt"
	"github.com/topvisor/prifma/pkg/conf"
)

type ErrWrongDirective struct {
	Command conf.Command
}

func NewErrWrongDirective(command conf.Command) error {
	return &ErrWrongDirective{
		Command: command,
	}
}

func (t *ErrWrongDirective) Error() string {
	return fmt.Sprintf("wrong directive: %s", t.Command.String())
}
