package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var Config ConfigStruct

type ConfigStruct struct {
	DeviceName      string `yaml:"deviceName"`      // for HA discovery
	Device          string `yaml:"device"`          // serial port
	ReadInterval    int    `yaml:"readInterval"`    // HP query interval
	ListenOnly      bool   `yaml:"listenOnly"`      // no commands at all
	OptionalPCB     bool   `yaml:"optionalPCB"`     // enable optional PCB emulation
	EnableOSCommand bool   `yaml:"enableOSCommand"` // enable OS commands

	MqttServer     string `yaml:"mqttServer"`
	MqttPort       string `yaml:"mqttPort"`
	MqttLogin      string `yaml:"mqttLogin"`
	MqttPass       string `yaml:"mqttPass"`
	MqttKeepalive  int    `yaml:"mqttKeepalive"`
	MqttTopicBase  string `yaml:"mqttTopicBase"`
	HAAutoDiscover bool   `yaml:"haAutoDiscover"`

	LogMqtt    bool `yaml:"logmqtt"`
	LogHexDump bool `yaml:"loghex"`

	//topics
	mqttWillTopic      string
	mqttLogTopic       string
	mqttValuesTopic    string
	mqttPcbValuesTopic string
	mqttCommandsTopic  string
}

func getStatusTopic(name string) string {
	return fmt.Sprintf("%s/%s", Config.mqttValuesTopic, name)
}

func getCommandTopic(name string) string {
	return fmt.Sprintf("%s/%s", Config.mqttCommandsTopic, name)
}

func getPcbStatusTopic(name string) string {
	return fmt.Sprintf("%s/%s", Config.mqttPcbValuesTopic, name)
}

func LogErrorPause(msg error) {
	log.Println(msg)
	log.Println("Cannot continue - awaiting new Config")
	for {
		time.Sleep(10 * time.Second)
	}
}

func ReadConfig(name string) (config ConfigStruct) {
	_, err := os.Stat(name)
	if err != nil {
		log.Fatalf("Config file is missing: %s ", name)
		return config
	}

	data, err := ioutil.ReadFile(name)
	if err != nil {
		LogErrorPause(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		LogErrorPause(err)
	}

	config.mqttWillTopic = config.MqttTopicBase + "/LWT"
	config.mqttLogTopic = config.MqttTopicBase + "/log"
	config.mqttValuesTopic = config.MqttTopicBase + "/main"
	config.mqttPcbValuesTopic = config.MqttTopicBase + "/optional"
	config.mqttCommandsTopic = config.MqttTopicBase + "/commands"
	log.Println("Config file loaded")

	return config
}
