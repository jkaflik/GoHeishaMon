package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

const topicsFileOther = "/etc/gh/topics.yaml"
const topicsFileWindows = "topics.yaml"

var allTopics []topicData
var topicNameLookup map[string]topicData

type topicData struct {
	SensorName     string   `yaml:"sensorName"`
	DecodeFunction string   `yaml:"decodeFunction"`
	EncodeFunction string   `yaml:"encodeFunction"`
	DecodeOffset   int      `yaml:"decodeOffset"`
	DisplayUnit    string   `yaml:"displayUnit"`
	Values         []string `yaml:"values"`
	Min            int      `yaml:"min"`
	Max            int      `yaml:"max"`
	Step           int      `yaml:"step"`
	StateClass     string   `yaml:"stateClass"`
	currentValue   string
}

func loadTopics(topicFile string) {
	log.Print("Loading topic data...")
	data, err := ioutil.ReadFile(topicFile)
	if err != nil {
		logErrorPause(err)
	}

	err = yaml.Unmarshal(data, &allTopics)
	if err != nil {
		logErrorPause(err)
	}

	topicNameLookup = make(map[string]topicData)
	for _, val := range allTopics {
		topicNameLookup[val.SensorName] = val
	}
	log.Println(" loaded.")
}
