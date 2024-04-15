package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFindRole_ShouldFindAdmin(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	UserRepo.Database = db
	var RoleData []Role
	RoleData = append(RoleData, Role{
		ID:          uuid.NewString(),
		RoleName:    "ADMIN",
		Permissions: ADMIN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	roles, err := BuildMockDBSelectRows(RoleData)
	assert.Nil(t, err)

	expectedSQL, err := BuildSelectQuery(RoleData[0])
	assert.Nil(t, err)

	mock.ExpectQuery(expectedSQL).WillReturnRows(roles)

	_, err = UserRepo.LoadRole("ADMIN")
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindRole_ShouldFindNoPermissions(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	UserRepo.Database = db
	var RoleData []Role
	RoleData = append(RoleData, Role{
		ID:          uuid.NewString(),
		RoleName:    "NO_PERMISSIONS",
		Permissions: NO_PERMISSIONS,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	roles, err := BuildMockDBSelectRows(RoleData)
	assert.Nil(t, err)

	expectedSQL, err := BuildSelectQuery(RoleData[0])
	assert.Nil(t, err)

	mock.ExpectQuery(expectedSQL).WillReturnRows(roles)

	_, err = UserRepo.LoadRole("NO_PERMISSIONS")
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAddRole_ShouldSucceed(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	UserRepo.Database = db
	var RoleData []Role
	RoleData = append(RoleData, Role{
		ID:          uuid.NewString(),
		RoleName:    "CATEGORY_WRITE",
		Permissions: CATEGORY_WRITE | CATEGORY_READ,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	selectRoles, err := BuildMockRowsTableHeader(RoleData, false)
	assert.Nil(t, err)
	// insertRoles, err := BuildMockDBInsertRows(RoleData)
	// assert.Nil(t, err)

	expectedSQLSelect, err := BuildSelectQuery(RoleData[0])
	assert.Nil(t, err)
	expectedSQLInsert, err := BuildInsertQuery(RoleData[0])
	assert.Nil(t, err)

	UserRepo.Database = UserRepo.Database.Model(&Role{})

	mock.ExpectQuery(expectedSQLSelect).WillReturnRows(selectRoles)
	mock.ExpectQuery(expectedSQLInsert).WillReturnRows()

	err = UserRepo.SaveRole(RoleData[0])

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
