package messageconsumer

import (
	"context"
	"fmt"

	"github.com/sunmiller/pizza-shop-eda/order-service/logger"
	"github.com/sunmiller/pizza-shop-eda/order-service/repository"
	"github.com/sunmiller/pizza-shop-eda/order-service/service"
)

type OrderMessageConsumer struct {
	consumer             service.IMessageConsumer
	orderConsumerChannel chan service.Message
	workerCount          int
	repositories         repository.Repositories
}

func (omc *OrderMessageConsumer) StartConsuming() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < omc.workerCount; i++ {
		go omc.registerConsumerWorker(i, ctx)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Log("Stopping Message Consumption..")
			return
		default:
			message, err := omc.consumer.ConsumeMessage()
			if err != nil {
				continue
			}

			select {
			case omc.orderConsumerChannel <- message:
			default:
				logger.Log("worker pool is busy, dropping message")
			}
		}
	}
}

func (omc *OrderMessageConsumer) registerConsumerWorker(id int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-omc.orderConsumerChannel:
			logger.Log(fmt.Sprintf("worker %d - Processed Message : %v", id, message.Data))
			_, err := omc.repositories.OrderRepository.Create(message.Data, nil)
			if err != nil {
				logger.Log(fmt.Sprintf("failed to save Data to MongoDb order collection %v, Consumer Id - %v", id, err))
			} else {
				omc.consumer.GetReader().CommitMessages(ctx, message.KafkaMessage)
			}

		}
	}
}

func GetOrderMessageConsumer(consumer service.IMessageConsumer, repositories repository.Repositories) *OrderMessageConsumer {
	return &OrderMessageConsumer{
		consumer:             consumer,
		repositories:         repositories,
		orderConsumerChannel: make(chan service.Message, 100),
		workerCount:          10,
	}
}
