package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sunmiller/pizza-shop-eda/order-service/handler"
	"github.com/sunmiller/pizza-shop-eda/order-service/service"
)

func registerOrderRoutes(r *gin.RouterGroup, publisher service.IMessagePublisher) {

	orderHandler := handler.GetOrderHandler(publisher)

	r.POST(
		"/create",
		orderHandler.CreateOrder,
	)
}
