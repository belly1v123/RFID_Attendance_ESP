package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	RFIDUID    string    `gorm:"unique;not null" json:"rfid_uid"`
	Department string    `json:"department"`
	Status     string    `gorm:"default:active" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
