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
