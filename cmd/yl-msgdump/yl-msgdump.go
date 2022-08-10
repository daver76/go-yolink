package main

// Connects to YoLink MQTT broker and displays messages as they are received

import (
	"os"
	"os/signal"
	"syscall"

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

	log.Info("Connecting to MQTT broker...")
	err = client.MQTTConnect()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Waiting for messages...")
	// Run forever, until a signal is received.
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)
	<-keepAlive
	log.Warn("Shutting down...")
}

func messageHandler(topic string, payload []byte) {
	log.WithFields(log.Fields{
		"topic":   topic,
		"payload": string(payload),
	}).Info("Received message")
}
