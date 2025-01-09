package services

func IsActionAllowed(role string, action string, mp ReadOnlyMap) bool {
	if actions, ok := mp.Load(role); ok {
		if actionSet, ok := actions.(map[string]struct{}); ok {
			_, exists := actionSet[action]
			return exists
		}
	}
	return false
}
