package models

import "time"

type Log struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RFIDUID    string    `gorm:"type:varchar(32);not null" json:"rfid_uid"`
	ScannedAt  time.Time `gorm:"autoCreateTime" json:"scanned_at"`
	Recognized bool      `gorm:"default:false" json:"recognized"`
}
