package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunmiller/pizza-shop-eda/order-service/config"
	"github.com/sunmiller/pizza-shop-eda/order-service/constants"

	// "github.com/sunmiller/pizza-shop-eda/order-service/repository"
	"github.com/sunmiller/pizza-shop-eda/order-service/routes"
	"github.com/sunmiller/pizza-shop-eda/order-service/service"
)

var webPort = 8001

func main() {
	log.Println("🟢 Docker Go app started, waiting...")
	log.Println("This is standard error (via log)")
	log.Println("Kafka host:", config.GetEnvProperty("kafka_host"))
	log.Println("Kafka port:", config.GetEnvProperty("kafka_port"))
	log.Println(os.Getenv("LOG"))
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(gin.Recovery())

	app.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "service is up and running fine",
		})
	})
	log.Println("This is start of repositories initialize")
	// var repositories = repository.GetRepositories()

	// consumer := service.GetNewKafkaConsumer(constants.TOPIC_ORDER, "order-message")
	// var orderMessageConsumer = messageconsumer1.GetOrderMessageConsumer(consumer, *repositories)

	// go orderMessageConsumer.StartConsuming()

	log.Println("This is completion of orderMessageConsumer")

	routes.RegisterRoutes(
		app,
		service.GetKafkaMessagePublisher(constants.TOPIC_ORDER, "order-service"),
	)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", strconv.Itoa(webPort)),
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		IdleTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, //1 mb
	}
	fmt.Println("started web server on port ", webPort)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}
