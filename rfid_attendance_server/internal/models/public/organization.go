package public

import "time"

type Organization struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Schema    string    `gorm:"uniqueIndex;not null" json:"schema"` // e.g. tenant_acme
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
