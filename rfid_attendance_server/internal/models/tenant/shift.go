package tenant

type Shift struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	StartTime string `json:"start_time"` // "09:00"
	EndTime   string `json:"end_time"`   // "17:00"
}
