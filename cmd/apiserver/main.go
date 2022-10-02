package main

import (
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"

	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

// @title Service App API
// @version 1.0
// @API для сбора данных сервисных станций

// @host localhost:443
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	config, err := model.NewConfig()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	if err := apiserver.Start(config); err != nil {
		logger.ErrorLogger.Println(err)
	}

}
