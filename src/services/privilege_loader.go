package services

import (
	"aura/src/utils"
	"sync"
)

type ReadOnlyMap interface {
	Load(key interface{}) (value interface{}, ok bool)
	Range(f func(key, value interface{}) bool)
}

type PrivilegeLoader struct {
	filename                 string
	effectivePrivilegesCache *sync.Map
	cacheMutex               sync.Mutex
}

func NewPrivilegeLoader(filename string) *PrivilegeLoader {
	return &PrivilegeLoader{
		filename:                 filename,
		effectivePrivilegesCache: &sync.Map{},
	}
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
	var newPrivileges sync.Map
	for role := range privileges.Roles {
		effective := utils.ComputeRolePrivilegesDFS(role, privileges, make(map[string]bool))
		newPrivileges.Store(role, effective)
	}

	pl.cacheMutex.Lock()
	defer pl.cacheMutex.Unlock()

	pl.effectivePrivilegesCache.Range(func(key, value interface{}) bool {
		pl.effectivePrivilegesCache.Delete(key)
		return true
	})
	newPrivileges.Range(func(key, value interface{}) bool {
		pl.effectivePrivilegesCache.Store(key, value)
		return true
	})
}

func (pl *PrivilegeLoader) GetEffectivePrivileges(role string) ([]string, bool) {

	privileges, ok := pl.effectivePrivilegesCache.Load(role)
	if !ok {
		return nil, false
	}
	return privileges.([]string), true
}

func (pl *PrivilegeLoader) GetEffectivePrivilegesCache() ReadOnlyMap {
	return pl.effectivePrivilegesCache
}
