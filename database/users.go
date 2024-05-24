package database

import (
	"bytes"
	"context"
	"errors"
	"time"
	"unicode"

	"math/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                 string `gorm:"primaryKey"`
	Username           string `gorm:"not null,type:text"`
	Email              string
	Password           []byte
	Role               Role `gorm:"embedded"`
	LoginAttempts      uint
	LastLogin          time.Time `gorm:"embedded"`
	ForcePasswordReset bool
	DisableAccount     bool
	CreatedAt          time.Time      `gorm:"embedded"`
	UpdatedAt          time.Time      `gorm:"embedded"`
	DeletedAt          gorm.DeletedAt `gorm:"index;embedded"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.NewString()
	tx.Statement.SetColumn("ID", user.ID)
	return
}

func (r *UserRepository) LogonUser(username string, password string) (User, error) {
	user, err := r.LoadUser(username)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, LogonErrorNew("Username or password is not valid", BAD_USER_CODE)
	}
	if !user.VerifyPassword(password) {
		user.LoginAttempts = user.LoginAttempts + 1
		err = r.SaveUser(user)
		if err != nil {
			return User{}, LogonErrorNew("Fail to save User account: "+err.Error(), FAILED_TO_SAVE_USER_CODE)
		}
		err = LogonErrorNew("Username or password is not valid", BAD_USER_CODE)
		return User{}, err
	}
	if user.DisableAccount {
		return user, LogonErrorNew("Account Locked", LOCKED_ACCOUNT_CODE)
	}
	if user.ForcePasswordReset {
		return user, LogonErrorNew("New Password need", FORCED_PASS_RESET_CODE)
	}
	if user.LoginAttempts > 5 {
		user.LoginAttempts = user.LoginAttempts + 1
		user.LastLogin = time.Now()
		err = r.SaveUser(user)
		if err != nil {
			return User{}, LogonErrorNew("Fail to save User account: "+err.Error(), FAILED_TO_SAVE_USER_CODE)
		}
		return user, LogonErrorNew("To many Failed Logins", LOGON_COUNT_FAILED_CODE)
	}
	user.LoginAttempts = 0
	user.LastLogin = time.Now()
	err = r.SaveUser(user)
	if err != nil {
		return User{}, LogonErrorNew("Fail to save User account: "+err.Error(), FAILED_TO_SAVE_USER_CODE)
	}
	return user, nil
}

func (r *UserRepository) CreateNewUser(username string, email string, password string) (User, error) {
	validate := r.validateUsername(username)
	var u User
	if validate != nil {
		return u, validate
	}
	u.Username = username
	u.Email = email
	u.DisableAccount = false
	u.ForcePasswordReset = false
	u.LastLogin = time.Now()
	u.LoginAttempts = 0
	noPerms, err := UserRepo.LoadRole("NO_PERMISSIONS")
	if err != nil {
		return u, err
	}
	u.Role = noPerms
	err = u.EncryptPassword(password)
	if err != nil {
		return u, err
	}
	err = r.SaveUser(u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (r *UserRepository) validateUsername(username string) error {
	if username == "" {
		return errors.New("username can not be empty")
	}
	var exists bool
	err := r.Database.Model(&User{}).
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

func (r *UserRepository) ChangeUserPassword(user User, oldPassword string, newPassword string) error {
	if user.VerifyPassword(oldPassword) {
		if oldPassword == newPassword {
			return errors.New("current and new password can't be the same")
		}
		err := user.EncryptPassword(newPassword)
		if err != nil {
			return err
		}
	} else {
		return errors.New("current password doesn't match for the user")
	}
	err := r.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) ResetUserPassword(user User) (string, error) {
	newPass := UserRepo.GenerateRandmoPassword()
	count := 0
	for count < 1000 {
		err := validatePassword(newPass)
		if err == nil {
			break
		}
		count += 1
		newPass = UserRepo.GenerateRandmoPassword()
	}
	err := validatePassword(newPass)
	if err == nil {
		return "", err
	}
	err = user.EncryptPassword(newPass)
	if err != nil {
		return "", err
	}
	user.ForcePasswordReset = false
	user.LoginAttempts = 0
	err = r.SaveUser(user)
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

func (r *UserRepository) SaveUser(user User) error {
	var testUser User
	record := r.Database.WithContext(context.Background()).Where("username = ?", user.Username).First(&testUser)
	if record.Error != nil && errors.Is(record.Error, gorm.ErrRecordNotFound) {
		record = r.Database.WithContext(context.Background()).Create(&user)
	} else if record.Error == nil {
		record = r.Database.WithContext(context.Background()).Save(&user)
	}
	err := record.Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) LoadUser(username string) (User, error) {
	var user User
	record := r.Database.Where("username = ?", username).First(&user)
	if record.Error != nil {
		return user, record.Error
	}
	return user, nil
}

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Role.Permissions == ADMIN {
		return errors.New("admin users are not allowed to be deleted")
	}
	return
}

func (r *UserRepository) GenerateRandmoPassword() string {
	literalList := "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890_*#^&@:<>.,?+="
	var s string
	for i := 1; i < 16; i++ {
		random := rand.Intn(len(literalList))
		s += string(literalList[random])
	}
	return s
}

func (r *UserRepository) MigrateUserModel() error {
	return UserRepo.Database.AutoMigrate(&User{})
}

func (r *UserRepository) InitUserModel() error {
	err := r.InitRoleModle()
	if err != nil {
		return err
	}
	var admin User
	record := r.Database.Where("username = ?", "admin").First(&admin)
	if record.Error != nil && errors.Is(record.Error, gorm.ErrRecordNotFound) {
		admin, err = r.CreateNewUser("admin", "admin@no.email", "Password_1")
		if err != nil {
			return err
		}
		admin, err = r.LoadUser("admin")
		if err != nil {
			return err
		}
		adminRole, err := r.LoadRole("ADMIN")
		if err != nil {
			return err
		}
		admin.Role = adminRole
		admin.ForcePasswordReset = true
		err = r.SaveUser(admin)
		if err != nil {
			return err
		}
	}
	return nil
}
