package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mylab021/ik8s-ms-proxy/pkg/common"
	"github.com/mylab021/ik8s-ms-proxy/pkg/config"
	"github.com/mylab021/ik8s-ms-proxy/pkg/handler"
	"github.com/mylab021/ik8s-ms-proxy/pkg/utils"
	"net/http"
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

	common.Run(router, cfg.ProductServiceConfig.Name, utils.GenListenAddress(cfg.ProductServiceConfig.BackendPort))
}
