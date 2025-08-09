package route

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func (c *RouteConfig) RegisterProductRoutes(rg *gin.RouterGroup, minioClient *minio.Client) {
	product := rg.Group("/products")

	product.GET("/", c.AuthMiddleware, c.ProductController.GetAllProducts)
	product.POST("/", c.AuthMiddleware, c.ProductController.CreateProduct)
}
