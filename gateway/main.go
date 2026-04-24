package main

import (
	"log"
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

	userTargetURL := utils.GenTargetURL(cfg.UserServiceConfig.BackendHost, cfg.UserServiceConfig.BackendPort)
	orderTargetURL := utils.GenTargetURL(cfg.OrderServiceConfig.BackendHost, cfg.OrderServiceConfig.BackendPort)
	productTargetURL := utils.GenTargetURL(cfg.ProductServiceConfig.BackendHost, cfg.ProductServiceConfig.BackendPort)

	log.Printf("UserService Backend URL: %s", userTargetURL)
	log.Printf("OrderService Backend URL: %s", orderTargetURL)
	log.Printf("ProductService Backend URL: %s", productTargetURL)

	serverHandler := handler.NewServerHandler(cfg.GatewayConfig.Name)

	router.GET("/", serverHandler.GetHttpInfo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	userHttpProxy := handler.NewHttpProxy(userTargetURL)
	orderHttpProxy := handler.NewHttpProxy(orderTargetURL)
	productHttpProxy := handler.NewHttpProxy(productTargetURL)

	router.GET("/user-service", userHttpProxy.GetBackendService)
	router.GET("/order-service", orderHttpProxy.GetBackendService)
	router.GET("/product-service", productHttpProxy.GetBackendService)

	common.Run(router, cfg.GatewayConfig.Name, utils.GenListenAddress(cfg.GatewayConfig.ServicePort))
}
