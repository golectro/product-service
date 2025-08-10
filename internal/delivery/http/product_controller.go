package http

import (
	"encoding/json"
	"fmt"
	"golectro-product/internal/constants"
	"golectro-product/internal/delivery/http/middleware"
	"golectro-product/internal/model"
	"golectro-product/internal/usecase"
	"golectro-product/internal/utils"
	"io"
	"math"
	"net/http"
	"path"
	"strconv"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ProductController struct {
	Log            *logrus.Logger
	ProductUseCase *usecase.ProductUseCase
	ImageUseCase   *usecase.ImageUseCase
	MinioUseCase   *usecase.MinioUseCase
	Viper          *viper.Viper
}

func NewProductController(userUseCase *usecase.ProductUseCase, minioUseCase *usecase.MinioUseCase, log *logrus.Logger, viper *viper.Viper, imageUseCase *usecase.ImageUseCase) *ProductController {
	return &ProductController{
		Log:            log,
		ProductUseCase: userUseCase,
		ImageUseCase:   imageUseCase,
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

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	productID := ctx.Param("productID")
	if productID == "" {
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductID, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		c.Log.WithError(err).Error("Invalid product ID format")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductIDFormat, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	product, err := c.ProductUseCase.GetProductByID(ctx, productUUID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get product by ID")
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.FailedGetProductByID, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.SuccessGetProductByID, product)
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

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
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

	productID := ctx.Param("productID")
	if productID == "" {
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductID, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		c.Log.WithError(err).Error("Invalid product ID format")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductIDFormat, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	request := new(model.UpdateProductRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("Failed to bind request")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidRequestData, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	result, err := c.ProductUseCase.UpdateProduct(ctx, productUUID, request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to update product")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedUpdateProduct, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.SuccessUpdateProduct, result)
	ctx.JSON(res.StatusCode, res)
}

func (c *ProductController) UploadProductImages(ctx *gin.Context) {
	auth := middleware.GetUser(ctx)

	productID := ctx.Param("productID")
	if productID == "" {
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductID, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		c.Log.WithError(err).Error("Invalid product ID format")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductIDFormat, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

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

	uploadedFilesAny, exists := ctx.Get("uploadedFiles")
	if !exists {
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.NoFilesUploaded, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	uploadedFiles := uploadedFilesAny.([]map[string]any)

	result, err := c.ProductUseCase.UploadProductImages(ctx, productUUID, uploadedFiles)
	if err != nil {
		c.Log.WithError(err).Error("Failed to upload product images")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.InternalServerError, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusCreated, constants.SuccessCreateProduct, result)
	ctx.JSON(res.StatusCode, res)
}

func (c *ProductController) GetProductImageURL(ctx *gin.Context) {
	imageID := ctx.Param("imageID")
	if imageID == "" {
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductID, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		c.Log.WithError(err).Error("Invalid image ID format")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductIDFormat, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}
	result, err := c.ImageUseCase.GetImageByID(ctx, imageUUID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get product image by ID")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedGetProductByID, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	if result == nil {
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.ProductNotFound, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	imageURL, err := c.MinioUseCase.GetPresignedURL(ctx, model.PresignedURLInput{
		Bucket:    c.Viper.GetString("MINIO_BUCKET_PRODUCT"),
		ObjectKey: result.ImageObject,
		Expiry:    int64((time.Hour * 24).Seconds()),
	})
	if err != nil {
		c.Log.WithError(err).Error("Failed to get presigned URL for product image")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedGetPresignedURL, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.SuccessGetProductByID, model.ProductImageURLResponse{
		ID:          result.ID,
		ProductID:   result.ProductID,
		ImageObject: result.ImageObject,
		URL:         imageURL,
	})
	ctx.JSON(res.StatusCode, res)
}

func (c *ProductController) GetObjectImage(ctx *gin.Context) {
	imageID := ctx.Param("imageID")
	if imageID == "" {
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductID, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	imageUUID, err := uuid.Parse(imageID)
	if err != nil {
		c.Log.WithError(err).Error("Invalid image ID format")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductIDFormat, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	result, err := c.ImageUseCase.GetImageByID(ctx, imageUUID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get product image by ID")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedGetImageByID, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	if result.ImageObject == "" {
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.ProductNotFound, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	object, err := c.MinioUseCase.GetObject(ctx, c.Viper.GetString("MINIO_BUCKET_PRODUCT"), result.ImageObject)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get object from Minio")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedGetPresignedURL, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}
	defer object.Close()

	info, err := object.Stat()
	if err != nil {
		c.Log.WithError(err).Error("Failed to get object info from Minio")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedGetPresignedURL, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	ctx.Header("Content-Type", info.ContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", info.Size))
	ctx.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", path.Base(result.ImageObject)))

	_, err = io.Copy(ctx.Writer, object)
	if err != nil {
		c.Log.WithError(err).Error("Failed to write object to response")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
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

	productID := ctx.Param("productID")
	if productID == "" {
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductID, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		c.Log.WithError(err).Error("Invalid product ID format")
		res := utils.FailedResponse(ctx, http.StatusBadRequest, constants.InvalidProductIDFormat, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	product, err := c.ProductUseCase.GetProduct(ctx, productUUID)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get product by ID")
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.FailedGetProductByID, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	err = c.ProductUseCase.DeleteProduct(ctx, product)
	if err != nil {
		c.Log.WithError(err).Error("Failed to delete product")
		res := utils.FailedResponse(ctx, http.StatusInternalServerError, constants.FailedDeleteProduct, err)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	res := utils.SuccessResponse(ctx, http.StatusOK, constants.SuccessDeleteProduct, true)
	ctx.JSON(res.StatusCode, res)
}
