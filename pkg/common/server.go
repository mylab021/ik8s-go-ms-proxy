package common

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, serverName string, addr string) {
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	//保证下面的优雅启停
	go func() {
		log.Printf("%s running in %s \n", serverName, addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	//SIGINT 用户发送INTR字符(Ctrl+C)触发
	//SIGTERM 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down project %s... \n", serverName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %v \n", serverName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("关闭超时")
	}
	log.Printf("%s stop success...", serverName)
}
