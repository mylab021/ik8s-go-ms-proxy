package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mylab021/ik8s-ms-proxy/pkg/config"
	"github.com/mylab021/ik8s-ms-proxy/pkg/logger"

	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func GetAllIPs() ([]string, error) {
	config.InitConfig()
	logger.NewLogger()
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

func GetK8SInfo() map[string]string {
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
	serverInfo["APP_NAME"] = "Gateway Service"
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
		"K8SInfo":        GetK8SInfo(),
		"ServerInfo":     GetServerInfo(),
		"ClientInfo":     GetClientInfo(ctx),
	})

}

func GetUserServiceInfo(ctx *gin.Context) {
	// 创建一个客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	targetURL := os.Getenv("USER_SERVICE_URL")
	if targetURL == "" {
		targetURL = "http://user-service:8080"
	}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		panic(err)
	}
	if ctx.Request.Header.Get("X-Forwarded-For") != "" {
		req.Header.Set("X-Forwarded-For", ctx.Request.Header.Get("X-Forwarded-For"))
	} else {
		req.Header.Set("X-Forwarded-For", ctx.ClientIP())
	}

	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch external data",
		})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read response",
		})
		return
	}
	// 解析 JSON
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse JSON",
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func GetOrderServiceInfo(ctx *gin.Context) {
	targetURL := os.Getenv("ORDER_SERVICE_URL")
	if targetURL == "" {
		targetURL = "http://order-service:8080"
	}
	target, err := url.Parse(targetURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse URL",
		})
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		if req.Header.Get("X-Forwarded-For") != "" {
			req.Header.Set("X-Forwarded-For", req.Header.Get("X-Forwarded-For"))
		} else {
			req.Header.Set("X-Forwarded-For", ctx.ClientIP())
		}
		req.URL.Path = "/"
	}
	proxy.Transport = &http.Transport{}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
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

	router.GET("/user-service", GetUserServiceInfo)
	router.GET("/order-service", GetOrderServiceInfo)

	srv := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
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
