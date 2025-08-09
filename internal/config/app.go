package config

import (
	"golectro-product/internal/delivery/http/middleware"
	"golectro-product/internal/delivery/http/route"

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
	authMiddleware := middleware.NewAuth(config.Viper)

	routeConfig := route.RouteConfig{
		App:               config.App,
		AuthMiddleware:    authMiddleware,
		Minio:             config.Minio,
		Viper:             config.Viper,
	}
	routeConfig.Setup()
}
