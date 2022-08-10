package main

// Connects to YoLink MQTT broker and writes messages to Postgres
// see ../../postgres_setup.sql

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/daver76/go-yolink"
	pgx "github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

var ErrInvalidFormat = errors.New("invalid format")

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

	// Example: YOLINK_DB_URL="postgres://username:password@localhost:5432/yolink"
	conn, err := pgx.Connect(context.Background(), os.Getenv("YOLINK_DB_URL"))
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close(context.Background())

	messageHandler := func(topic string, payload []byte) {
		err := saveMessage(conn, topic, payload)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Unable to save message")
		}
		log.WithFields(log.Fields{
			"topic":   topic,
			"payload": string(payload),
		}).Info("Received message")
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

func saveMessage(conn *pgx.Conn, topic string, payload []byte) error {
	data := make(map[string]any)
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return err
	}
	split_topic := strings.Split(topic, "/")
	if len(split_topic) < 4 {
		return ErrInvalidFormat
	}
	home_id := split_topic[1]

	var dev_id string
	var timestamp float64
	var ok bool

	if timestamp, ok = data["time"].(float64); !ok {
		return ErrInvalidFormat
	}
	timestamp /= 1000.0

	if dev_id, ok = data["deviceId"].(string); !ok {
		return ErrInvalidFormat
	}
	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO mqtt_messages(home_id, dev_id, time, data) VALUES ($1, $2, to_timestamp($3), $4)",
		home_id,
		dev_id,
		timestamp,
		string(payload),
	)

	return err
}
