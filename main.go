package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/Shopify/sarama"

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
		panic(err)
	}

	return config
}

func setupRouter(config *viper.Viper) *gin.Engine {
	gin.SetMode(config.GetString("gin.release_mode"))

	r := gin.Default()

	bootstrapServers := strings.Split(config.GetString("kafka.bootstrap_servers"), ",")
	topic := config.GetString("kafka.topic")

	producer, err := sarama.NewSyncProducer(bootstrapServers, nil)
	if err != nil {
		log.Error("Error creating the sync producer")
		panic(err)
	}

	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.POST("/produce", func(c *gin.Context) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)

		if err != nil {
			log.Error("Error parsing JSON body")
			c.String(http.StatusBadRequest, "can't parse JSON body")
			return
		}

		message := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(jsonData)}
		partition, offset, err := producer.SendMessage(&message)

		if err != nil {
			log.Error("Error writing to Kafka")
			log.Error(err)
			c.String(http.StatusInternalServerError, "can't write to Kafka")
			return
		}

		c.JSON(http.StatusOK, gin.H{"partition": partition, "offset": offset})
	})

	return r
}

func main() {
	config := getConfig()

	addr := config.GetString("server.address") + ":" + config.GetString("server.port")
	log.Info("Server addr: " + addr)

	r := setupRouter(config)
	err := r.Run(addr)

	if err != nil {
		log.Error("Error starting gin server")
		panic(err)
	}
}
