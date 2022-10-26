package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/giannisalinetti/edge-retail-consumer/pkg/webserver"
)

// loadConf loads configurations from YAML file
func loadConf(path string) (map[interface{}]interface{}, error) {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := make(map[interface{}]interface{})

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {

	fmt.Println("Starting edge retail consumer...")

	yamlConfigMap := flag.String("config", "", "Client YAML config file")
	listenPort := flag.String("port", "8080", "Default listen port")
	flag.Parse()

	// Load configuration from YAML
	cfg, err := loadConf(*yamlConfigMap)
	if err != nil {
		log.Printf("Error loading conf: %v\n", err)
		os.Exit(1)
	}

	// Populate kafkaConfigMap
	var edgeRetailConfig = kafka.ConfigMap{
		"bootstrap.servers": cfg["bootstrap-servers"],
		"group.id":          cfg["group-id"],
		"auto.offset.reset": cfg["auto-offset-reset"],
		"security.protocol": cfg["security-protocol"],
		"sasl.mechanisms":   cfg["sasl-mechanisms"],
		"sasl.username":     cfg["sasl-username"],
		"sasl.password":     cfg["sasl-password"],
	}

	// Define topic from config
	var topic = cfg["topic"].(string)

	// Create new consumer
	c, err := kafka.NewConsumer(&edgeRetailConfig)
	if err != nil {
		log.Printf("Failed to create consumer. %v", err)
		os.Exit(1)
	}

	// Subscribe target topic(s)
	c.Subscribe(topic, nil)

	// Collect broker metadata
	meta, err := c.GetMetadata(&topic, true, -1)
	if err != nil {
		log.Printf("Error getting metadata: %v", err)
	}

	// Print broker metadata
	for _, broker := range meta.Brokers {
		log.Printf("Connected to host %v:%v, ID: %v\n", broker.Host, broker.Port, broker.ID)
	}

	// Start webserver as asynchronous goroutine
	go webserver.StartServer(*listenPort)

	log.Printf("Begin ReadMessage loop")

	// Loop over messages
	for {

		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}

		log.Printf("Key: %s, Message: %s\n", string(msg.Key), string(msg.Value))

		webserver.MessageMap[string(msg.Key)] = msg.Value
	}

}
