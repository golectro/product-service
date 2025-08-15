package converter

import (
	"golectro-product/internal/entity"
	"golectro-product/internal/model"
)

func ToProductResponse(product *entity.Product) *model.ProductResponse {
	return &model.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		Brand:       product.Brand,
		Color:       product.Color,
		Specs:       product.Specs,
		Quantity:    product.Quantity,
		CreatedBy:   product.CreatedBy,
	}
}
