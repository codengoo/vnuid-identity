package entities

type Session struct {
	ID          string `json:"id" gorm:"primaryKey;type:uuid;"`
	DeviceId    string `json:"device_id" gorm:"not null"`
	DeviceName  string `json:"device_name" gorm:"not null"`
	UserId      string `json:"user_id" gorm:""`
	LoginMethod string `json:"login_method" gorm:"not null"`
	SavedDevice bool   `json:"saved_device" gorm:"not null"`
	// IpAddress   string    `json:"ip_address" gorm:"not null"`
	// Location    string    `json:"location" gorm:"not null"`

	User User `json:"user" gorm:""`
}
