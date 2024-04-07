package models

import (
	"blue-beetle/database"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	RoleName    string `gorm:"primaryKey"`
	Permissions permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (r *Role) Save() error {
	record := database.Instance.Where("role_name = ?", r.RoleName).First(&r)
	if record.Error != nil {
		record = database.Instance.Create(&r)
	} else {
		record = database.Instance.Save(&r)
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
	record := database.Instance.Where("role_name = ?", name).First(&r)
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

func MigrateRoleModel() error {
	return database.Instance.AutoMigrate(&Role{})
}

func InitRoleModle() error {
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