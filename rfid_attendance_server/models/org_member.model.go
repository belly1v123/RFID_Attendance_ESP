package models

import "time"

// type User struct {
// 	ID         uint      `gorm:"primaryKey" json:"id"`
// 	Name       string    `json:"name"`
// 	RFIDUID    string    `gorm:"unique;not null" json:"rfid_uid"`
// 	Department string    `json:"department"`
// 	Status     string    `gorm:"default:active" json:"status"`
// 	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
// 	Email      string    `json:"email"`
// 	Phone      string    `json:"phone"`

// 	AttendanceRecords []AttendanceRecord `gorm:"foreignKey:UserID" json:"attendance_records,omitempty"`
// }

type OrganizationMember struct {
	BaseModel
	OrganizationID string       `json:"organization_id" gorm:"not null;index"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID;references:ID"`
	Name           string       `json:"name"`
	Email          string       `json:"email" gorm:"unique;not null;index"`
	PhoneNumber    string       `json:"phone_number"`
	Role           string       `json:"role"`
	IsActive       bool         `json:"is_active" gorm:"default:true"`
	RFIDUID        string       `json:"rfid_uid" gorm:"unique;not null;index"`
	Department     string       `json:"department"`
	Status         string       `gorm:"default:active" json:"status"` //active, inactive, disabled, deleted
	Shift          string       `json:"shift"`

	Expected_checkin_time  time.Time `json:"expected_checkin_time"`
	Expected_checkout_time time.Time `json:"expected_checkout_time"`

	CreatedByID   string `json:"created_by_id" gorm:"not null;index"`
	CreatedByType string `json:"created_by_type"`
	UpdatedByID   string `json:"updated_by_id"`
	UpdatedByType string `json:"updated_by_type"`
	DeletedByID   string `json:"deleted_by_id"`
	DeletedByType string `json:"deleted_by_type"`
}
