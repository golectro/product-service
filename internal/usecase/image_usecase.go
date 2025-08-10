package usecase

import (
	"context"
	"golectro-product/internal/constants"
	"golectro-product/internal/entity"
	"golectro-product/internal/repository"
	"golectro-product/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImageUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	ImageRepository      *repository.ImageRepository
	ElasticsearchUseCase *ElasticsearchUseCase
}

func NewImageUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, productRepository *repository.ImageRepository, elasticsearchUseCase *ElasticsearchUseCase) *ImageUseCase {
	return &ImageUseCase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		ImageRepository:      productRepository,
		ElasticsearchUseCase: elasticsearchUseCase,
	}
}

func (uc *ImageUseCase) GetImageByID(ctx context.Context, imageID uuid.UUID) (*entity.ProductImage, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	image, err := uc.ImageRepository.FindImageById(tx, imageID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find image by ID")
		return nil, utils.WrapMessageAsError(constants.FailedGetImageByID)
	}

	if image == nil {
		return nil, utils.WrapMessageAsError(constants.ImageNotFound)
	}

	return image, nil
}

func (uc *ImageUseCase) DeleteImage(ctx context.Context, imageID *entity.ProductImage) error {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.ImageRepository.Delete(tx, imageID); err != nil {
		uc.Log.WithError(err).Error("Failed to delete image")
		return utils.WrapMessageAsError(constants.FailedDeleteImage, err)
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return utils.WrapMessageAsError(constants.FailedCommitTransaction, err)
	}

	uc.Log.Info("Image deleted successfully")
	return nil
}
