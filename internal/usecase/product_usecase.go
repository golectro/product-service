package usecase

import (
	"context"
	"golectro-product/internal/entity"
	"golectro-product/internal/repository"

	"github.com/go-playground/validator/v10"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	ProductRepository *repository.ProductRepository
}

func NewProductUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, productRepository *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		ProductRepository: productRepository,
	}
}

func (uc *ProductUseCase) GetAllProducts(ctx context.Context, limit, offset int) ([]entity.Product, int64, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	products, total, err := uc.ProductRepository.GetAll(tx, limit, offset)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to get all products")
		return nil, 0, err
	}
	return products, total, nil
}
