package gasbuddy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Station struct {
	Id      int
	Name    string
	Address string
	Fuels   []Fuel
}

type Fuel struct {
	Name    string
	Price   float32
	Updated time.Time
}

type gpResponse struct {
	Station gbStation
}

type gbStation struct {
	Id      int
	Name    string
	City    string
	Address string
	ZipCode string
	APIFuel []gbApiFuel
	Fuels   []gbFuel
}

type gbApiFuel struct {
	Id          int
	DisplayName string
}

type gbFuel struct {
	FuelType    int
	CreditPrice gbCreditPrice
}

type gbCreditPrice struct {
	Amount     float32
	TimePosted string
}

func GetStationInfo(stationId int) (Station, error) {

	// encode the payload
	postBody, _ := json.Marshal(map[string]string{
		"id":         strconv.Itoa(stationId),
		"fuelTypeId": "1",
	})

	// make post request
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("https://www.gasbuddy.com/gaspricemap/station", "application/json", responseBody)
	if err != nil {
		log.Printf("Error calling GasBuddy [%s]!\n", err.Error())
		return Station{}, err
	}

	// ensure call was successfull
	if resp.StatusCode != 200 {
		log.Printf("Calling GasBuddy failed with status code [%d]!\n", resp.StatusCode)
		return Station{}, err
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error processing response from GasBuddy [%s]!\n", err.Error())
		return Station{}, err
	}

	// unmarshall
	var response gpResponse
	json.Unmarshal(body, &response)

	// map
	var station Station
	station.Id = response.Station.Id
	station.Name = response.Station.Name
	station.Address = strings.TrimSpace(response.Station.Address) + ", " + strings.TrimSpace(response.Station.ZipCode) + " " + strings.TrimSpace(response.Station.City)
	station.Fuels = make([]Fuel, 0)
	for _, gbFuel := range response.Station.Fuels {
		if gbFuel.CreditPrice.Amount > 0 {
			var fuel Fuel
			for _, apiFuel := range response.Station.APIFuel {
				if gbFuel.FuelType == apiFuel.Id {
					fuel.Name = apiFuel.DisplayName
				}
			}
			fuel.Price = gbFuel.CreditPrice.Amount
			var timeString = strings.Trim(gbFuel.CreditPrice.TimePosted, "/Date(")
			timeString = strings.Trim(timeString, ")/")
			timeInt, _ := strconv.ParseInt(timeString, 10, 64)
			fuel.Updated = time.UnixMilli(timeInt)
			station.Fuels = append(station.Fuels, fuel)
		}
	}

	return station, nil
}
