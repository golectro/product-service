package handler

import (
	"context"
	"fmt"
	"golectro-product/internal/constants"
	proto "golectro-product/internal/delivery/grpc/proto/product"
	"golectro-product/internal/usecase"
	"golectro-product/internal/utils"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	proto.UnimplementedProductServiceServer
	ProductUseCase *usecase.ProductUseCase
}

func (h *ProductHandler) GetProductById(ctx context.Context, req *proto.GetProductByIdRequest) (*proto.GetProductByIdResponse, error) {
	productID, err := utils.ParseUUID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	product, err := h.ProductUseCase.GetProductByID(ctx, productID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s: %v", constants.FailedGetProductByID, err)
	}

	if product == nil {
		return nil, status.Errorf(codes.NotFound, "%s: product with ID %s not found", constants.ProductNotFound, req.Id)
	}

	return &proto.GetProductByIdResponse{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Category:    string(product.Category),
		Brand:       product.Brand,
		Color:       string(product.Color),
		Specs:       string(product.Specs),
		Price:       product.Price,
		Quantity:    int32(product.Quantity),
		CreatedBy:   product.CreatedBy.String(),
	}, nil
}

func (h *ProductHandler) GetProductByIds(ctx context.Context, req *proto.GetProductByIdsRequest) (*proto.GetProductByIdsResponse, error) {
	productUUIDs := make([]uuid.UUID, len(req.Ids))
	for i, id := range req.Ids {
		parsedID, err := utils.ParseUUID(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid product ID at index %d: %v", i, err)
		}
		productUUIDs[i] = parsedID
	}

	products, err := h.ProductUseCase.GetProductsByIDs(ctx, productUUIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s: %v", constants.FailedGetProductsByIDs, err)
	}

	if len(products) == 0 {
		return &proto.GetProductByIdsResponse{}, nil
	}

	response := &proto.GetProductByIdsResponse{}
	for _, product := range products {
		response.Products = append(response.Products, &proto.GetProductByIdResponse{
			Id:          product.ID.String(),
			Name:        product.Name,
			Description: product.Description,
			Category:    string(product.Category),
			Brand:       product.Brand,
			Color:       string(product.Color),
			Specs:       string(product.Specs),
			Price:       product.Price,
			Quantity:    int32(product.Quantity),
			CreatedBy:   product.CreatedBy.String(),
		})
	}

	return response, nil
}

func (h *ProductHandler) DecreaseQuantity(ctx context.Context, req *proto.DecreaseQuantityRequest) (*proto.DecreaseQuantityResponse, error) {
	productID, err := utils.ParseUUID(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	quantity, err := h.ProductUseCase.DecreaseProductQuantity(ctx, productID, int(req.Quantity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s: %v", constants.FailedDecreaseProductQuantity, err)
	}

	return &proto.DecreaseQuantityResponse{
		Success:     true,
		Message:     "Product quantity decreased successfully",
		NewQuantity: int32(quantity),
	}, nil
}

func (h *ProductHandler) DecreaseQuantityByIds(ctx context.Context, req *proto.DecreaseQuantityByIdsRequest) (*proto.DecreaseQuantityByIdsResponse, error) {
	if len(req.Items) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no product items provided")
	}

	var results []*proto.DecreaseQuantityResult
	allSuccess := true

	for i, item := range req.Items {
		productID, err := utils.ParseUUID(item.ProductId)
		if err != nil {
			allSuccess = false
			results = append(results, &proto.DecreaseQuantityResult{
				ProductId: item.ProductId,
				Success:   false,
				Message:   fmt.Sprintf("invalid product ID at index %d: %v", i, err),
			})
			continue
		}

		newQty, err := h.ProductUseCase.DecreaseProductQuantity(ctx, productID, int(item.Quantity))
		if err != nil {
			allSuccess = false
			results = append(results, &proto.DecreaseQuantityResult{
				ProductId: item.ProductId,
				Success:   false,
				Message:   fmt.Sprintf("failed to decrease quantity: %v", err),
			})
			continue
		}

		results = append(results, &proto.DecreaseQuantityResult{
			ProductId:   item.ProductId,
			Success:     true,
			NewQuantity: int32(newQty),
			Message:     "quantity decreased successfully",
		})
	}

	return &proto.DecreaseQuantityByIdsResponse{
		Success: allSuccess,
		Message: func() string {
			if allSuccess {
				return "All quantities decreased successfully"
			}
			return "Some quantities failed to decrease"
		}(),
		Results: results,
	}, nil
}
