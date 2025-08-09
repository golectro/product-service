package main

import (
	"fmt"
	"golectro-product/internal/command"
	"golectro-product/internal/config"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogger(viper)
	db := config.NewDatabase(viper, log)
	mongo := config.NewMongoDB(viper, log)
	validate := config.NewValidator(viper)
	redis := config.NewRedis(viper, log)
	vault := config.NewVaultClient(viper, log)
	app := config.NewGin(viper, log, mongo, redis)
	minio := config.NewMinioClient(viper, log)
	elasticsearch := config.NewElasticSearch(viper, log)
	executor := command.NewCommandExecutor(db)

	config.Bootstrap(&config.BootstrapConfig{
		Viper:    viper,
		Log:      log,
		DB:       db,
		Mongo:    mongo,
		Validate: validate,
		App:      app,
		Redis:    redis,
		Minio:    minio,
		Vault:    vault,
		Elastic:  elasticsearch,
	})

	if !executor.Execute(log) {
		return
	}

	webPort := viper.GetInt("PORT")
	err := app.Run(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
