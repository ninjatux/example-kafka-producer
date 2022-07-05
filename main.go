package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func getConfig() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config.yaml")
	config.AddConfigPath(".")
	config.SetConfigType("yaml")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := config.ReadInConfig()
	if err != nil {
		log.Error("Error loading config")
		log.Error(err)
	}

	return config
}

func setupRouter(config *viper.Viper) *gin.Engine {
	gin.SetMode(config.GetString("gin.release_mode"))

	r := gin.Default()

	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.POST("/produce", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}

func main() {
	config := getConfig()

	addr := config.Get("server.address").(string) + ":" + config.Get("server.port").(string)
	log.Info("Server addr: " + addr)

	r := setupRouter(config)
	err := r.Run(addr)

	if err != nil {
		log.Error("Error starting gin server")
		log.Error(err)
	}
}
