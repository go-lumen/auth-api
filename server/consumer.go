package server

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	log "github.com/sirupsen/logrus"

	"fmt"
)

func (a *API) SetupConsumer() {
	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init consumer
	brokerList := a.Config.GetStringSlice("brokers")
	topics := []string{a.Config.GetString("kafka_topic")}
	consumer, err := cluster.NewConsumer(brokerList, a.Config.GetString("group_id"), topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Error(fmt.Sprintf("Error: %s", err.Error()))
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Debug(fmt.Sprintf("Rebalanced: %+v", ntf))
		}
	}()

	// consume messages, watch signals
	for {
		msg, ok := <-consumer.Messages()
		if ok {
			a.TopicRouter.CallHandler(msg.Topic, msg.Value)
			consumer.MarkOffset(msg, "") // mark message as processed
		}
	}
}
