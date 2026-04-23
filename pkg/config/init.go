package config

func InitConfig() *Config {
	gatewayConfig := NewGatewayConfig("Gateway", "localhost", 8080)
	userServiceConfig := NewUserServiceConfig("User Service", "localhost", 8081)
	orderServiceConfig := NewOrderServiceConfig("Order Service", "localhost", 8082)
	productServiceConfig := NewProductServiceConfig("Product Service", "localhost", 8083)

	config := NewConfig()
	config.UserServiceConfig = *userServiceConfig
	config.OrderServiceConfig = *orderServiceConfig
	config.ProductServiceConfig = *productServiceConfig
	config.GatewayConfig = *gatewayConfig

	return config
}
