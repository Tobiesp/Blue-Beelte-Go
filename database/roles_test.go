package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFindRole_ShouldFind(t *testing.T) {
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

	expectedSQL := "SELECT (.+) FROM \"roles\" WHERE id =(.+)"

	mock.ExpectQuery(expectedSQL).WillReturnRows(roles)

	var role Role
	err = role.Load("ADMIN")
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
