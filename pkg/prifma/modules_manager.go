package prifma

import "net/http"

type ModulesManager interface {
	GetModule(directive string, conds ...Condition) Module
	GetModulesForRequest(req *http.Request) []Module
}

func NewModulesManager(modules ...Module) *DefaultModulesManager {
	mainModulesMap := make(map[string]int, len(modules))
	for i, module := range modules {
		mainModulesMap[module.GetDirective()] = i
	}

	return &DefaultModulesManager{
		ModulesArray: modules,
		ModulesMap:   mainModulesMap,
		CondModules:  make(map[Condition]ModulesManager),
	}
}

type DefaultModulesManager struct {
	ModulesArray []Module
	ModulesMap   map[string]int
	CondModules  map[Condition]ModulesManager
}

func (t *DefaultModulesManager) GetModule(directive string, conds ...Condition) Module {
	if conds == nil || len(conds) == 0 {
		return t.ModulesArray[t.ModulesMap[directive]]
	}

	cond := conds[0]
	conds = conds[1:]

	if _, ok := t.CondModules[cond]; !ok {
		modules := make([]Module, len(t.ModulesArray))
		for i, module := range t.ModulesArray {
			modules[i] = module.Clone()
		}

		t.CondModules[cond] = NewModulesManager(modules...)
	}

	return t.CondModules[cond].GetModule(directive, conds...)
}

func (t *DefaultModulesManager) GetModulesForRequest(req *http.Request) []Module {
	for cond, manager := range t.CondModules {
		if cond.Test(req) {
			return manager.GetModulesForRequest(req)
		}
	}

	return t.ModulesArray
}
