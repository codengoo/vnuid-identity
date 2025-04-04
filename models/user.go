package models

import "time"

const (
	Admin   string = "ADMIN"
	Student string = "STUDENT"
	Teacher string = "TEACHER"
)

type User struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid;"`
	Email         string    `json:"email" gorm:"uniqueIndex;not null"`
	SID           string    `json:"sid" gorm:"uniqueIndex;not null"`
	GID           string    `json:"gid" gorm:"uniqueIndex;not null"`
	Password      string    `json:"password"`
	Name          string    `json:"name" gorm:"not null"`
	DOB           *string   `json:"dob"`
	OfficialClass string    `json:"official_class" gorm:"not null"`
	Type          string    `json:"type" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
