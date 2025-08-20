package tenant

import "time"

type Member struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name       string `json:"name"`
	Email      string `gorm:"uniqueIndex;not null" json:"email"`
	Phone      string `json:"phone"`
	Role       string `json:"role"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
	RFIDUID    string `gorm:"uniqueIndex;not null" json:"rfid_uid"`
	Department string `json:"department"`
	Status     string `gorm:"default:active" json:"status"`

	ExpectedCheckIn  string `json:"expected_checkin"`  // "09:00"
	ExpectedCheckOut string `json:"expected_checkout"` // "17:00"

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
