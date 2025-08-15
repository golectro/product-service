package route

import (
	"golectro-product/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func (c *RouteConfig) RegisterProductRoutes(rg *gin.RouterGroup, minioClient *minio.Client) {
	product := rg.Group("/products")

	product.GET("/", c.ProductController.GetAllProducts)
	product.GET("/:productID", c.ProductController.GetProductByID)
	product.GET("/search", c.ProductController.SearchProducts)
	product.POST("/", c.ProductController.CreateProduct)
	product.PUT("/:productID", c.AuthMiddleware, c.ProductController.UpdateProduct)
	product.POST("/:productID/images", c.AuthMiddleware, middleware.MultipleFileUpload(minioClient, middleware.UploadOptions{
		FieldName:     "images",
		MaxFileSizeMB: 5,
		MaxFiles:      5,
		BucketName:    c.Viper.GetString("MINIO_BUCKET_PRODUCT"),
		AllowedTypes:  []string{"image/jpeg", "image/png", "image/gif"},
	}), c.ProductController.UploadProductImages)
	product.GET("/image/:imageID/url", c.ProductController.GetProductImageURL)
	product.GET("/image/:imageID/preview", c.ProductController.GetObjectImage)
	product.DELETE("/:productID", c.AuthMiddleware, c.ProductController.DeleteProduct)
	product.DELETE("/image/:imageID", c.AuthMiddleware, c.ProductController.DeleteProductImage)
}
