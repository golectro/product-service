package repository

import (
	"golectro-product/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImageRepository struct {
	Repository[entity.ProductImage]
	Log *logrus.Logger
}

func NewImageRepository(log *logrus.Logger) *ImageRepository {
	return &ImageRepository{Log: log}
}

func (r *ImageRepository) FindImageById(db *gorm.DB, imageID uuid.UUID) (*entity.ProductImage, error) {
	var image entity.ProductImage

	if err := db.First(&image, "id = ?", imageID).Error; err != nil {
		r.Log.WithError(err).Error("Failed to find image by ID")
		return nil, err
	}

	return &image, nil
}
