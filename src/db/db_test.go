package db

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name, userID, resourceID, role string
	mockSetup                      func(mock sqlmock.Sqlmock)
	expectedError                  error
	expectedResult                 string
}

func setupTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, DB, func()) {
	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err, "Failed to open mock DB connection")

	sqlDB := NewDB(dbConn)
	cleanup := func() {
		sqlDB.Close()
		require.NoError(t, mock.ExpectationsWereMet(), "DB expectations not met")
	}

	return dbConn, mock, sqlDB, cleanup
}

func TestAssignRole(t *testing.T) {
	tests := []testCase{
		{
			name:       "success",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			role:       "admin",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(addRoleQuery)).
					WithArgs("user-uuid", "admin", "resource-uuid").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:          "invalid input - empty userID",
			userID:        "",
			resourceID:    "resource-uuid",
			role:          "admin",
			expectedError: errors.New("invalid input"),
		},
		{
			name:       "database error",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			role:       "admin",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(addRoleQuery)).
					WithArgs("user-uuid", "admin", "resource-uuid").
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("failed to assign role"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock, db, cleanup := setupTest(t)
			defer cleanup()

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			err := db.AssignRole(tt.userID, tt.role, tt.resourceID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetRole(t *testing.T) {
	tests := []testCase{
		{
			name:       "success",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"role"}).AddRow("admin")
				mock.ExpectQuery(regexp.QuoteMeta(selectRoleQuery)).
					WithArgs("user-uuid", "resource-uuid").
					WillReturnRows(rows)
			},
			expectedResult: "admin",
		},
		{
			name:          "invalid input - empty userID",
			userID:        "",
			resourceID:    "resource-uuid",
			expectedError: errors.New("invalid input"),
		},
		{
			name:       "no rows",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectRoleQuery)).
					WithArgs("user-uuid", "resource-uuid").
					WillReturnRows(sqlmock.NewRows([]string{"role"}))
			},
			expectedResult: "",
		},
		{
			name:       "database error",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(selectRoleQuery)).
					WithArgs("user-uuid", "resource-uuid").
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("failed to get role"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock, db, cleanup := setupTest(t)
			defer cleanup()

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			role, err := db.GetRole(tt.userID, tt.resourceID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, role)
			}
		})
	}
}

func TestRemoveRole(t *testing.T) {
	tests := []testCase{
		{
			name:       "success",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteRoleQuery)).
					WithArgs("user-uuid", "resource-uuid").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:          "invalid input - empty userID",
			userID:        "",
			resourceID:    "resource-uuid",
			expectedError: errors.New("invalid input"),
		},
		{
			name:       "no rows affected",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteRoleQuery)).
					WithArgs("user-uuid", "resource-uuid").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name:       "database error",
			userID:     "user-uuid",
			resourceID: "resource-uuid",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(deleteRoleQuery)).
					WithArgs("user-uuid", "resource-uuid").
					WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("failed to remove role"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mock, db, cleanup := setupTest(t)
			defer cleanup()

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			err := db.RemoveRole(tt.userID, tt.resourceID)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
