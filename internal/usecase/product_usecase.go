package usecase

import (
	"context"
	"golectro-product/internal/entity"
	"golectro-product/internal/model"
	"golectro-product/internal/model/converter"
	"golectro-product/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

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

func (uc *ProductUseCase) CreateProduct(ctx context.Context, request *model.CreateProductRequest, userID uuid.UUID) (*model.CreateProductResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation failed for product")
		return nil, err
	}

	entityProduct := &entity.Product{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		Category:    request.Category,
		Brand:       request.Brand,
		Color:       request.Color,
		Specs:       request.Specs,
		Price:       request.Price,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(entityProduct).Error; err != nil {
		uc.Log.WithError(err).Error("Failed to create product")
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return converter.ToCreateProductResponse(entityProduct), nil
}
