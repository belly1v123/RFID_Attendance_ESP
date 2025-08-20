package tenant

import "time"

type ScanLog struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RFIDUID    string    `gorm:"type:varchar(32);not null" json:"rfid_uid"`
	DeviceID   string    `gorm:"index" json:"device_id"`
	ScannedAt  time.Time `gorm:"autoCreateTime" json:"scanned_at"`
	Recognized bool      `gorm:"default:false" json:"recognized"`
}
