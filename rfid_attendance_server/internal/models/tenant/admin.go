package tenant

import "time"

type AdminLevel string

const (
	SuperAdmin AdminLevel = "super_admin"
	OrgAdmin   AdminLevel = "org_admin"
)

type Admin struct {
	ID       string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name     string     `gorm:"not null" json:"name"`
	Email    string     `gorm:"uniqueIndex;not null" json:"email"`
	Password string     `gorm:"not null" json:"-"`
	Level    AdminLevel `gorm:"not null" json:"level"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
