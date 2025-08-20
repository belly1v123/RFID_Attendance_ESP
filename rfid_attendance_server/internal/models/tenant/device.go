package tenant

type Device struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
