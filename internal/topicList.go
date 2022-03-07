package internal

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

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

func LoadTopics(topicFile string) {
	log.Print("Loading topic data...")
	data, err := ioutil.ReadFile(topicFile)
	if err != nil {
		LogErrorPause(err)
	}

	err = yaml.Unmarshal(data, &allTopics)
	if err != nil {
		LogErrorPause(err)
	}

	topicNameLookup = make(map[string]topicData)
	for _, val := range allTopics {
		topicNameLookup[val.SensorName] = val
	}
	log.Println(" loaded.")
}
