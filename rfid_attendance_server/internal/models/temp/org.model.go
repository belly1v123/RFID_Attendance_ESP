package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

type JSONBMap map[string]any

func (j JSONBMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONBMap) Scan(src any) error {
	return json.Unmarshal(src.([]byte), j)
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if o.AvailableRoles == nil {
		o.AvailableRoles = &JSONBMap{}
	}
	if o.AvailableShifts == nil {
		o.AvailableShifts = &JSONBMap{}
	}
	return nil
}

type Organization struct {
	BaseModel

	Name    string `json:"name" binding:"required" gorm:"uniqueIndex;not null;index"`
	Email   string `json:"email" gorm:"uniqueIndex;not null;index"`
	Phone   string `json:"phone" gorm:"default:null"`
	Address string `json:"address"`

	AvailableRoles  *JSONBMap `gorm:"type:jsonb;" json:"available_roles"`
	AvailableShifts *JSONBMap `gorm:"type:jsonb;" json:"available_shifts"`

	EntryDuplicationDelay int `json:"entry_duplication_delay" gorm:"default:5"`

	CreatedByID *string `json:"created_by_id" gorm:"not null;index"`
	CreatedBy   *Admin  `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID; constraint:OnDelete:SET NULL"`

	DeletedByID *string `json:"deleted_by_id"`
	DeletedBy   *Admin  `json:"deleted_by" gorm:"foreignKey:DeletedByID;references:ID; constraint:OnDelete:SET NULL"`

	UpdatedByID *string `json:"updated_by_id"`
	UpdatedBy   *Admin  `json:"updated_by" gorm:"foreignKey:UpdatedByID;references:ID; constraint:OnDelete:SET NULL"`
}
