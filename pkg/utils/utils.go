package utils

import "fmt"

func GenListenAddress(port int) string {
	return fmt.Sprintf(":%d", port)
}

func GenTargetURL(backendHost string, backendPort int) string {
	return fmt.Sprintf("http://%s:%d", backendHost, backendPort)
}
