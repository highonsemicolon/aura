package db

import (
	"database/sql"
	"fmt"
)

type assignRoleInput struct {
	UserID     string `validate:"required"`
	Role       string `validate:"required"`
	ResourceID string `validate:"required"`
}

type removeRoleInput struct {
	UserID     string `validate:"required"`
	ResourceID string `validate:"required"`
}

type getRoleInput struct {
	UserID     string `validate:"required"`
	ResourceID string `validate:"required"`
}

func (db *sqlDB) AssignRole(userID, role, resourceID string) error {
	if err := db.validateInput(assignRoleInput{UserID: userID, Role: role, ResourceID: resourceID}); err != nil {
		return err
	}

	_, err := db.conn.Exec(addRoleQuery, userID, role, resourceID)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

func (db *sqlDB) RemoveRole(userID, resourceID string) error {

	if err := db.validateInput(removeRoleInput{UserID: userID, ResourceID: resourceID}); err != nil {
		return err
	}

	_, err := db.conn.Exec(deleteRoleQuery, userID, resourceID)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}
	return err
}

func (db *sqlDB) GetRole(userID, resourceID string) (string, error) {

	if err := db.validateInput(getRoleInput{UserID: userID, ResourceID: resourceID}); err != nil {
		return "", err
	}

	var role string
	err := db.conn.QueryRow(selectRoleQuery, userID, resourceID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get role: %w", err)
	}
	return role, err
}


func (db *sqlDB) CheckIfResourceExists(resourceID string) (bool, error) {
	var exists bool
	err := db.conn.QueryRow(getResourceQuery, resourceID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check resource: %w", err)
	}
	return exists, nil
}