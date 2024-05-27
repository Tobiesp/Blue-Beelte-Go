package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFindUser_ShouldFindAdmin(t *testing.T) {
	db := DbMock(t)

	UserRepo.Database = db
	UserRepo.AutoMigrate()
	UserRepo.InitiateModels()
	r, err := UserRepo.LoadRole("ADMIN")
	assert.Nil(t, err)
	Expected := User{
		ID:                 uuid.NewString(),
		Username:           "admin",
		Email:              "admin@no.email",
		Role:               r,
		LoginAttempts:      0,
		LastLogin:          time.Now(),
		ForcePasswordReset: false,
		DisableAccount:     false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	usr, err := UserRepo.LogonUser("admin", "Password_1")
	if err != nil {
		lerr := err.(*LogonError)
		if lerr.ErrorCode() == 3 {
			err = UserRepo.ChangeUserPassword(usr, "Password_1", "Password_2")
			assert.Nil(t, err)
		}
		usr, err = UserRepo.LogonUser("admin", "Password_2")
	}
	assert.Nil(t, err)
	Compare_Users(t, Expected, usr)
}

func TestAddUser_ShouldSucceed(t *testing.T) {
	db := DbMock(t)

	UserRepo.Database = db
	UserRepo.AutoMigrate()
	UserRepo.InitiateModels()
	userTime := time.Now().UTC().String()
	Expected, err := UserRepo.CreateNewUser("newUser"+userTime, "newUser@no.email", "Password_1")
	assert.Nil(t, err)

	err = UserRepo.SaveUser(Expected)
	assert.Nil(t, err)

	rle, err := UserRepo.LoadUser(Expected.Username)

	assert.Nil(t, err)
	Compare_Users(t, Expected, rle)
}

func TestUserLogon_ShouldSucceed(t *testing.T) {
	db := DbMock(t)

	UserRepo.Database = db
	UserRepo.AutoMigrate()
	UserRepo.InitiateModels()

	usr, err := UserRepo.LogonUser("admin", "Password_1")
	if err != nil {
		lerr := err.(*LogonError)
		if lerr.ErrorCode() == 3 {
			err = UserRepo.ChangeUserPassword(usr, "Password_1", "Password_2")
			assert.Nil(t, err)
		}
		usr, err = UserRepo.LogonUser("admin", "Password_2")
	}
	assert.Nil(t, err)
	expected := usr.VerifyPassword("Password_2")
	assert.True(t, expected)
}

func TestUserLogon_ShouldFail(t *testing.T) {
	db := DbMock(t)

	UserRepo.Database = db
	UserRepo.AutoMigrate()
	UserRepo.InitiateModels()

	usr, err := UserRepo.LogonUser("admin", "Password_3")
	assert.NotNil(t, err)
	lerr := err.(*LogonError)
	assert.Equal(t, lerr.ErrorCode(), 1)
	assert.Equal(t, usr.LoginAttempts, uint(0))
}

func Compare_Users(t *testing.T, Expected User, Actual User) {
	assert.Equal(t, Expected.Username, Actual.Username)
	assert.Equal(t, Expected.Email, Actual.Email)
	Compare_Roles(t, Expected.Role, Actual.Role)
	assert.Equal(t, Expected.LoginAttempts, Actual.LoginAttempts)
	assert.Equal(t, Expected.ForcePasswordReset, Actual.ForcePasswordReset)
	assert.Equal(t, Expected.DisableAccount, Actual.DisableAccount)
}
