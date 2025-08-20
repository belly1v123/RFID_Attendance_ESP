package tenant

import "time"

type AttendanceRecord struct {
	ID       string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MemberID string     `gorm:"index;not null" json:"member_id"`
	CheckIn  *time.Time `json:"check_in,omitempty"`
	CheckOut *time.Time `json:"check_out,omitempty"`
	Day      time.Time  `gorm:"type:date;not null" json:"day"`
	Status   string     `gorm:"not null;default:absent" json:"status"`

	UpdatedBy     string `json:"updated_by"`
	UpdatedByType string `json:"updated_by_type"`
}
