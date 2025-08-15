package main

import (
	"golectro-product/internal/command"
	"golectro-product/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	elasticsearch := config.NewElasticSearch(viperConfig, log)

	if !command.NewCommandExecutor(viperConfig, db).Execute(log) {
		return
	}

	config.StartGRPC(viperConfig, db, validate, log, elasticsearch)
}
