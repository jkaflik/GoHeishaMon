package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type configStruct struct {
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
	return fmt.Sprintf("%s/%s", config.mqttValuesTopic, name)
}

func getCommandTopic(name string) string {
	return fmt.Sprintf("%s/%s", config.mqttCommandsTopic, name)
}

func getPcbStatusTopic(name string) string {
	return fmt.Sprintf("%s/%s", config.mqttPcbValuesTopic, name)
}

func logErrorPause(msg error) {
	log.Println(msg)
	log.Println("Cannot continue - awaiting new config")
	for {
		time.Sleep(10 * time.Second)
	}
}

func readConfig(name string) configStruct {
	_, err := os.Stat(name)
	if err != nil {
		log.Printf("Config file is missing: %s ", name)
		updateConfig(name)
		// it's either it reboots or we can't continue
		for {
			time.Sleep(10 * time.Second)
		}
	}

	var config configStruct

	data, err := ioutil.ReadFile(name)
	if err != nil {
		logErrorPause(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logErrorPause(err)
	}

	config.mqttWillTopic = config.MqttTopicBase + "/LWT"
	config.mqttLogTopic = config.MqttTopicBase + "/log"
	config.mqttValuesTopic = config.MqttTopicBase + "/main"
	config.mqttPcbValuesTopic = config.MqttTopicBase + "/optional"
	config.mqttCommandsTopic = config.MqttTopicBase + "/commands"
	log.Println("Config file loaded")

	return config
}

func updateConfig(name string) {
	log.Println("Config updater - checking USB media")
	err := exec.Command("/usr/bin/usb_mount.sh").Run()
	if err != nil {
		return
	}
	defer exec.Command("/usr/bin/usb_umount.sh").Run()

	_, err = os.Stat("/mnt/usb/GoHeishaMonConfig.new")
	if err != nil {
		return
	}
	if !bytes.Equal(getFileChecksum(name), getFileChecksum("/mnt/usb/GoHeishaMonConfig.new")) {
		log.Println("Updated configuration detected on USB media... will reboot")

		err = exec.Command("/bin/cp", "/mnt/usb/GoHeishaMonConfig.new", name).Run()
		if err != nil {
			log.Printf("Can't update config file %s", name)
			return
		}
		exec.Command("sync").Run()
		exec.Command("/usr/bin/usb_umount.sh").Run()
		exec.Command("reboot").Run()
	}
}

func getFileChecksum(f string) []byte {
	input := strings.NewReader(f)

	hash := md5.New()
	if _, err := io.Copy(hash, input); err != nil {
		log.Println(err)
	}
	return hash.Sum(nil)
}

func updateConfigLoop(name string) {
	for {
		updateConfig(name)
		time.Sleep(time.Minute * 5)
	}
}
