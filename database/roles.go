package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          string `gorm:"primaryKey"`
	RoleName    string `gorm:"not null,type:text"`
	Permissions permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID version 4
	role.ID = uuid.NewString()
	tx.Statement.SetColumn("ID", role.ID)
	return
}

func (r *UserRepository) SaveRole(role Role) error {
	record := r.Database.WithContext(context.Background()).Where("role_name = ?", role.RoleName).First(&role)
	var err error = nil
	if record.Error != nil && errors.Is(record.Error, gorm.ErrRecordNotFound) {
		log.Default().Println("Creating new Role")
		recordC := r.Database.WithContext(context.Background()).Create(&role)
		err = recordC.Error
		if err != nil {
			log.Default().Println("err: ", err)
		}
		log.Default().Println("Role created...")
	} else if record.Error == nil {
		log.Default().Println("Saving existing Role: " + role.RoleName)
		recordS := r.Database.WithContext(context.Background()).Save(&role)
		err = recordS.Error
	}
	if err != nil {
		log.Default().Println("Error during save operation")
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

func (r *UserRepository) LoadRole(name string) (Role, error) {
	var role Role
	record := r.Database.Where("role_name = ?", name).First(&role)
	if record.Error != nil {
		return role, record.Error
	}
	return role, nil
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
	no_perm.Permissions = NO_PERMISSIONS
	err := r.SaveRole(no_perm)
	if err != nil {
		return err
	}
	var admin Role
	admin.RoleName = "ADMIN"
	admin.Permissions = ADMIN
	err = r.SaveRole(admin)
	if err != nil {
		return err
	}
	return nil
}
