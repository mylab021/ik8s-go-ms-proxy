package config

import (
	"os"
	"strconv"
)

func InitConfig() *Config {
	gatewayConfig := NewGatewayConfig("Gateway", "localhost", 8080)
	userServiceConfig := NewUserServiceConfig("User Service", "localhost", 8081)
	orderServiceConfig := NewOrderServiceConfig("Order Service", "localhost", 8082)
	productServiceConfig := NewProductServiceConfig("Product Service", "localhost", 8083)

	if os.Getenv("USER_SERVICE_BACKEND_HOST") != "" {
		userServiceConfig.BackendHost = os.Getenv("USER_SERVICE_BACKEND_HOST")
	}
	if os.Getenv("USER_SERVICE_BACKEND_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("USER_SERVICE_BACKEND_PORT"))
		if err != nil {
			port = 8080
		}
		userServiceConfig.BackendPort = port
	}

	if os.Getenv("ORDER_SERVICE_BACKEND_HOST") != "" {
		orderServiceConfig.BackendHost = os.Getenv("ORDER_SERVICE_BACKEND_HOST")
	}
	if os.Getenv("ORDER_SERVICE_BACKEND_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("ORDER_SERVICE_BACKEND_PORT"))
		if err != nil {
			port = 8080
		}
		orderServiceConfig.BackendPort = port
	}

	if os.Getenv("PRODUCT_SERVICE_BACKEND_HOST") != "" {
		productServiceConfig.BackendHost = os.Getenv("PRODUCT_SERVICE_BACKEND_HOST")
	}
	if os.Getenv("PRODUCT_SERVICE_BACKEND_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("PRODUCT_SERVICE_BACKEND_PORT"))
		if err != nil {
			port = 8080
		}
		productServiceConfig.BackendPort = port
	}

	if os.Getenv("GATEWAY_SERVICE_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("GATEWAY_SERVICE_PORT"))
		if err != nil {
			port = 80
		}
		gatewayConfig.ServicePort = port
	} else {
		gatewayConfig.ServicePort = 80
	}

	if os.Getenv("USER_SERVICE_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("USER_SERVICE_PORT"))
		if err != nil {
			port = 8080
		}
		userServiceConfig.ServicePort = port
	} else {
		userServiceConfig.ServicePort = 8081
	}

	if os.Getenv("ORDER_SERVICE_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("ORDER_SERVICE_PORT"))
		if err != nil {
			port = 8080
		}
		orderServiceConfig.ServicePort = port
	} else {
		orderServiceConfig.ServicePort = 8082
	}

	if os.Getenv("PRODUCT_SERVICE_PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("PRODUCT_SERVICE_PORT"))
		if err != nil {
			port = 8080
		}
		productServiceConfig.ServicePort = port
	} else {
		productServiceConfig.ServicePort = 8083
	}

	config := NewConfig()
	config.UserServiceConfig = *userServiceConfig
	config.OrderServiceConfig = *orderServiceConfig
	config.ProductServiceConfig = *productServiceConfig
	config.GatewayConfig = *gatewayConfig

	return config
}
