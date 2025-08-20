package models

type ActionLog struct {
	BaseModel

	UserId   string `json:"user_id"`
	UserType string `json:"user_type" `
	Action   string `json:"action" `

	TargetType     string       `json:"target_type" `
	TargetId       string       `json:"target_id"`
	OrganizationId string       `json:"organization_id"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationId;references:ID"`

	IPAddress string `json:"ip_address"`
	Location  string `json:"location"`

	UserAgent string         `json:"user_agent"`
	ExtraData map[string]any `gorm:"type:jsonb" json:"extra_data"`
	Message   string         `json:"message"` //human-readable summary
	IsSuccess bool           `json:"is_success" gorm:"default:true"`
}
