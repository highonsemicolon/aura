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
	RemoveRole(assignerID, userID, resourceID string) error
	GetRole(requesterRole, userID, resourceID string) (string, error)
}

type privilegeService struct {
	pc services.PrivilegeChecker
	DB db.DB
}

func NewPrivilegeService(pc services.PrivilegeChecker, db db.DB) *privilegeService {
	return &privilegeService{DB: db, pc: pc}
}

func (ps *privilegeService) AssignRole(assignerID, userID, role, resourceID string) error {
	if err := ps.validateInputs(assignerID, resourceID, userID); err != nil {
		return err
	}

	_, err := ps.GetRole(assignerID, assignerID, resourceID)
	if err != nil {
		return err
	}

	if !ps.pc.IsRoleAllowed(role) {
		return ErrRoleNotAllowed
	}

	return ps.DB.AssignRole(userID, role, resourceID)
}

func (ps *privilegeService) RemoveRole(assignerID, userID, resourceID string) error {
	if err := ps.validateInputs(assignerID, userID, resourceID); err != nil {
		return err
	}

	_, err := ps.GetRole(assignerID, assignerID, resourceID)
	if err != nil {
		return err
	}

	return ps.DB.RemoveRole(userID, resourceID)
}

func (ps *privilegeService) GetRole(requesterRole, userID, resourceID string) (string, error) {
	if err := ps.validateInputs(userID, resourceID); err != nil {
		return "", err
	}

	if !ps.pc.IsActionAllowed(requesterRole, SuperAction) {
		return "", ErrUnauthorized
	}
	return ps.DB.GetRole(userID, resourceID)
}

func (ps *privilegeService) validateInputs(input ...string) error {
	for _, s := range input {
		if s == "" {
			return ErrInvalidInput
		}
	}
	return nil
}
