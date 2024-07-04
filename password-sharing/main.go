package main

import (
	"github.com/getkrabhijeet/go101/password-sharing/config"
	"github.com/getkrabhijeet/go101/password-sharing/controller"
	"github.com/getkrabhijeet/go101/password-sharing/database"
	"github.com/getkrabhijeet/go101/password-sharing/health"
	"github.com/getkrabhijeet/go101/password-sharing/helper"
	"github.com/getkrabhijeet/go101/password-sharing/logger"
	"github.com/getkrabhijeet/go101/password-sharing/server"
	"github.com/getkrabhijeet/go101/password-sharing/service"
)

func main() {
	appConfiguration, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	appLogger := logger.NewLoggerFactory(appConfiguration)

	encoder := helper.NewEncoder(appConfiguration)
	databaseFactory := database.NewFactory(appConfiguration, appLogger)
	randomFactory := helper.NewRandomFactory()
	service := service.NewPasswordService(databaseFactory, appConfiguration, randomFactory, appLogger, encoder)

	pgHealthCheck := health.NewPgHealthCheck(databaseFactory, appLogger)

	server := server.NewServer(
		appLogger,
		appConfiguration,
		controller.NewCreateLinkController(service, appConfiguration),
		controller.NewGetLinkController(service),
		controller.NewHealthController(pgHealthCheck),
	)

	if err = server.Run(); err != nil {
		panic(err)
	}
}
