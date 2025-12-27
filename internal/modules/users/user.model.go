package users

import (
	"encoding/json"
	"time"
)

type UserRole int

const (
	USER UserRole = iota
	ADMIN
)

func (r UserRole) String() string {
	switch r {
	case USER:
		return "USER"
	case ADMIN:
		return "ADMIN"
	default:
		return "USER"
	}
}

func (r UserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FullName  string    `gorm:"size:100;not null" json:"fullName"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `json:"password" validate:"required,min=6"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Address   string    `gorm:"size:255" json:"address"`
	Age       int       `json:"age"`
	Role      UserRole  `gorm:"default:0"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
