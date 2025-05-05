package entities

type NFC struct {
	ID     string `json:"id" gorm:"primaryKey;type:uuid;"`
	UserId string `json:"user_id" gorm:"not null;"`
	Active bool   `json:"active" gorm:"not null;"`
	Status string `json:"status" gorm:"not null;"`

	User User `json:"user" gorm:""`
}
