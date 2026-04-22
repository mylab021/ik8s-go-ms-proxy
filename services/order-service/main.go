package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllIPs() ([]string, error) {
	var ips []string
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网络接口失败: %v", err)
	}

	for _, iface := range interfaces {
		// 过滤掉回环接口和未启用接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口地址
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 过滤IPv6和空地址
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	return ips, nil
}

func GetRequestHeaders(ctx *gin.Context) map[string]string {
	headers := map[string]string{}
	for key, value := range ctx.Request.Header {
		headers[key] = value[0]
	}
	return headers
}

func GetK8SInfo(ctx *gin.Context) map[string]string {
	k8sInfo := make(map[string]string)
	if os.Getenv("NODE_NAME") != "" {
		k8sInfo["K8S Node Name"] = os.Getenv("NODE_NAME")
	}
	if os.Getenv("NODE_IP") != "" {
		k8sInfo["K8S Node IP"] = os.Getenv("NODE_IP")
	}
	if os.Getenv("POD_NAMESPACE") != "" {
		k8sInfo["K8S Pod Namespace"] = os.Getenv("POD_NAMESPACE")
	}
	if os.Getenv("POD_NAME") != "" {
		k8sInfo["K8S Pod Name"] = os.Getenv("POD_NAME")
	}
	if os.Getenv("POD_IP") != "" {
		k8sInfo["K8S Pod IP"] = os.Getenv("POD_IP")
	}
	return k8sInfo
}

func GetServerInfo() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	serverInfo := make(map[string]interface{})
	serverInfo["APP_NAME"] = "Order Service"
	serverInfo["HostName"], _ = os.Hostname()
	ips, err := GetAllIPs()
	if err == nil {
		serverInfo["Server IPs"] = ips
	}
	serverInfo["GoVersion"] = runtime.Version()
	serverInfo["OS"] = runtime.GOOS
	serverInfo["Architecture"] = runtime.GOARCH
	serverInfo["NumCPU"] = strconv.Itoa(runtime.NumCPU())
	serverInfo["Goroutines"] = strconv.Itoa(runtime.NumGoroutine())
	serverInfo["Uptime"] = strconv.FormatInt(time.Now().Unix(), 10)
	return serverInfo
}

func GetClientInfo(ctx *gin.Context) map[string]interface{} {
	clientInfo := make(map[string]interface{})
	clientInfo["Client IP"] = ctx.ClientIP()
	clientInfo["url"] = ctx.Request.URL.String()
	clientInfo["User Agent"] = ctx.Request.UserAgent()
	clientInfo["Host"] = ctx.Request.Host
	clientInfo["Method"] = ctx.Request.Method
	clientInfo["X-Forwarded-For"] = ctx.Request.Header.Get("X-Forwarded-For")
	return clientInfo
}

func GetInfo(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"RequestHeaders": GetRequestHeaders(ctx),
		"K8SInfo":        GetK8SInfo(ctx),
		"ServerInfo":     GetServerInfo(),
		"ClientInfo":     GetClientInfo(ctx),
	})

}

func main() {
	router := gin.Default()

	router.GET("/", GetInfo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
