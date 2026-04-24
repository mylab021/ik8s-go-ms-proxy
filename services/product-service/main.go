package main

import (
	"net/http"

	"pkg/common"
	"pkg/config"
	"pkg/handler"
	"pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	cfg := config.LoadConfig()

	serverHandler := handler.NewServerHandler(cfg.ProductServiceConfig.Name)

	router.GET("/", serverHandler.GetHttpInfo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	common.Run(router, cfg.ProductServiceConfig.Name, utils.GenListenAddress(cfg.ProductServiceConfig.ServicePort))
}
