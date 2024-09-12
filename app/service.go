package main

import (
	"context"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/authnull0/database-service/src/controller"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"log"
)

var env string

func start() {
	loadConfig()

	gin.DisableConsoleColor()

	r := gin.Default()

	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowAllOrigins = true
	r.Use(CORSMiddleware())
	setupRoutes(r)
	startServer(r)
}

// startServer - Start server
func startServer(r *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + viper.GetString(env+"server.port"),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Default().Printf("Shutting down server...\n")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Default().Printf("Server forced to shutdown: %s\n", err)
	}

	log.Default().Printf("Server exiting\n")
}

func setupRoutes(r *gin.Engine) *gin.Engine {
	dbController := controller.DbController{}

	//API to sync database
	r.POST("api/v1/databaseService/dbSync", dbController.DbSync)
	//API to sync users
	r.POST("api/v1/databaseService/dbUser", dbController.DbUser)
	//API to sync privilege
	r.POST("api/v1/databaseService/dbPrivilege", dbController.DbPrivilege)
	//API to list Database
	r.POST("api/v1/databaseService/listDatabase", dbController.ListDatabase)
	//API to list Privilege
	r.POST("api/v1/databaseService/listUserPrivilege", dbController.ListUserPrivilege)

	return r
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("..")
	viper.AddConfigPath("conf")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file, %s", err)
	}
	env = viper.GetString("env") + "."

}

func main() {
	start()
}

func CORSMiddleware() gin.HandlerFunc {
	logrus.Info("Middleware:CORSMiddleware")
	return func(c *gin.Context) {

		// allowOrigin := viper.GetString(env + "cors.allowOrigin")
		logrus.Info("Middleware:CORSMiddleware: allowOrigin: ", "*")

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Authorization, withCredentials ,User-Agent")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
