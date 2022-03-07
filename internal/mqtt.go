package internal

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttPublish(mclient mqtt.Client, topic string, data interface{}, qos byte) {
	token := mclient.Publish(topic, qos, true, data)
	if token.Wait() && token.Error() != nil {
		log.Printf("Failed to publish, %v", token.Error())
	}
}

func MakeMQTTConn(commandsChannel chan []byte) mqtt.Client {
	log.Println("Setting up MQTT...")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%v", "tcp", Config.MqttServer, Config.MqttPort))
	opts.SetPassword(Config.MqttPass)
	opts.SetUsername(Config.MqttLogin)
	opts.SetClientID("heatpump2mqtt-pub")
	opts.SetWill(Config.mqttWillTopic, "offline", 0, true)
	opts.SetKeepAlive(time.Duration(Config.MqttKeepalive) * time.Second)

	opts.SetCleanSession(true)  // don't want to receive entire backlog of setting changes
	opts.SetAutoReconnect(true) // default, but I want it explicit
	opts.SetConnectRetry(true)
	opts.SetOnConnectHandler(func(mclient mqtt.Client) {
		mqttPublish(mclient, Config.mqttWillTopic, "online", 0)
		if !Config.ListenOnly {
			mclient.Subscribe(getCommandTopic("+"), 0, onGenericCommand)
			mclient.Subscribe(getStatusTopic("+/set"), 0, onAquareaCommand(commandsChannel))
		}
		log.Println("MQTT connected")
	})

	// connect to broker
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Fail to connect broker, %v", token.Error())
		//should not happen - SetConnectRetry=true
	}
	log.Println("Done.")

	return client
}
