package migrations

import (
	"golectro-product/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.Product{}, &entity.ProductImage{})
}
