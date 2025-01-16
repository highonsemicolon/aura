package db

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, DB) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	sqlDB := NewDB(dbConn)
	return dbConn, mock, sqlDB
}

func TestAssignRole_Success(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	mock.ExpectExec(regexp.QuoteMeta(assignRoleQuery)).
		WithArgs("user-uuid", "admin", "resource-uuid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := sqlDB.AssignRole("user-uuid", "admin", "resource-uuid")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAssignRole_InvalidInput(t *testing.T) {
	_, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	err := sqlDB.AssignRole("", "admin", "resource-uuid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestAssignRole_DBError(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	mock.ExpectExec(regexp.QuoteMeta(assignRoleQuery)).
		WithArgs("user-uuid", "admin", "resource-uuid").
		WillReturnError(errors.New("database error"))

	err := sqlDB.AssignRole("user-uuid", "admin", "resource-uuid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to assign role")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveRole_Success(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	mock.ExpectExec(regexp.QuoteMeta(removeRoleQuery)).
		WithArgs("user-uuid", "resource-uuid").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := sqlDB.RemoveRole("user-uuid", "resource-uuid")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveRole_InvalidInput(t *testing.T) {
	_, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	err := sqlDB.RemoveRole("", "resource-uuid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestRemoveRole_NoRowsAffected(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	mock.ExpectExec(regexp.QuoteMeta(removeRoleQuery)).
		WithArgs("user-uuid", "resource-uuid").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := sqlDB.RemoveRole("user-uuid", "resource-uuid")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveRole_DBError(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	mock.ExpectExec(regexp.QuoteMeta(removeRoleQuery)).
		WithArgs("user-uuid", "resource-uuid").
		WillReturnError(errors.New("database error"))

	err := sqlDB.RemoveRole("user-uuid", "resource-uuid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to remove role")
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetRole_Success(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	rows := sqlmock.NewRows([]string{"role"}).AddRow("admin")
	mock.ExpectQuery(regexp.QuoteMeta(getRoleQuery)).
		WithArgs("user-uuid", "resource-uuid").
		WillReturnRows(rows)

	role, err := sqlDB.GetRole("user-uuid", "resource-uuid")
	assert.NoError(t, err)
	assert.Equal(t, "admin", role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRole_InvalidInput(t *testing.T) {
	_, _, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	_, err := sqlDB.GetRole("", "resource-uuid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestGetRole_NoRows(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(getRoleQuery)).
		WithArgs("user-uuid", "resource-uuid").
		WillReturnRows(sqlmock.NewRows([]string{"role"}))

	role, err := sqlDB.GetRole("user-uuid", "resource-uuid")
	assert.NoError(t, err)
	assert.Equal(t, "", role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRole_DBError(t *testing.T) {
	_, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(getRoleQuery)).
		WithArgs("user-uuid", "resource-uuid").
		WillReturnError(errors.New("database error"))

	_, err := sqlDB.GetRole("user-uuid", "resource-uuid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get role")
	assert.NoError(t, mock.ExpectationsWereMet())
}
