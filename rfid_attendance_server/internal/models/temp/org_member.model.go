package models

type OrganizationMember struct {
	BaseModel
	OrganizationID string       `json:"organization_id" gorm:"not null;index"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID;references:ID"`
	Name           string       `json:"name"`
	Email          string       `json:"email" gorm:"unique;not null;index"`
	PhoneNumber    string       `json:"phone_number"`
	Role           string       `json:"role"`
	IsActive       bool         `json:"is_active" gorm:"default:true"`
	RFIDUID        string       `json:"rfid_uid" gorm:"unique;not null;index"`
	Department     string       `json:"department"`
	Status         string       `gorm:"default:active" json:"status"` //active, inactive, disabled, deleted
	Shift          string       `json:"shift"`

	Expected_checkin_time  string `json:"expected_checkin_time"`
	Expected_checkout_time string `json:"expected_checkout_time"`

	CreatedByID string `json:"created_by_id" gorm:"not null;index"`
	CreatedBy   Admin  `json:"created_by_type" gorm:"foreignKey:CreatedByID;references:ID;constraint:OnDelete:SET NULL"`
	UpdatedByID string `json:"updated_by_id"`
	UpdatedBy   Admin  `json:"updated_by_type" gorm:"foreignKey:UpdatedByID;references:ID;constraint:OnDelete:SET NULL"`
	DeletedByID string `json:"deleted_by_id"`
	DeletedBy   Admin  `json:"deleted_by_type" gorm:"foreignKey:DeletedByID;references:ID;constraint:OnDelete:SET NULL"`
}

type Status string

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	Disabled Status = "disabled"
	Deleted  Status = "deleted"
)
