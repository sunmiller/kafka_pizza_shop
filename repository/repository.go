package repository

import "github.com/sunmiller/pizza-shop-eda/order-service/config"

type Repositories struct {
	OrderRepository IRepository
}

func GetRepositories() *Repositories {
	return &Repositories{
		OrderRepository: GetMongoRepository(config.GetEnvProperty("database_name"), "order"),
	}
}
