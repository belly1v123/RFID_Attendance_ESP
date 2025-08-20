package models

import "time"

type AttendanceRecord struct {
	BaseModel

	OrganizationID string       `json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"organization"`

	UserID string             `json:"user_id"`
	User   OrganizationMember `gorm:"foreignKey:UserID;references:ID" json:"user"`

	CheckIn  *time.Time `json:"check_in,omitempty"`
	CheckOut *time.Time `json:"check_out,omitempty"`

	PresentDay time.Time `gorm:"type:date;not null" json:"day"`

	Status string `gorm:"not null;default:absent" json:"status"`

	UpdatedBy     string `json:"updated_by"`
	UpdatedByType string `json:"updated_by_type"`

	// User User `gorm:"foreignKey:UserID" json:"-"`
}
