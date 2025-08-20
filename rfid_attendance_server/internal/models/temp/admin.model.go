package models

type Admin struct {
	BaseModel

	Name        string `json:"name" gorm:"not null"`
	Email       string `json:"email" gorm:"uniqueIndex;not null"`
	Password    string `json:"-" gorm:"not null"`
	PhoneNumber string `json:"phone_number"`
	AuthToken   string `json:"auth_token,omitempty"`

	AdminLevel     AdminLevel    `json:"admin_level" gorm:"not null"` // either super_admin or org_admin
	OrganizationID *string       `json:"organization_id"`
	Organization   *Organization `json:"organization" gorm:"foreignKey:OrganizationID;references:ID;constraint:OnDelete:SET NULL"`

	CreatedByID *string `json:"created_by_id"`
	CreatedBy   *Admin  `json:"created_by" gorm:"foreignKey:CreatedByID;references:ID;constraint:OnDelete:SET NULL"`

	UpdatedByID *string `json:"updated_by_id"`
	UpdatedBy   *Admin  `json:"updated_by" gorm:"foreignKey:UpdatedByID;references:ID"`

	DeletedByID *string `json:"deleted_by_id"`
	DeletedBy   *Admin  `json:"deleted_by" gorm:"foreignKey:DeletedByID;references:ID"`
}

type AdminLevel string

const (
	SuperAdmin AdminLevel = "super_admin"
	OrgAdmin   AdminLevel = "org_admin"
)
