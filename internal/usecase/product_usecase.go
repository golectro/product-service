package usecase

import (
	"context"
	"golectro-product/internal/constants"
	"golectro-product/internal/entity"
	"golectro-product/internal/model"
	"golectro-product/internal/model/converter"
	"golectro-product/internal/repository"
	"golectro-product/internal/utils"
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

func (uc *ProductUseCase) GetProductByID(ctx context.Context, productID uuid.UUID) (*model.CreateProductResponse, error) {
	product, err := uc.ProductRepository.FindProductById(uc.DB.WithContext(ctx), productID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find product by ID")
		return nil, utils.WrapMessageAsError(constants.FailedGetProductByID)
	}

	if product == nil {
		return nil, utils.WrapMessageAsError(constants.ProductNotFound)
	}

	return converter.ToProductResponse(product), nil
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, request *model.CreateProductRequest, userID uuid.UUID) (*model.CreateProductResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		message := utils.TranslateValidationError(uc.Validate, err)
		return nil, utils.WrapMessageAsError(message)
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
		return nil, utils.WrapMessageAsError(constants.FailedToCreateProduct, err)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return nil, utils.WrapMessageAsError(constants.FailedCreateProduct, err)
	}

	return converter.ToProductResponse(entityProduct), nil
}

func (uc *ProductUseCase) UploadProductImages(ctx context.Context, productID uuid.UUID, images []map[string]any) (*model.UploadFilesResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product, err := uc.ProductRepository.FindProductById(tx, productID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find product by ID")
		return nil, utils.WrapMessageAsError(constants.FailedGetProductByID)
	}

	var uploadedImages []entity.ProductImage

	for _, img := range images {
		fileName, ok := img["file_name"].(string)
		if !ok || fileName == "" {
			continue
		}

		productImage := entity.ProductImage{
			ID:          uuid.New(),
			ProductID:   product.ID,
			ImageObject: fileName,
		}

		if err := uc.ProductRepository.CreateImage(tx, &productImage); err != nil {
			uc.Log.WithError(err).Error("Failed to create product image record")
			tx.Rollback()
			return nil, utils.WrapMessageAsError(constants.FailedCreateProduct, err)
		}

		uploadedImages = append(uploadedImages, productImage)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction for product image upload")
		return nil, utils.WrapMessageAsError(constants.FailedCreateProduct, err)
	}

	return &model.UploadFilesResponse{
		ProductID: productID,
		Images:    extractImageURLs(images),
	}, nil
}

func extractImageURLs(images []map[string]any) []string {
	var urls []string
	for _, img := range images {
		if url, ok := img["file_name"].(string); ok && url != "" {
			urls = append(urls, url)
		}
	}
	return urls
}
