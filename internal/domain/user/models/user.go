package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"user-service/internal/shared"
)

type RoleEnum string

const (
	Admin  RoleEnum = "admin"
	Member RoleEnum = "member"
)

type User struct {
	ID string `json:"id" gorm:"type:uuid;primaryKey;not null"`
	shared.BaseModel
	Password      string     `json:"password,omitempty" gorm:"type:varchar(255);not null"`
	FirstName     string     `json:"first_name,omitempty" gorm:"type:varchar(255);not null"`
	LastName      *string    `json:"last_name,omitempty" gorm:"type:varchar(255)"`
	PhoneNumber   *string    `json:"phone_number,omitempty" gorm:"type:varchar(50);unique"`
	Email         string     `json:"email,omitempty" gorm:"type:varchar(255);unique"`
	Role          RoleEnum   `json:"role" gorm:"type:varchar(255)"`
	IsActive      bool       `json:"is_active" gorm:"type:boolean;not null"`
	LastLoginDate *time.Time `json:"last_login_date,omitempty" gorm:"type:timestamptz"`
	LastLoginIP   *string    `json:"last_login_ip,omitempty" gorm:"type:varchar(255)"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String() // Generate UUID
	}
	return
}
