package models

type OrganizationAdmin struct {
	BaseModel
	OrganizationID string       `json:"organization_id" gorm:"not null;index"`
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationID;references:ID"`
	Name           string       `json:"name"`
	Email          string       `json:"email" gorm:"unique;not null;index"`
	Password       string       `json:"password"`
	AuthToken      string       `json:"auth_token"`
	PhoneNumber    string       `json:"phone_number"`
	CreatedByID    string       `json:"created_by_id" gorm:"not null;index"`
	CreatedBy      SystemAdmin  `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID"`
	DeletedByID    string       `json:"deleted_by_id"`
	DeletedBy      SystemAdmin  `json:"deleted_by" gorm:"foreignKey:DeletedByID;references:ID"`
	UpdatedByID    string       `json:"updated_by_id"`
	UpdatedBy      SystemAdmin  `json:"updated_by" gorm:"foreignKey:UpdatedByID;references:ID"`
}
