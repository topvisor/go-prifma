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

type ErrModuleDirectiveNotFound struct {
	Command conf.Command
}

func NewErrModuleDirectiveNotFound(command conf.Command) error {
	return &ErrModuleDirectiveNotFound{
		Command: command,
	}
}

func (t *ErrModuleDirectiveNotFound) Error() string {
	return fmt.Sprintf("directive not found: %s", t.Command.String())
}
