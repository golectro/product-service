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
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	ProductRepository      *repository.ProductRepository
	ProductImageRepository *repository.ImageRepository
	ElasticsearchUseCase   *ElasticsearchUseCase
}

func NewProductUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, productRepository *repository.ProductRepository, productImageRepository *repository.ImageRepository, elasticsearchUseCase *ElasticsearchUseCase) *ProductUseCase {
	return &ProductUseCase{
		DB:                     db,
		Log:                    log,
		Validate:               validate,
		ProductRepository:      productRepository,
		ProductImageRepository: productImageRepository,
		ElasticsearchUseCase:   elasticsearchUseCase,
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

func (uc *ProductUseCase) GetProductByID(ctx context.Context, productID uuid.UUID) (*model.ProductResponse, error) {
	product, err := uc.ProductRepository.FindProductById(uc.DB.WithContext(ctx), productID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find product by ID")
		return nil, utils.WrapMessageAsError(constants.FailedGetProductByID, err)
	}

	if product == nil {
		return nil, utils.WrapMessageAsError(constants.ProductNotFound)
	}

	return converter.ToProductResponse(product), nil
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product, err := uc.ProductRepository.FindProductById(tx, productID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find product by ID")
		return nil, utils.WrapMessageAsError(constants.FailedGetProductByID, err)
	}

	if product == nil {
		return nil, utils.WrapMessageAsError(constants.ProductNotFound)
	}

	return product, nil
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, request *model.ProductRequest, userID uuid.UUID) (*model.ProductResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		message := utils.TranslateValidationError(uc.Validate, err)
		return nil, utils.WrapMessageAsError(message)
	}

	productID := uuid.New()
	entityProduct := &entity.Product{
		ID:          productID,
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

	if err := uc.ElasticsearchUseCase.InsertDocument(productID, entityProduct); err != nil {
		uc.Log.WithError(err).Error("Failed to insert product into Elasticsearch")
		return nil, utils.WrapMessageAsError(constants.FailedInsertProductToElasticsearch, err)
	}

	return converter.ToProductResponse(entityProduct), nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, productID uuid.UUID, request *model.UpdateProductRequest) (*model.ProductResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product, err := uc.ProductRepository.FindProductById(tx, productID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find product by ID")
		return nil, utils.WrapMessageAsError(constants.FailedGetProductByID)
	}

	if product == nil {
		return nil, utils.WrapMessageAsError(constants.ProductNotFound)
	}

	if err := uc.Validate.Struct(request); err != nil {
		message := utils.TranslateValidationError(uc.Validate, err)
		return nil, utils.WrapMessageAsError(message)
	}

	if request.Name != nil {
		product.Name = *request.Name
	}
	if request.Description != nil {
		product.Description = *request.Description
	}
	if request.Category != nil {
		product.Category = *request.Category
	}
	if request.Brand != nil {
		product.Brand = *request.Brand
	}
	if request.Color != nil {
		product.Color = *request.Color
	}
	if request.Specs != nil {
		product.Specs = *request.Specs
	}
	if request.Price != nil {
		product.Price = *request.Price
	}
	product.UpdatedAt = time.Now()

	if err := tx.Save(product).Error; err != nil {
		uc.Log.WithError(err).Error("Failed to update product")
		return nil, utils.WrapMessageAsError(constants.FailedUpdateProduct, err)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction for product update")
		return nil, utils.WrapMessageAsError(constants.FailedUpdateProduct, err)
	}

	if err := uc.ElasticsearchUseCase.InsertDocument(product.ID, product); err != nil {
		uc.Log.WithError(err).Error("Failed to update product in Elasticsearch")
		return nil, utils.WrapMessageAsError(constants.FailedUpdateProductInElasticsearch, err)
	}

	return converter.ToProductResponse(product), nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, product *entity.Product) error {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.ProductRepository.Delete(tx, product); err != nil {
		uc.Log.WithError(err).Error("Failed to delete product")
		return utils.WrapMessageAsError(constants.FailedDeleteProduct, err)
	}

	if err := uc.ElasticsearchUseCase.DeleteDocumentByID(product.ID.String()); err != nil {
		uc.Log.WithError(err).Error("Failed to delete product from Elasticsearch")
		return utils.WrapMessageAsError(constants.FailedDeleteProductFromElasticsearch, err)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction for product deletion")
		return utils.WrapMessageAsError(constants.FailedDeleteProduct, err)
	}

	return nil
}

func (uc *ProductUseCase) UploadProductImages(ctx context.Context, productID uuid.UUID, images []map[string]any) (*model.UploadFilesResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	product, err := uc.ProductRepository.FindProductById(tx, productID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find product by ID")
		return nil, utils.WrapMessageAsError(constants.FailedGetProductByID)
	}

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
