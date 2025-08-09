package converter

import (
	"golectro-product/internal/entity"
	"golectro-product/internal/model"
)

func ToCreateProductResponse(product *entity.Product) *model.CreateProductResponse {
	return &model.CreateProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		Brand:       product.Brand,
		Color:       product.Color,
		Specs:       product.Specs,
		CreatedBy:   product.CreatedBy,
	}
}
