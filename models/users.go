package models

import (
	"blue-beetle/database"
	"bytes"
	"errors"
	"time"
	"unicode"

	"math/rand"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username           string `gorm:"primaryKey"`
	Email              string
	Password           []byte
	Role               Role
	LoginAttempts      uint
	LastLogin          time.Time
	ForcePasswordReset bool
	DisableAccount     bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func CreateNewUser(username string, password string) (User, error) {
	validate := validateUsername(username)
	var u User
	if validate != nil {
		return u, validate
	}
	u.Username = username
	var noPerms Role
	noPerms.Load("NO_PERMISSIONS")
	u.Role = noPerms
	err := u.EncryptPassword(password)
	if err != nil {
		return u, err
	}
	u.Save()
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

func validatePassword(password string) error {
	// Must have at least 1 Upper case
	// Must have at least 1 lower case
	// Must have at least 1 symbol
	// Must have at least 1 number
	// Must be at least 8 chars long
	if len(password) < 8 {
		return errors.New("password must be at least 8 charators long")
	}
	var flags uint16 = 0x0
	for _, c := range password {
		if unicode.IsUpper(c) {
			flags = flags | 0x0001
		}
		if unicode.IsLower(c) {
			flags = flags | 0x0010
		}
		if unicode.IsNumber(c) {
			flags = flags | 0x0100
		}
		if !unicode.IsLetter(c) {
			flags = flags | 0x1000
		}
	}
	if (flags & 0x0001) == 0x0000 {
		return errors.New("password must have at least 1 uppercase letter")
	}
	if (flags & 0x0010) == 0x0000 {
		return errors.New("password must have at least 1 lowercase letter")
	}
	if (flags & 0x0100) == 0x0000 {
		return errors.New("password must have at least 1 number")
	}
	if (flags & 0x1000) == 0x0000 {
		return errors.New("password must have at least 1 symbole")
	}
	return nil
}

func (u *User) VerifyPassword(password string) bool {
	pass, err := encryptPassword(password)
	if err != nil {
		return false
	}
	return bytes.Equal(u.Password, pass)
}

func (u *User) EncryptPassword(password string) error {
	err := validatePassword(password)
	if err != nil {
		return err
	}
	ep, err := encryptPassword(password)
	if err != nil {
		return err
	}
	u.Password = ep
	return nil
}

func (u *User) ChangePassword(oldPassword string, newPassword string) error {
	if u.VerifyPassword(oldPassword) {
		err := u.EncryptPassword(newPassword)
		if err != nil {
			return err
		}
	} else {
		return errors.New("current password doesn't match for the user")
	}
	return nil
}

func (u *User) ResetPassword() (string, error) {
	newPass := GenerateRandmoPassword()
	count := 0
	for count < 1000 {
		err := validatePassword(newPass)
		if err == nil {
			break
		}
		count += 1
		newPass = GenerateRandmoPassword()
	}
	err := validatePassword(newPass)
	if err == nil {
		return "", err
	}
	err = u.EncryptPassword(newPass)
	if err != nil {
		return "", err
	}
	err = u.Save()
	if err != nil {
		return "", err
	}
	return newPass, nil
}

func encryptPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return bytes, nil
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

func GenerateRandmoPassword() string {
	literalList := "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890_*#^&@:<>.,?+="
	var s string
	for i := 1; i < 16; i++ {
		random := rand.Intn(len(literalList))
		s += string(literalList[random])
	}
	return s
}

func MigrateUserModel() error {
	return database.Instance.AutoMigrate(&User{})
}

func InitUserModle() error {
	err := InitRoleModle()
	if err != nil {
		return err
	}
	var admin User
	admin.Email = "admin@no.email"
	admin.Username = "admin"
	err = admin.EncryptPassword("Password_1")
	if err != nil {
		return err
	}
	var adminRole Role
	adminRole.Load("ADMIN")
	admin.Role = adminRole
	err = admin.Save()
	if err != nil {
		return err
	}
	return nil
}
