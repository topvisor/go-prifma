package prifma_new

import "net/http"

type ModulesManager interface {
	GetModules(conds ...Condition) map[string]Module
	GetModulesForRequest(req *http.Request) map[string]Module
}

func NewModulesManager(modules ...Module) ModulesManager {
	mainModulesMap := make(map[string]Module, len(modules))
	for _, module := range modules {
		mainModulesMap[module.GetDirective()] = module
	}

	return &DefaultModulesManager{
		Modules:     mainModulesMap,
		CondModules: make(map[Condition]ModulesManager),
	}
}

type DefaultModulesManager struct {
	Modules     map[string]Module
	CondModules map[Condition]ModulesManager
}

func (t *DefaultModulesManager) GetModules(conds ...Condition) map[string]Module {
	if conds == nil || len(conds) == 0 {
		return t.Modules
	}

	cond := conds[0]
	conds = conds[1:]

	if _, ok := t.CondModules[cond]; !ok {
		modules := make([]Module, 0, len(t.Modules))
		for _, module := range t.Modules {
			modules = append(modules, module.Clone())
		}

		t.CondModules[cond] = NewModulesManager(modules...)
	}

	return t.CondModules[cond].GetModules(conds...)
}

func (t *DefaultModulesManager) GetModulesForRequest(req *http.Request) map[string]Module {
	for cond, manager := range t.CondModules {
		if cond.Test(req) {
			return manager.GetModulesForRequest(req)
		}
	}

	return t.Modules
}
