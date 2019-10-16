package prifma_new

import "net/http"

type ModulesManager interface {
	AddModules(cond Condition, modules ...Module)
	GetModules(req *http.Request) []Module
}

func NewModulesManager(modules ...Module) ModulesManager {
	return &DefaultModulesManager{
		Modules: map[CompiledCondition][]Module{
			nil: modules,
		},
	}
}

type DefaultModulesManager struct {
	Modules map[CompiledCondition][]Module
}

func (t *DefaultModulesManager) AddModules(cond Condition, modules ...Module) {
	t.Modules[cond] = append(t.Modules[cond], modules...)
}

func (t *DefaultModulesManager) GetModules(req *http.Request) []Module {
	if req == nil {
		return t.Modules[nil]
	}

	panic("implement me")
}
