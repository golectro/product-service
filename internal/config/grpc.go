package config

import (
	"golectro-product/internal/delivery/grpc"
	"golectro-product/internal/repository"
	"golectro-product/internal/usecase"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func StartGRPC(viper *viper.Viper, db *gorm.DB, validate *validator.Validate, log *logrus.Logger, elastic *elasticsearch.Client) {
	productRepository := repository.NewProductRepository(log)
	imageRepository := repository.NewImageRepository(log)

	elasticsearchUseCase := usecase.NewElasticsearchUsecase(elastic, log, validate, viper)
	productUseCase := usecase.NewProductUsecase(db, log, validate, productRepository, imageRepository, elasticsearchUseCase)

	port := viper.GetInt("GRPC_PORT")
	grpc.StartGRPCServer(productUseCase, port, viper)
}
