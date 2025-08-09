package repository

import (
	"golectro-product/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepository struct {
	Repository[entity.Product]
	Log *logrus.Logger
}

func NewProductRepository(log *logrus.Logger) *ProductRepository {
	return &ProductRepository{Log: log}
}

func (r *ProductRepository) GetAll(db *gorm.DB, limit, offset int) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	if err := db.Model(&entity.Product{}).Count(&total).Error; err != nil {
		r.Log.WithError(err).Error("Failed to count products")
		return nil, 0, err
	}

	err := r.Repository.FindByConditionWithPagination(db, &products, "", offset/limit+1, limit)
	if err != nil {
		r.Log.WithError(err).Error("Failed to find products with pagination")
		return nil, 0, err
	}

	return products, total, nil
}
