package product

import (
	"errors"

	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func (s *ProductService) CreateProduct(payload *Product) error {

	if s.DB == nil {

		return errors.New("database not initialized")
	}
	return s.DB.Create(payload).Error
}
