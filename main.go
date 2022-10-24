package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var edgeRetailConfig = kafka.ConfigMap{
	"bootstrap.servers": "rh-summit-cdbcno-om-e-i---m-dg.bf2.kafka.rhcloud.com:443",
	"group.id":          "retail-store",
	"auto.offset.reset": "earliest",
	"security.protocol": "sasl_ssl",
	"sasl.mechanisms":   "PLAIN",
	"sasl.username":     "8efb13d6-8a0d-41fe-ba90-f9101270dbcc",
	"sasl.password":     "X8VwUUz8VHahOJOe6xFPKb125XueRhox",
}

var topicList = []string{"pickup"}

func main() {

	fmt.Println("Starting edge retail consumer...")

	// Create new consumer
	c, err := kafka.NewConsumer(&edgeRetailConfig)
	if err != nil {
		fmt.Printf("Failed to create consumer. %v", err)
	}

	// Subscribe target topic(s)
	c.SubscribeTopics(topicList, nil)

	// Collect broker metadata
	meta, err := c.GetMetadata(&topicList[0], true, -1)
	if err != nil {
		fmt.Printf("Error getting metadata: %v", err)
	}

	// Print broker metadata
	for _, broker := range meta.Brokers {
		log.Printf("Connected to host %v:%v, ID: %v\n", broker.Host, broker.Port, broker.ID)
	}

	// Loop over messages
	for {

		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		log.Printf("Message: %s\n", string(msg.Value))
	}

}
