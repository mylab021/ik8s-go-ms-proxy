package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mylab021/ik8s-go-ms-proxy/pkg/common"
	"github.com/mylab021/ik8s-go-ms-proxy/pkg/config"
	"github.com/mylab021/ik8s-go-ms-proxy/pkg/handler"
	"github.com/mylab021/ik8s-go-ms-proxy/pkg/utils"
)

func main() {
	router := gin.Default()

	cfg := config.LoadConfig()

	serverHandler := handler.NewServerHandler(cfg.OrderServiceConfig.Name)

	router.GET("/", serverHandler.GetHttpInfo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	common.Run(router, cfg.OrderServiceConfig.Name, utils.GenListenAddress(cfg.OrderServiceConfig.BackendPort))
}
