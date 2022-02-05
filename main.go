package main

import (
	"flag"
	"log"
	"time"

	"github.com/dtrumpfheller/gas-tracker/gasbuddy"
	"github.com/dtrumpfheller/gas-tracker/helpers"
	"github.com/dtrumpfheller/gas-tracker/influxdb"
)

var (
	configFile = flag.String("config", "config.yml", "configuration file")
	config     helpers.Config
)

func main() {

	// load arguments into variables
	flag.Parse()

	// load config file
	config = helpers.ReadConfig(*configFile)

	for {
		updateMetrics()
		time.Sleep(time.Duration(config.SleepDuration) * time.Minute)
	}
}

func updateMetrics() {

	log.Println("Getting gas prices... ")
	start := time.Now()

	// get gas prices and export to influxdb
	for _, stationId := range config.Stations {
		station, err := gasbuddy.GetStationInfo(stationId)
		if err == nil {
			influxdb.ExportStation(station, config)
		}
	}

	log.Printf("Finished in %s\n", time.Since(start))
}
