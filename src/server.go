package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	apicontroller "github.com/gregorsoll/gomicroservice/controllers/api_controller"
	indexController "github.com/gregorsoll/gomicroservice/controllers/index_controller"
	zaploki "github.com/paul-milne/zap-loki"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	viper.SetConfigFile("server.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	logger, _ := initLogger()
	defer logger.Sync() // flushes buffer, if any

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// push gin logs to zap
	router.Use(ginzap.RecoveryWithZap(logger, true))

	router.StaticFile("/favicon.ico", "./favicon.ico")

	indexController.IndexRoutes(&router.RouterGroup, "/")
	apicontroller.IndexRoutes(&router.RouterGroup, "/api")

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", viper.GetString("server.ip"), viper.GetString("server.port")),
		Handler: router,
	}

	go func() {
		// service connections
		logger.Info("Listen at " + srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.DPanic("listen", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 2 seconds.")
	}
	log.Println("Server exiting")
}

func initLogger() (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	loki := zaploki.New(context.Background(), zaploki.Config{
		Url:          viper.GetString("loki.url"),
		BatchMaxSize: 1000,
		BatchMaxWait: 10 * time.Second,
		Labels:       map[string]string{"app": viper.GetString("loki.appname")},
	})

	return loki.WithCreateLogger(zapConfig)
}
