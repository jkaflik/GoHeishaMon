package main

import (
	"github.com/jkaflik/heatpump2mqtt/internal"
	"io/ioutil"
	"log"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/tarm/serial"
)

const serialTimeout = 2 * time.Second
const optionalPCBSaveTime = 5 * time.Minute
const optionalPCBFile = "/etc/gh/optionalpcb.raw"

var panasonicQuery []byte = []byte{0x71, 0x6c, 0x01, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var optionalPCBQuery []byte = []byte{0xF1, 0x11, 0x01, 0x50, 0x00, 0x00, 0x40, 0xFF, 0xFF, 0xE5, 0xFF, 0xFF, 0x00, 0xFF, 0xEB, 0xFF, 0xFF, 0x00, 0x00}

var serialPort *serial.Port
var commandsChannel chan []byte

var args struct {
	ConfigFile string `arg:"-c,--config" help:"config file path" default:"config.yaml"`
	TopicsFile string `arg:"-t,--topics" help:"topics file path" default:"topics.yaml"`
}

func main() {
	arg.MustParse(&args)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("heatpump2mqtt loading...")
	internal.Config = internal.ReadConfig(args.ConfigFile)

	serialConfig := &serial.Config{Name: internal.Config.Device, Baud: 9600, Parity: serial.ParityEven, StopBits: serial.Stop1, ReadTimeout: serialTimeout}
	var err error
	serialPort, err = serial.OpenPort(serialConfig)
	if err != nil {
		// no point in continuing, wating for new config file
		internal.LogErrorPause(err)
	}
	log.Println("Serial port open")

	commandsChannel = make(chan []byte, 100)
	internal.LoadTopics(args.TopicsFile)
	if internal.Config.OptionalPCB {
		loadOptionalPCB()
	}

	mclient := internal.MakeMQTTConn(commandsChannel)
	internal.RedirectLog(mclient)
	if internal.Config.HAAutoDiscover {
		internal.PublishDiscoveryTopics(mclient)
	}

	queryTicker := time.NewTicker(time.Second * time.Duration(internal.Config.ReadInterval))
	optionPCBSaveTicker := time.NewTicker(optionalPCBSaveTime)
	log.Println("Entering main loop")
	internal.SendCommand(serialPort, panasonicQuery)
	for {
		time.Sleep(serialTimeout)
		internal.ReadSerial(serialPort, mclient)

		var queueLen = len(commandsChannel)
		if queueLen > 10 {
			log.Println("Command queue length: ", len(commandsChannel))
		}

		select {
		case <-optionPCBSaveTicker.C:
			if internal.Config.OptionalPCB {
				saveOptionalPCB()
			}

		case value := <-commandsChannel:
			internal.SendCommand(serialPort, value)

		case <-queryTicker.C:
			commandsChannel <- panasonicQuery

		default:
			if internal.Config.OptionalPCB && !internal.Config.ListenOnly {
				commandsChannel <- optionalPCBQuery
			}
		}
	}
}

func saveOptionalPCB() {
	err := ioutil.WriteFile(optionalPCBFile, optionalPCBQuery, 0644)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Optional PCB data stored")
	}
}

func loadOptionalPCB() {
	data, err := ioutil.ReadFile(optionalPCBFile)
	if err != nil {
		log.Println(err)
	} else {
		optionalPCBQuery = data
		log.Println("Optional PCB data loaded")
	}

}
