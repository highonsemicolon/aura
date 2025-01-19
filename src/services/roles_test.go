package services

import (
	"testing"

	role "aura/src/services/role"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) AssignRole(userID, resourceID, role string) error {
	args := m.Called(userID, resourceID, role)
	return args.Error(0)
}

func (m *MockDB) RemoveRole(userID, resourceID string) error {
	args := m.Called(userID, resourceID)
	return args.Error(0)
}

func (m *MockDB) GetRole(userID, resourceID string) (string, error) {
	args := m.Called(userID, resourceID)
	return args.String(0), args.Error(1)
}

type MockChecker struct {
	mock.Mock
}

func NewMockChecker() role.PrivilegeChecker {
	return &MockChecker{}
}

func (m *MockChecker) IsActionAllowed(role, action string) bool {
	args := m.Called(role, action)
	return args.Get(0).(bool)
}

func (m *MockChecker) IsRoleAllowed(role string) bool {
	args := m.Called(role)
	return args.Get(0).(bool)
}

type testCase struct {
	name, assignerID, targetUserID, role, resourceID string
	setupMock                                        func(*MockDB, *MockChecker)
	expectedError                                    error
}

func TestAssignRole(t *testing.T) {
	tests := []testCase{
		{
			name:         "success - owner assigning role",
			assignerID:   "admin-id",
			targetUserID: "user-id",
			role:         "editor",
			resourceID:   "resource-id",
			setupMock: func(m *MockDB, pc *MockChecker) {
				m.On("GetRole", "admin-id", "resource-id").Return("owner", nil)
				pc.On("IsActionAllowed", mock.Anything, mock.Anything).Return(true)
				pc.On("IsRoleAllowed", mock.Anything).Return(true)
				m.On("AssignRole", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:         "error - invalid role",
			assignerID:   "admin-id",
			targetUserID: "user-id",
			role:         "invalid-role",
			resourceID:   "resource-id",
			setupMock: func(m *MockDB, pc *MockChecker) {
				m.On("GetRole", mock.Anything, mock.Anything, mock.Anything).Return("admin", nil)
				pc.On("IsActionAllowed", mock.Anything, mock.Anything).Return(true)
				pc.On("IsRoleAllowed", "invalid-role").Return(false)

			},
			expectedError: ErrRoleNotAllowed,
		},
		{
			name:         "error - unauthorized assigner",
			assignerID:   "editor-id",
			targetUserID: "user-id",
			role:         "editor",
			resourceID:   "resource-id",
			setupMock: func(m *MockDB, pc *MockChecker) {
				pc.On("IsActionAllowed", mock.Anything, mock.Anything).Return(false)

			},
			expectedError: ErrUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := new(MockDB)
			pc := new(MockChecker)
			service := NewPrivilegeService(pc, store)

			if tt.setupMock != nil {
				tt.setupMock(store, pc)
			}

			err := service.AssignRole(tt.assignerID, tt.targetUserID, tt.role, tt.resourceID)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			store.AssertExpectations(t)
		})
	}
}

func TestPrivilegeService_RemoveRole(t *testing.T) {
	tests := []testCase{
		{
			name:         "success - owner removing role",
			assignerID:   "admin-id",
			targetUserID: "user-id",
			role:         "editor",
			resourceID:   "resource-id",
			setupMock: func(m *MockDB, pc *MockChecker) {
				m.On("GetRole", "admin-id", "resource-id").Return("owner", nil)
				pc.On("IsActionAllowed", mock.Anything, mock.Anything).Return(true)
				m.On("RemoveRole", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:          "error - invalid input",
			assignerID:    "",
			targetUserID:  "user-id",
			role:          "editor",
			resourceID:    "resource-id",
			setupMock:     nil,
			expectedError: ErrInvalidInput,
		},
		{
			name:         "error - unauthorized assigner",
			assignerID:   "editor-id",
			targetUserID: "user-id",
			role:         "editor",
			resourceID:   "resource-id",
			setupMock: func(m *MockDB, pc *MockChecker) {
				pc.On("IsActionAllowed", mock.Anything, mock.Anything).Return(false)

			},
			expectedError: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := new(MockDB)
			pc := new(MockChecker)
			service := NewPrivilegeService(pc, store)

			if tt.setupMock != nil {
				tt.setupMock(store, pc)
			}

			err := service.RemoveRole(tt.assignerID, tt.targetUserID, tt.resourceID)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			store.AssertExpectations(t)
		})
	}
}

func TestPrivilegeService_GetRole(t *testing.T) {
	tests := []testCase{
		{
			name:         "success - get role",
			assignerID:   "admin-id",
			targetUserID: "user-id",
			role:         "editor",
			resourceID:   "resource-id",
			setupMock: func(m *MockDB, pc *MockChecker) {
				m.On("GetRole", "admin-id", "resource-id").Return("owner", nil)
				pc.On("IsActionAllowed", mock.Anything, mock.Anything).Return(true)
			},
		},
		{
			name:          "error - invalid input",
			assignerID:    "",
			targetUserID:  "user-id",
			role:          "editor",
			resourceID:    "resource-id",
			setupMock:     nil,
			expectedError: ErrInvalidInput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := new(MockDB)
			pc := new(MockChecker)
			service := NewPrivilegeService(pc, store)

			if tt.setupMock != nil {
				tt.setupMock(store, pc)
			}

			_, err := service.GetRole(tt.assignerID, tt.assignerID, tt.resourceID)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			store.AssertExpectations(t)
		})
	}
}
