package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFindRole_ShouldFindAdmin(t *testing.T) {
	db := DbMock(t)

	UserRepo.Database = db
	UserRepo.AutoMigrate()
	UserRepo.InitiateModels()
	Expected := Role{
		ID:          uuid.NewString(),
		RoleName:    "ADMIN",
		Permissions: ADMIN,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rle, err := UserRepo.LoadRole("ADMIN")
	assert.Nil(t, err)
	Compare_Roles(t, Expected, rle)
}

func TestFindRole_ShouldFindNoPermissions(t *testing.T) {
	db := DbMock(t)

	UserRepo.Database = db
	UserRepo.AutoMigrate()
	UserRepo.InitiateModels()
	Expected := Role{
		ID:          uuid.NewString(),
		RoleName:    "NO_PERMISSIONS",
		Permissions: NO_PERMISSIONS,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	rle, err := UserRepo.LoadRole("NO_PERMISSIONS")
	assert.Nil(t, err)
	Compare_Roles(t, Expected, rle)
}

func TestAddRole_ShouldSucceed(t *testing.T) {
	db := DbMock(t)

	UserRepo.Database = db
	UserRepo.AutoMigrate()
	UserRepo.InitiateModels()
	Expected := Role{
		ID:          uuid.NewString(),
		RoleName:    "CATEGORY_WRITE",
		Permissions: CATEGORY_WRITE | CATEGORY_READ,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := UserRepo.SaveRole(Expected)
	assert.Nil(t, err)

	rle, err := UserRepo.LoadRole(Expected.RoleName)

	assert.Nil(t, err)
	Compare_Roles(t, Expected, rle)
}

func Compare_Roles(t *testing.T, Expected Role, Actual Role) {
	assert.Equal(t, Expected.RoleName, Actual.RoleName)
	assert.Equal(t, Expected.Permissions, Actual.Permissions)
}
