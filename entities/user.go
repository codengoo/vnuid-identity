package entities

const (
	Admin   string = "ADMIN"
	Student string = "STUDENT"
	Teacher string = "TEACHER"
)

type User struct {
	ID            string `json:"id" gorm:"primaryKey;type:uuid;"`
	Email         string `json:"email" gorm:"uniqueIndex;not null"`
	Sid           string `json:"sid" gorm:"uniqueIndex;not null"`
	Gid           string `json:"gid" gorm:"uniqueIndex;not null"`
	Password      string `json:"password"`
	Type          string `json:"type" gorm:"not null"`
	Authenticator string `json:"authenticator"`
	BiometricKey  string `json:"biometric_key"`
	ProfileId     string `json:"profile_id" gorm:"uniqueIndex;not null"`

	Sessions []Session `json:"sessions" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Nfcs     []NFC     `json:"nfcs" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	Profile  Profile   `json:"profile" gorm:"foreignKey:ProfileId;references:ID;constraint:OnDelete:CASCADE"`
}
