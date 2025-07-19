package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string         `gorm:"index;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SystemAdmin struct {
	BaseModel
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	AuthToken   string `json:"auth_token"`
}
