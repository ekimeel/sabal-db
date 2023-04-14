package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat/distuv"
	"math/rand"
	"time"
)

type DataPoint struct {
	Timestamp            time.Time
	SupplyAirTemperature float64
	ReturnAirTemperature float64
	SupplyAirFlow        float64
	ReturnAirFlow        float64
	Humidity             float64
}

func main() {
	startDate := time.Date(2023, 3, 18, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 4, 17, 23, 59, 0, 0, time.UTC)

	data := generateAHUTimeSeriesData(startDate, endDate, 1*time.Minute)

	for _, point := range data {
		fmt.Printf("%v: Supply Air Temp: %.2f°F, Return Air Temp: %.2f°F, Supply Air Flow: %.2f CFM, Return Air Flow: %.2f CFM, Humidity: %.2f%%\n",
			point.Timestamp, point.SupplyAirTemperature, point.ReturnAirTemperature, point.SupplyAirFlow, point.ReturnAirFlow, point.Humidity)
	}
}

func generateAHUTimeSeriesData(startDate, endDate time.Time, interval time.Duration) []DataPoint {
	data := []DataPoint{}

	rand.Seed(time.Now().UnixNano())

	supplyAirTempDist := distuv.Normal{Mu: 55, Sigma: 5}
	returnAirTempDist := distuv.Normal{Mu: 75, Sigma: 5}
	supplyAirFlowDist := distuv.Normal{Mu: 5000, Sigma: 500}
	returnAirFlowDist := distuv.Normal{Mu: 4800, Sigma: 500}

	for t := startDate; t.Before(endDate) || t.Equal(endDate); t = t.Add(interval) {
		dataPoint := DataPoint{
			Timestamp:            t,
			SupplyAirTemperature: supplyAirTempDist.Rand(),
			ReturnAirTemperature: returnAirTempDist.Rand(),
			SupplyAirFlow:        supplyAirFlowDist.Rand(),
			ReturnAirFlow:        returnAirFlowDist.Rand(),
			Humidity:             rand.Float64()*(60-30) + 30,
		}
		data = append(data, dataPoint)
	}

	return data
}
