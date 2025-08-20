package public

import "time"

type ActionLog struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserId         string         `json:"user_id"`
	UserType       string         `json:"user_type"`
	Action         string         `json:"action"`
	TargetType     string         `json:"target_type"`
	TargetId       string         `json:"target_id"`
	OrganizationID string         `json:"organization_id"`
	IPAddress      string         `json:"ip_address"`
	Location       string         `json:"location"`
	UserAgent      string         `json:"user_agent"`
	ExtraData      map[string]any `gorm:"type:jsonb" json:"extra_data"`
	Message        string         `json:"message"`
	IsSuccess      bool           `gorm:"default:true" json:"is_success"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
}
