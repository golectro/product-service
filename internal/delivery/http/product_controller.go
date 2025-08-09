package http

import (
	"encoding/json"
	"golectro-product/internal/constants"
	"golectro-product/internal/delivery/http/middleware"
	"golectro-product/internal/model"
	"golectro-product/internal/usecase"
	"golectro-product/internal/utils"
	"math"
	"net/http"
	"strconv"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ProductController struct {
	Log            *logrus.Logger
	ProductUseCase *usecase.ProductUseCase
	MinioUseCase   *usecase.MinioUseCase
	Viper          *viper.Viper
}

func NewProductController(userUseCase *usecase.ProductUseCase, minioUseCase *usecase.MinioUseCase, log *logrus.Logger, viper *viper.Viper) *ProductController {
	return &ProductController{
		Log:            log,
		ProductUseCase: userUseCase,
		MinioUseCase:   minioUseCase,
		Viper:          viper,
	}
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := c.ProductUseCase.GetAllProducts(ctx, limit, offset)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get addresses")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedGetProducts, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := model.PageMetadata{
		CurrentPage: page,
		PageSize:    limit,
		TotalPage:   int64(totalPages),
		TotalItem:   total,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	res := utils.SuccessWithPaginationResponse(ctx, http.StatusOK, constants.SuccessGetProducts, products, pagination)
	ctx.JSON(res.StatusCode, res)
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	var roles []string
	if err := json.Unmarshal(auth.Roles, &roles); err != nil {
		c.Log.WithError(err).Error("Failed to decode roles")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	isAdmin := slices.Contains(roles, "admin")
	if !isAdmin {
		res := utils.FailedResponse(ctx, http.StatusForbidden, constants.AccessDenied, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	request := new(model.CreateProductRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("Failed to bind request")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidRequestData, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	result, err := c.ProductUseCase.CreateProduct(ctx, request, auth.ID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to create product")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedCreateProduct, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusCreated, constants.SuccessCreateProduct, result)
	ctx.JSON(res.StatusCode, res)
}
