package services

type PrivilegeChecker interface {
	IsActionAllowed(role string, action string) bool
	IsRoleAllowed(role string) bool
}

type Checker struct {
	fw *FileWatcher
}

func NewChecker(fw *FileWatcher) *Checker {
	return &Checker{fw: fw}
}

func (pc *Checker) IsActionAllowed(role, action string) bool {
	mp := pc.fw.GetEffectivePrivilegesCache()
	if actions, ok := mp.Load(role); ok {
		if actionSet, ok := actions.(map[string]struct{}); ok {
			_, exists := actionSet[action]
			return exists
		}
	}
	return false
}

func (pc *Checker) IsRoleAllowed(role string) bool {
	mp := pc.fw.GetEffectivePrivilegesCache()
	_, exists := mp.Load(role)
	return exists
}
