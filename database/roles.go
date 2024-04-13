package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID `gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"`
	RoleName    string    `gorm:"not null,type:text"`
	Permissions permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (r *Role) Save() error {
	record := UserRepo.Database.Where("role_name = ?", r.RoleName).First(&r)
	if record.Error != nil {
		record = UserRepo.Database.Create(&r)
	} else {
		record = UserRepo.Database.Save(&r)
	}
	err := record.Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) SetPermission(flag permission) {
	r.Permissions = r.Permissions | flag
}

func (r *Role) UnsetPermission(flag permission) {
	r.Permissions = r.Permissions & ^flag
}

func (r *Role) Load(name string) error {
	record := UserRepo.Database.Where("role_name = ?", name).First(&r)
	if record.Error != nil {
		return record.Error
	}
	return nil
}

func (r *Role) BeforeDelete(tx *gorm.DB) (err error) {
	if r.Permissions == ADMIN {
		return errors.New("admin role can not be deleted")
	}
	return
}

func (r *UserRepository) MigrateRoleModel() error {
	return UserRepo.Database.AutoMigrate(&Role{})
}

func (r *UserRepository) InitRoleModle() error {
	var no_perm Role
	no_perm.RoleName = "NO_PERMISSIONS"
	no_perm.Permissions = NO_PERMISSION
	err := no_perm.Save()
	if err != nil {
		return err
	}
	var admin Role
	admin.RoleName = "ADMIN"
	admin.Permissions = ADMIN
	err = admin.Save()
	if err != nil {
		return err
	}
	return nil
}