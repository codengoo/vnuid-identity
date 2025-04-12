package entities

import "time"

type NFC struct {
	ID        string `json:"id" gorm:"primaryKey;type:uuid;"`
	UserId    string `json:"user_id" gorm:"not null;"`
	Active    bool   `json:"active" gorm:"not null;"`
	Status    string `json:"status" gorm:"not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
