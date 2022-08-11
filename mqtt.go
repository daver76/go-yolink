package yolink

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	random_charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type MessageHandler func(topic string, body []byte)

func (client *APIClient) MQTTConnect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", API_HOST, MQTT_PORT))
	opts.SetClientID(randomString(20))
	opts.SetUsername(client.AccessToken)
	opts.SetPassword("")
	opts.SetDefaultPublishHandler(client.messageHandler)
	opts.SetKeepAlive(30 * time.Second)
	opts.OnConnect = client.mqttConnectionHandler
	opts.OnConnectionLost = client.mqttConnectionLost
	opts.OnReconnecting = client.mqttReconnectHandler

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := mqttClient.Subscribe(fmt.Sprintf("yl-home/%s/+/report", client.HomeId),
		0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (client *APIClient) messageHandler(mqttClient mqtt.Client, msg mqtt.Message) {
	if client.MqttMessageHandler != nil {
		client.MqttMessageHandler(msg.Topic(), msg.Payload())
	} else {
		log.WithFields(log.Fields{
			"topic": msg.Topic(),
			"size":  len(msg.Payload()),
		}).Warn("Unhandled message received")
	}
}

func (client *APIClient) mqttConnectionHandler(mqttClient mqtt.Client) {
	log.Info("MQTT Connected")
}

func (client *APIClient) mqttReconnectHandler(mqttClient mqtt.Client, options *mqtt.ClientOptions) {
	log.Warn("MQTT Reconnecting...")
	// It is possible the old token expired, so just in case:
	_, err := client.getToken(false)
	if err != nil {
		log.Error("Error during reconnect: ", err)
		return
	}
	options.Username = client.AccessToken
}

func (client *APIClient) mqttConnectionLost(mqttClient mqtt.Client, err error) {
	log.Warn("MQTT Connection Lost: ", err)
}

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(random_charset[rand.Intn(len(random_charset))])
	}
	return sb.String()
}
