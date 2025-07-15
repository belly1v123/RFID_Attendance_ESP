package models

import "time"

type AttendanceRecord struct {
	ID       uint       `gorm:"primaryKey" json:"id"`
	UserID   uint       `gorm:"not null" json:"user_id"`
	CheckIn  *time.Time `json:"check_in,omitempty"`
	CheckOut *time.Time `json:"check_out,omitempty"`
	Day      time.Time  `gorm:"type:date;not null" json:"day"`
	Status   string     `gorm:"type:varchar(20);not null" json:"status"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
