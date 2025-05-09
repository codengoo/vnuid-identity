package entities

type Profile struct {
	ID            string  `json:"id" gorm:"primaryKey;type:uuid;"`
	Sid           string  `json:"sid" gorm:"not null"`
	Email         string  `json:"email" gorm:"not null"`
	DOB           *string `json:"dob"`
	OfficialClass string  `json:"official_class" gorm:"not null"`
	Name          string  `json:"name" gorm:"not null"`
	Phone         string  `json:"phone"`
	Address       string  `json:"address"`
}
