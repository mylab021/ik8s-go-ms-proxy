package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type httpProxy struct {
	targetURL string
}

func NewHttpProxy(targetURL string) *httpProxy {
	return &httpProxy{targetURL: targetURL}
}

func (ht *httpProxy) GetBackendService(ctx *gin.Context) {
	targetURL := ht.targetURL
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
