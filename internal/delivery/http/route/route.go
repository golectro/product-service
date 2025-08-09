package route

import (
	"golectro-product/internal/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App               *gin.Engine
	Minio             *minio.Client
	AuthMiddleware    gin.HandlerFunc
	Viper             *viper.Viper
	SwaggerController *http.SwaggerController
}

func (c *RouteConfig) Setup() {
	api := c.App.Group("/api")

	c.RegisterCommonRoutes(c.App)
	c.RegisterSwaggerRoutes(api)
}
