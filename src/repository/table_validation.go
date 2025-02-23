package repository

import "regexp"

func isValidTableName(tableName string) bool {
	validTableName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validTableName.MatchString(tableName)
}
