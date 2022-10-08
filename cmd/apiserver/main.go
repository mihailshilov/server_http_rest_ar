package main

import (
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"

	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

// @title       Repair App API
// @version     9.0
// @description API Server for TodoList Application

// @host     carsrv.shilov.pro
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization

func main() {

	config, err := model.NewConfig()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	if err := apiserver.Start(config); err != nil {
		logger.ErrorLogger.Println(err)
	}

}
