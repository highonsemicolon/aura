package services

import (
	"aura/src/utils"
	"sync"
	"sync/atomic"
)

type PrivilegeLoader struct {
	filename                 string
	effectivePrivilegesCache atomic.Pointer[sync.Map]
}

func NewPrivilegeLoader(filename string) *PrivilegeLoader {
	pl := PrivilegeLoader{
		filename: filename,
	}

	pl.effectivePrivilegesCache.Store(&sync.Map{})
	return &pl
}

func (pl *PrivilegeLoader) LoadAndComputePrivileges() error {
	privileges, err := utils.LoadPrivileges(pl.filename)
	if err != nil {
		return err
	}
	pl.computeEffectivePrivileges(privileges)
	return nil
}

func (pl *PrivilegeLoader) computeEffectivePrivileges(privileges *utils.RolePrivileges) {
	newPrivileges := &sync.Map{}
	for role := range privileges.Roles {
		effective := utils.ComputeRolePrivilegesDFS(role, privileges, make(map[string]bool))
		newPrivileges.Store(role, effective)
	}

	pl.effectivePrivilegesCache.Store(newPrivileges)
}

func (pl *PrivilegeLoader) GetEffectivePrivileges(role string) ([]string, bool) {

	effectivePrivileges := pl.effectivePrivilegesCache.Load()
	privileges, ok := effectivePrivileges.Load(role)
	if !ok {
		return nil, false
	}
	return privileges.([]string), true
}

func (pl *PrivilegeLoader) GetEffectivePrivilegesCache() *sync.Map {
	return pl.effectivePrivilegesCache.Load()
}
