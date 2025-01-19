package services

import (
	"aura/src/db"
	"errors"

	services "aura/src/services/role"
)

var (
	SuperAction       = "assign"
	ErrInvalidInput   = errors.New("invalid input parameters")
	ErrRoleNotAllowed = errors.New("role not allowed")
	ErrUnauthorized   = errors.New("unauthorized to perform this action")
)

type PrivilegeServiceInterface interface {
	AssignRole(assignerID, userID, role, resourceID string) error
	RemoveRole(userID, resourceID string) error
	GetRole(userID, resourceID string) (string, error)
}

type PrivilegeService struct {
	pc services.PrivilegeChecker
	DB db.DB
}

func NewPrivilegeService(pc services.PrivilegeChecker, db db.DB) PrivilegeServiceInterface {
	return &PrivilegeService{DB: db, pc: pc}
}

func (ps *PrivilegeService) AssignRole(assignerID, userID, role, resourceID string) error {
	if err := ps.validateInputs(assignerID, resourceID); err != nil {
		return err
	}

	assignerRole, err := ps.GetRole(assignerID, resourceID)
	if err != nil {
		return err
	}

	if !ps.pc.IsActionAllowed(assignerRole, SuperAction) {
		return ErrUnauthorized
	}

	if !ps.pc.IsRoleAllowed(role) {
		return ErrRoleNotAllowed
	}

	return ps.DB.AssignRole(userID, role, resourceID)
}

func (ps *PrivilegeService) RemoveRole(userID, resourceID string) error {
	return ps.DB.RemoveRole(userID, resourceID)
}

func (ps *PrivilegeService) GetRole(userID, resourceID string) (string, error) {
	return ps.DB.GetRole(userID, resourceID)
}

func (ps *PrivilegeService) validateInputs(assignerID, resourceID string) error {
	if assignerID == "" || resourceID == "" {
		return ErrInvalidInput
	}
	return nil
}
