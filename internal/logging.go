package internal

import (
	"io"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mLogger struct {
	mclient mqtt.Client
}

var logger mLogger

func (m mLogger) Write(p []byte) (n int, err error) {
	mqttPublish(m.mclient, Config.mqttLogTopic, p, 0)
	return len(p), nil
}

func LogHex(name string, command []byte) {
	if Config.LogHexDump {
		log.Printf("%s: %X\n", name, command)
	}
}

func RedirectLog(mclient mqtt.Client) {
	logger.mclient = mclient

	if Config.LogMqtt {
		log.Println("Enabling logging to MQTT")
		log.SetOutput(io.MultiWriter(logger, log.Writer()))
	}
}
