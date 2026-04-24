package handler

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type serverHandler struct {
	ServiceName string
}

func NewServerHandler(serviceName string) *serverHandler {
	return &serverHandler{
		ServiceName: serviceName,
	}
}

func getAllIPs() ([]string, error) {
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

func getRequestHeaders(ctx *gin.Context) map[string]interface{} {
	headers := make(map[string]interface{})
	for key, value := range ctx.Request.Header {
		headers[key] = value[0]
	}
	return headers
}

func getK8SInfo(ctx *gin.Context) map[string]interface{} {
	k8sInfo := make(map[string]interface{})
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

func getServerInfo(serviceName string) map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	serverInfo := make(map[string]interface{})
	serverInfo["ServiceName"] = serviceName
	serverInfo["HostName"], _ = os.Hostname()
	ips, err := getAllIPs()
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

func getClientInfo(ctx *gin.Context) map[string]interface{} {
	clientInfo := make(map[string]interface{})
	clientInfo["Client IP"] = ctx.ClientIP()
	clientInfo["url"] = ctx.Request.URL.String()
	clientInfo["User Agent"] = ctx.Request.UserAgent()
	clientInfo["Host"] = ctx.Request.Host
	clientInfo["Method"] = ctx.Request.Method
	clientInfo["X-Forwarded-For"] = ctx.Request.Header.Get("X-Forwarded-For")
	return clientInfo
}

func (s *serverHandler) GetHttpInfo(ctx *gin.Context) {
	//response := make(map[string]interface{})
	//response["RequestHeaders"] = getRequestHeaders(ctx)
	//response["K8SInfo"] = getK8SInfo(ctx)
	//response["ServerInfo"] = getServerInfo(s.ServiceName)
	//response["ClientInfo"] = getClientInfo(ctx)
	//ctx.JSON(http.StatusOK, response)
	//log.Printf("K8S Server Info %v", getK8SInfo(ctx))
	ctx.HTML(http.StatusOK, "server_info.html", gin.H{
		"title":         s.ServiceName,
		"k8sInfo":       getK8SInfo(ctx),
		"ServerInfo":    getServerInfo(s.ServiceName),
		"ClientInfo":    getClientInfo(ctx),
		"RequestHeader": getRequestHeaders(ctx),
	})
}
