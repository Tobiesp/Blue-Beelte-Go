package models

import (
	"blue-beetle/database"
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username  string `gorm:"primaryKey"`
	Email     string
	Password  []byte
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func CreateNewUser(username string, password string) (User, error) {
	validate := validateUsername(username)
	var u User
	if validate != nil {
		return u, validate
	}
	u.Username = username
	u.Password = encryptPassword(password)
	var noPerms Role
	noPerms.Load("NO_PERMISSIONS")
	u.Role = noPerms
	return u, nil
}

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username can not be empty")
	}
	var exists bool
	err := database.Instance.Model(&User{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists).
		Error
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exits")
	}
	return nil
}

func encryptPassword(password string) []byte {
	//TODO: Add code here to encrypt the password
	return []byte(password)
}

func (u *User) Save() error {
	record := database.Instance.Where("username = ?", u.Username).First(&u)
	if record.Error != nil {
		record = database.Instance.Create(&u)
	} else {
		record = database.Instance.Save(&u)
	}
	err := record.Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Load(username string) error {
	record := database.Instance.Where("username = ?", u.Username).First(&u)
	if record.Error != nil {
		return record.Error
	}
	return nil
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role.Permissions == ADMIN {
		return errors.New("admin users are not allowed to be deleted")
	}
	return
}
