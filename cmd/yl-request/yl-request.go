package main

// CLI utility to query device state

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/daver76/go-yolink"
	log "github.com/sirupsen/logrus"
)

var CLI struct {
	ListDevices struct {
	} `cmd:"" help:"List devices."`

	GetState struct {
		DeviceId string `arg:"" name:"deviceid" help:"DeviceID" required:""`
	} `cmd:"" help:"Get device state."`
}

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

	ctx := kong.Parse(&CLI)
	var output any

	switch ctx.Command() {
	case "list-devices":
		output, err = client.GetDevices()
	case "get-state <deviceid>":
		output, err = getDeviceState(client, ctx.Args[1])
	default:
		panic(ctx.Command())
	}

	if err != nil {
		log.Error(err)
	} else if output != nil {
		js, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Error(err)
		} else {
			fmt.Println(string(js))
		}
	}
}

func getDeviceState(client *yolink.APIClient, deviceID string) (any, error) {
	var ret any = nil
	devices, err := client.GetDevices()
	if err != nil {
		return ret, err
	}
	var device *yolink.Device
	for _, d := range devices {
		if d.DeviceID == deviceID {
			device = &d
			break
		}
	}
	if device == nil {
		return ret, fmt.Errorf("unknown deviceID: %s", deviceID)
	}
	return client.GetState(*device)
}
