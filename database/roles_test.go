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
		ID:          uuid.New(),
		RoleName:    "ADMIN",
		Permissions: ADMIN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	roles, err := BuildMockDBRows(RoleData)
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
		ID:          uuid.New(),
		RoleName:    "NO_PERMISSIONS",
		Permissions: NO_PERMISSIONS,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	roles, err := BuildMockDBRows(RoleData)
	assert.Nil(t, err)

	expectedSQL, err := BuildSelectQuery(RoleData[0])
	assert.Nil(t, err)

	mock.ExpectQuery(expectedSQL).WillReturnRows(roles)

	_, err = UserRepo.LoadRole("NO_PERMISSIONS")
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
