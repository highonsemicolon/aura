package services

func IsActionAllowed(role string, action string, fw *FileWatcher) bool {
	mp := fw.GetEffectivePrivilegesCache()
	if actions, ok := mp.Load(role); ok {
		if actionSet, ok := actions.(map[string]struct{}); ok {
			_, exists := actionSet[action]
			return exists
		}
	}
	return false
}
