package config

import (
	"golectro-product/internal/delivery/http"
	"golectro-product/internal/delivery/http/middleware"
	"golectro-product/internal/delivery/http/route"
	"golectro-product/internal/repository"
	"golectro-product/internal/usecase"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/vault/api"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB         *gorm.DB
	Mongo      *mongo.Database
	App        *gin.Engine
	Redis      *redis.Client
	Log        *logrus.Logger
	Validate   *validator.Validate
	Viper      *viper.Viper
	GRPCClient *grpc.ClientConn
	Elastic    *elasticsearch.Client
	Minio      *minio.Client
	Vault      *api.Client
}

func Bootstrap(config *BootstrapConfig) {
	productRepository := repository.NewProductRepository(config.Log)
	minioRepository := repository.NewMinioRepository(config.Minio)
	imageRepository := repository.NewImageRepository(config.Log)

	elasticsearchUseCase := usecase.NewElasticsearchUsecase(config.Elastic, config.Log, config.Validate, config.Viper)
	productUseCase := usecase.NewProductUsecase(config.DB, config.Log, config.Validate, productRepository, imageRepository, elasticsearchUseCase)
	minioUseCase := usecase.NewMinioUsecase(minioRepository, config.Validate, config.Log)
	imageUseCase := usecase.NewImageUsecase(config.DB, config.Log, config.Validate, imageRepository)

	productController := http.NewProductController(productUseCase, minioUseCase, config.Log, config.Viper, imageUseCase)

	authMiddleware := middleware.NewAuth(config.Viper)

	routeConfig := route.RouteConfig{
		App:               config.App,
		AuthMiddleware:    authMiddleware,
		Minio:             config.Minio,
		Viper:             config.Viper,
		ProductController: productController,
	}
	routeConfig.Setup()
}
