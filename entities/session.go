package entities

import "time"

type Session struct {
	ID          string `json:"id" gorm:"primaryKey;type:uuid;"`
	DeviceId    string `json:"device_id" gorm:"not null"`
	UserId      string `json:"user_id" gorm:"not null"`
	LoginMethod string `json:"login_method" gorm:"not null"`
	SavedDevice bool   `json:"saved" gorm:"not null"`
	// IpAddress   string    `json:"ip_address" gorm:"not null"`
	// Location    string    `json:"location" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
