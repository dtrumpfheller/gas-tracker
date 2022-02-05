package influxdb

import (
	"context"
	"log"
	"strconv"

	"github.com/dtrumpfheller/gas-tracker/gasbuddy"
	"github.com/dtrumpfheller/gas-tracker/helpers"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func ExportStation(station gasbuddy.Station, config helpers.Config) {

	// create client objects
	client := influxdb2.NewClient(config.URL, config.Token)
	queryAPI := client.QueryAPI(config.Organization)
	writeAPI := client.WriteAPI(config.Organization, config.Bucket)

	for _, fuel := range station.Fuels {

		// check if entry is already stored, only consider last 7 days
		query := `from(bucket: "` + config.Bucket + `") 
		|> range(start: -2d) 
		|> filter(fn: (r) => r["_measurement"] == "gas_buddy") 
		|> filter(fn: (r) => r["stationId"] == "` + strconv.Itoa(station.Id) + `")
		|> filter(fn: (r) => r["fuelType"] == "` + fuel.Name + `")
		|> last()`
		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			log.Printf("Error calling InfluxDB [%s]!\n", err.Error())
			return
		}

		for result.Next() {
			if fuel.Updated.After(result.Record().Time()) {
				log.Printf("Updating [%s] for [%s] - [%s]\n", fuel.Name, station.Name, station.Address)
				// write gas price
				point := influxdb2.NewPointWithMeasurement("gas_buddy").
					AddTag("stationId", strconv.Itoa(station.Id)).
					AddTag("name", station.Name).
					AddTag("address", station.Address).
					AddTag("fuelType", fuel.Name).
					AddField("fuelPrice", fuel.Price).
					SetTime(fuel.Updated)
				writeAPI.WritePoint(point)
			}
		}
	}

	// force all unwritten data to be sent
	writeAPI.Flush()

	// ensures background processes finishes
	client.Close()
}
