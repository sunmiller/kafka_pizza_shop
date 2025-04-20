package logger

import (
	"log"
	"os"

	"github.com/sunmiller/pizza-shop-eda/order-service/utils"
)

func Log(message any) {
	isLoggedEnabled := os.Getenv("LOG")

	if isLoggedEnabled != "" {
		log.Println(message)
		utils.AppendToFile("log.txt", message)
	}
}
