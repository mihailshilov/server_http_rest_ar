package main

import (
	"github.com/webdevolegkuprianov/server_http_rest_ar/app/apiserver"
	"github.com/webdevolegkuprianov/server_http_rest_ar/app/apiserver/model"

	logger "github.com/webdevolegkuprianov/server_http_rest_ar/app/apiserver/logger"
)

func main() {

	config, err := model.NewConfig()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	if err := apiserver.Start(config); err != nil {
		logger.ErrorLogger.Println(err)
	}

}
