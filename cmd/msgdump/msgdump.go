package main

import (
	"os"
	"time"

	"github.com/daver76/go-yolink"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(
		&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true},
	)

	client, err := yolink.NewAPIClient(
		os.Getenv("YOLINK_UAID"),
		os.Getenv("YOLINK_SECRET_KEY"),
	)
	if err != nil {
		log.Error(err)
		return
	}

	client.MqttMessageHandler = messageHandler

	log.Info("Connecting to MQTT...")
	err = client.MQTTConnect()
	if err != nil {
		log.Error(err)
		return
	}

	for {
		log.Debug("Waiting...")
		time.Sleep(10 * time.Second)
	}
}

func messageHandler(topic string, payload []byte) {
	log.WithFields(log.Fields{
		"topic":   topic,
		"payload": string(payload),
	}).Info("Received message")
}
