package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type (
	CreateProductRequest struct {
		Name        string         `json:"name" validate:"required,max=255"`
		Description string         `json:"description" validate:"max=2000"`
		Category    datatypes.JSON `json:"category"`
		Brand       string         `json:"brand" validate:"required,max=100"`
		Color       datatypes.JSON `json:"color"`
		Specs       datatypes.JSON `json:"specs"`
		Price       float64        `json:"price" validate:"required"`
	}

	CreateProductResponse struct {
		ID          uuid.UUID      `json:"id"`
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Category    datatypes.JSON `json:"category"`
		Brand       string         `json:"brand"`
		Color       datatypes.JSON `json:"color"`
		Specs       datatypes.JSON `json:"specs"`
		Price       float64        `json:"price"`
		CreatedBy   uuid.UUID      `json:"created_by"`
	}

	UpdateProductRequest struct {
		Name        *string         `json:"name,omitempty" validate:"max=255"`
		Description *string         `json:"description,omitempty" validate:"max=2000"`
		Category    *datatypes.JSON `json:"category,omitempty"`
		Brand       *string         `json:"brand,omitempty" validate:"max=100"`
		Color       *datatypes.JSON `json:"color,omitempty"`
		Specs       *datatypes.JSON `json:"specs,omitempty"`
		Price       *float64        `json:"price,omitempty"`
	}

	UploadFilesResponse struct {
		ProductID uuid.UUID `json:"product_id"`
		Images    []string  `json:"images"`
	}

	ProductImageURLResponse struct {
		ID          uuid.UUID `json:"id"`
		ProductID   uuid.UUID `json:"product_id"`
		ImageObject string    `json:"image_object"`
		URL         string    `json:"url"`
	}
)
