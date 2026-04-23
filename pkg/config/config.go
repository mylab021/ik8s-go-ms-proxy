package config

type Config struct {
	UserServiceConfig    ServiceConfig
	OrderServiceConfig   ServiceConfig
	ProductServiceConfig ServiceConfig
	GatewayConfig        ServiceConfig
}

func NewConfig() *Config {
	return &Config{}
}

type ServiceConfig struct {
	Name        string
	BackendHost string `default:"localhost"`
	BackendPort int    `default:"8080"`
	ServicePort int    `default:"8080"`
}

func NewUserServiceConfig(name string, backendHost string, backendPort int) *ServiceConfig {
	return &ServiceConfig{
		Name:        name,
		BackendHost: backendHost,
		BackendPort: backendPort,
	}
}

func NewOrderServiceConfig(name string, backendHost string, backendPort int) *ServiceConfig {
	return &ServiceConfig{
		Name:        name,
		BackendHost: backendHost,
		BackendPort: backendPort,
	}
}

func NewProductServiceConfig(name string, backendHost string, backendPort int) *ServiceConfig {
	return &ServiceConfig{
		Name:        name,
		BackendHost: backendHost,
		BackendPort: backendPort,
	}
}

func NewGatewayConfig(name string, backendHost string, backendPort int) *ServiceConfig {
	return &ServiceConfig{
		Name:        name,
		BackendHost: backendHost,
		BackendPort: backendPort,
	}
}
