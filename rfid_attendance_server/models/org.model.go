package models

type Organization struct {
	BaseModel

	Name    string `json:"name" binding:"required" gorm:"unique;not null;index"`
	Email   string `json:"email" gorm:"unique;not null;index"`
	Phone   string `json:"phone"`
	Address string `json:"address"`

	AvailableRoles        map[string]any `gorm:"type:jsonb" json:"available_roles"`
	AvailableShifts       map[string]any `gorm:"type:jsonb" json:"available_shifts"`
	EntryDuplicationDelay int            `json:"entry_duplication_delay" gorm:"default:5"`

	CreatedByID string      `json:"created_by_id" gorm:"not null;index"`
	CreatedBy   SystemAdmin `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`

	DeletedByID string      `json:"deleted_by_id"`
	DeletedBy   SystemAdmin `json:"deleted_by" gorm:"foreignKey:DeletedByID;references:ID"`

	UpdatedBy     string `json:"updated_by"`
	UpdatedByType string `json:"updated_by_type"`
}
