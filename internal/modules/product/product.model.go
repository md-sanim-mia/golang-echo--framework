package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"type:varchar(150);not null"`
	Slug        string    `gorm:"uniqueIndex;type:varchar(150)"`
	Description string    `gorm:"type:text"`
	Price       float64   `gorm:"not null"`
	Stock       int       `gorm:"default:0"`
	Category    string    `gorm:"type:varchar(100)"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   gorm.DeletedAt
	UpdatedAt   gorm.DeletedAt
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}
