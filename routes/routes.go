package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sunmiller/pizza-shop-eda/order-service/service"
)

func RegisterRoutes(r *gin.Engine, publisher service.IMessagePublisher) {
	orderRoutes := r.Group("/order-service")
	{
		registerOrderRoutes(orderRoutes, publisher)
	}
}
