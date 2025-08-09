package repository

import (
	"golectro-product/internal/entity"

	"github.com/google/uuid"
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

	err := db.Preload("Images").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	if err != nil {
		r.Log.WithError(err).Error("Failed to find products with images and pagination")
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) FindProductById(db *gorm.DB, productID uuid.UUID) (*entity.Product, error) {
	var product entity.Product

	if err := db.Preload("Images").First(&product, "id = ?", productID).Error; err != nil {
		r.Log.WithError(err).Error("Failed to find product by ID")
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) CreateImage(db *gorm.DB, productImage *entity.ProductImage) error {
	return db.Create(productImage).Error
}
